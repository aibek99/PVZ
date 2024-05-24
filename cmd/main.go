package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"Homework-1/internal/app"
	"Homework-1/internal/app/command/order/accept"
	"Homework-1/internal/app/command/order/issue"
	"Homework-1/internal/app/command/order/order"
	"Homework-1/internal/app/command/order/receive"
	"Homework-1/internal/app/command/order/returns"
	"Homework-1/internal/app/command/order/turnin"
	"Homework-1/internal/app/command/pvz/create"
	"Homework-1/internal/app/command/pvz/get"
	"Homework-1/internal/cache"
	"Homework-1/internal/config"
	"Homework-1/internal/connection"
	"Homework-1/internal/database/postgres"
	"Homework-1/internal/kafka/consumer"
	"Homework-1/internal/kafka/producer"
	"Homework-1/internal/metrics"
	"Homework-1/internal/server"
	file2 "Homework-1/internal/storage/order/file"
	"Homework-1/internal/storage/pvz/file"
)

func main() {
	ctx := context.Background()

	if err := bootstrap(ctx); err != nil {
		log.Fatalf("[main] bootstrap: %v", err)
	}
}

func bootstrap(ctx context.Context) error {
	defer log.Printf("[main][bootstrap] Graceful Shutdown complete\n")

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var option, ENV string
	flag.StringVar(&option, "option", "rest", `There are 2 options: 1 - "rest", 2 - "cli"`)
	flag.StringVar(&ENV, "env", "prod", `There are 3 env: 1 - "prod", 2 - "local", 3 - "testing"`)
	flag.Parse()

	ENV = ".env." + ENV
	if err := godotenv.Load(ENV); err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}

	switch {
	case option == "rest":
		configPath := os.Getenv("CONFIG_PATH")

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("config.LoadConfig: %w", err)
		}

		ctxTime, timeCancel := context.WithTimeout(ctx, 10*time.Second)
		defer timeCancel()

		exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
		if err != nil {
			log.Printf("Cannot create Jaeger exporter: %s", err.Error())
		}

		tp := trace.NewTracerProvider(
			trace.WithBatcher(exporter),
			trace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("Homework-1"),
			)),
		)
		otel.SetTracerProvider(tp)

		defer func() {
			if err = otel.GetTracerProvider().(*trace.TracerProvider).Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracing provider: %v", err)
			}
		}()

		kafkaProducer, err := producer.NewProducer(cfg.Kafka)
		if err != nil {
			log.Printf("[main][bootstrap] Failed to create Kafka Producer: %v\n", err)
		}

		defer func(_ *producer.Producer) {
			err = kafkaProducer.Close()
			if err != nil {
				log.Printf("[main][bootstrap] kafkaProducer.Close: %v\n", err)
			}
		}(kafkaProducer)

		// Create Kafka Consumer
		kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka.Brokers, os.Stdout)
		if err != nil {
			errPro := kafkaProducer.Close()
			if errPro != nil {
				log.Printf("[main][bootstrap] Failed to close Kafka Producer: %v\n", errPro)
			}

			log.Printf("[main][bootstrap] Failed to create Kafka Consumer: %v\n", err)
		}
		// Start Kafka message Consumer
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err = kafkaConsumer.ReadMessages(ctx, cfg.Kafka.Topic); err != nil {
				// Last closing Producer in order to prevent data lost
				if err = kafkaProducer.Close(); err != nil {
					log.Printf("[main][bootstrap] Error closing Kafka Producer: %v\n", err)
				}
				// First closing Consumer
				if err = kafkaConsumer.Close(); err != nil {
					log.Printf("[main][bootstrap] Error closing Kafka Consumer: %v\n", err)
				}

				log.Printf("[main][bootstrap] Closed Consumer: %v\n", err)
			}
		}()

		psqlDB, err := connection.NewDB(ctxTime, cfg.Postgres)
		if err != nil {
			return fmt.Errorf("connection.NewDB: %w", err)
		}

		defer func() {
			err = psqlDB.Close()
			if err != nil {
				log.Printf("[main][bootstrap] psqlDB.Close: %v\n", err)
			}
		}()

		// if our server stops to run then metrics will not affect from server, it will run forever
		go func() {
			err = metrics.Listen(cfg.Server.MetricsPort)
			log.Printf("[metrics][Listen] %v", err)
		}()

		var rdb connection.Cache
		if cfg.CacheType == "inMemory" {
			rdb = connection.NewInMemoryCache(ctx, cfg.InMemoryCache)
		} else {
			rdb, err = connection.NewCache(ctx, cfg.Redis)
			if err != nil {
				return fmt.Errorf("connection.NewCache: %w", err)
			}
		}

		defer func() {
			err = rdb.Close()
			if err != nil {
				log.Printf("redis.Close: %v\n", err)
			}
		}()

		dataStore := postgres.NewDataStore(psqlDB)
		cacheStore := cache.NewClientRDRepository(rdb)

		source := server.NewServer(cfg, dataStore, cacheStore, kafkaProducer)

		err = source.Run(ctx)
		if err != nil {
			return fmt.Errorf("source.Run: %w", err)
		}

		wg.Wait()
	case option == "cli":
		pvzStoreName := os.Getenv("PVZ_STORE_NAME")
		pvzStore, err := file.New(pvzStoreName)
		if err != nil {
			return fmt.Errorf("pvzStore.New: %w", err)
		}

		defer func() {
			if err = pvzStore.Close(); err != nil {
				log.Println(fmt.Errorf("[main][bootstrap] pvzStore.Close: %w", err))
			}
		}()

		orderStoreName := os.Getenv("ORDER_STORE_NAME")
		orderStore, err := file2.New(orderStoreName)
		if err != nil {
			return fmt.Errorf("orderStore.New: %w", err)
		}

		defer func() {
			if err = orderStore.Close(); err != nil {
				log.Println(fmt.Errorf("[main][bootstrap] orderStore.Close: %w", err))
			}
		}()

		cli, err := app.New(orderStore, pvzStore)
		if err != nil {
			return fmt.Errorf("app.New: %w", err)
		}
		cli.AddCommand(accept.New(orderStore))
		cli.AddCommand(issue.New(orderStore))
		cli.AddCommand(order.New(orderStore))
		cli.AddCommand(receive.New(orderStore))
		cli.AddCommand(returns.New(orderStore))
		cli.AddCommand(turnin.New(orderStore))
		cli.AddCommand(create.New(pvzStore))
		cli.AddCommand(get.New(pvzStore))

		err = cli.PVZRun(ctx)
		if err != nil {
			return fmt.Errorf("cli.PVZRun: %w", err)
		}
	default:
		return errors.New("invalid option")
	}

	return nil
}
