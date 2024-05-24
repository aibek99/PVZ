//go:build integration
// +build integration

package box

import (
	"bytes"
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"Homework-1/internal/box/delivery"
	BoxUseCase "Homework-1/internal/box/usecase"
	"Homework-1/internal/cache"
	"Homework-1/internal/config"
	"Homework-1/internal/connection"
	"Homework-1/internal/database/postgres"
	"Homework-1/internal/kafka"
	"Homework-1/internal/kafka/consumer"
	"Homework-1/internal/kafka/producer"
	"Homework-1/internal/server"
	"Homework-1/pkg/api/abstract"
	"Homework-1/pkg/api/box_v1"
)

// TestIntegrationBoxHandler_CreateBox is
func TestIntegrationBoxHandler_CreateBox(t *testing.T) {
	// get connection from database
	ENV := ".env.prod"
	if os.Getenv("ENV") == "testing" {
		ENV = ".env.testing"
	}

	err := godotenv.Load("./../../" + ENV)
	require.NoError(t, err)

	configPath := os.Getenv("CONFIG_PATH")
	if ENV == ".env.prod" {
		configPath = "./../." + configPath
	}

	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	tdb, err := connection.NewTDB(context.Background(), cfg.Postgres, t)
	require.NoError(t, err)

	rdb, err := connection.NewCache(context.Background(), cfg.Redis)
	require.NoError(t, err)

	defer tdb.Close()

	setup := func(t *testing.T, tdb *connection.TDB) (kafka.Producer, box_v1.BoxServiceClient) {
		kafkaProducer, err := producer.NewProducer(cfg.Kafka)
		require.NoError(t, err)
		dataStore := postgres.NewDataStore(tdb)
		cacheStore := cache.NewClientRDRepository(rdb)
		useCase := BoxUseCase.NewBoxUseCase(dataStore, cacheStore)
		handler := delivery.NewBoxHandler(useCase)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.KafkaInterceptor(kafkaProducer)))
		box_v1.RegisterBoxServiceServer(grpcServer, handler)
		listener, err := net.Listen("tcp", "localhost:9000")
		require.NoError(t, err)
		go func() {
			require.NoError(t, grpcServer.Serve(listener))
		}()
		conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		client := box_v1.NewBoxServiceClient(conn)
		t.Cleanup(func() {
			grpcServer.Stop()
			conn.Close()
			listener.Close()
		})
		return kafkaProducer, client
	}

	tests := []struct {
		description   string
		requestBody   box_v1.BoxCreateRequest
		wantCreateErr error
		wantGetResp   *box_v1.BoxAllInfo
		wantGetErr    error
		assertFunc    func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer)
	}{
		{
			description: "Success to created box_v1",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantCreateErr: nil,
			wantGetResp: &box_v1.BoxAllInfo{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantGetErr: nil,
			assertFunc: func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer) {
				assert.Equal(t, wantGetResp.Box.Name, getResp.Box.Name)
				assert.Equal(t, wantGetResp.Box.Weight, getResp.Box.Weight)
				assert.Equal(t, wantGetResp.Box.Cost, getResp.Box.Cost)
				assert.Equal(t, wantGetResp.Box.IsCheck, getResp.Box.IsCheck)

				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true,"weight":10.1}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {"boxID":`+strconv.Itoa(int(getResp.ID))+`}`)

				t.Cleanup(func() {
					err = tdb.DropRowByID(context.Background(), "box", getResp.ID)
					require.NoError(t, err)
				})
			},
		},
		{
			description: "Failed to create Box",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
				},
			},
			wantCreateErr: status.Errorf(codes.InvalidArgument, "reqvalidator.ValidateRequest Key: 'Request.Weight' Error:Field validation for 'Weight' failed on the 'required' tag"),
			wantGetResp:   nil,
			wantGetErr:    status.Errorf(codes.InvalidArgument, "invalid request id"),
			assertFunc: func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer) {
				assert.Equal(t, wantGetResp, getResp)
				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {"boxID":`)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			kafkaProducer, client := setup(t, tdb)

			var messages bytes.Buffer
			kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka.Brokers, &messages)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err = kafkaConsumer.ReadMessages(ctx, cfg.Kafka.Topic); err != nil {
					err = kafkaProducer.Close()
					require.NoError(t, err)
					err = kafkaConsumer.Close()
					require.NoError(t, err)
				}
			}()

			createResp, err := client.CreateBox(context.Background(), &tc.requestBody)
			assert.Equal(t, tc.wantCreateErr, err)

			boxID := int64(-1)
			if err == nil {
				boxID, err = strconv.ParseInt(createResp.Message, 10, 64)
				require.NoError(t, err)
			}

			getResp, err := client.GetBoxByID(ctx, &box_v1.BoxIDRequest{BoxID: boxID})
			assert.Equal(t, tc.wantGetErr, err)
			wg.Wait()

			tc.assertFunc(t, getResp, tc.wantGetResp, &messages)
		})
	}
}

// TestIntegrationBoxHandler_DeleteBox is
func TestIntegrationBoxHandler_DeleteBox(t *testing.T) {
	// get connection from database
	ENV := ".env.prod"
	if os.Getenv("ENV") == "testing" {
		ENV = ".env.testing"
	}

	err := godotenv.Load("./../../" + ENV)
	require.NoError(t, err)

	configPath := os.Getenv("CONFIG_PATH")
	if ENV == ".env.prod" {
		configPath = "./../." + configPath
	}

	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	tdb, err := connection.NewTDB(context.Background(), cfg.Postgres, t)
	require.NoError(t, err)

	rdb, err := connection.NewCache(context.Background(), cfg.Redis)
	require.NoError(t, err)

	defer tdb.Close()

	setup := func(t *testing.T, tdb *connection.TDB) (kafka.Producer, box_v1.BoxServiceClient) {
		kafkaProducer, err := producer.NewProducer(cfg.Kafka)
		require.NoError(t, err)
		dataStore := postgres.NewDataStore(tdb)
		cacheStore := cache.NewClientRDRepository(rdb)
		useCase := BoxUseCase.NewBoxUseCase(dataStore, cacheStore)
		handler := delivery.NewBoxHandler(useCase)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.KafkaInterceptor(kafkaProducer)))
		box_v1.RegisterBoxServiceServer(grpcServer, handler)
		listener, err := net.Listen("tcp", "localhost:9000")
		require.NoError(t, err)
		go func() {
			require.NoError(t, grpcServer.Serve(listener))
		}()
		conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		client := box_v1.NewBoxServiceClient(conn)
		t.Cleanup(func() {
			grpcServer.Stop()
			conn.Close()
			listener.Close()
		})
		return kafkaProducer, client
	}

	tests := []struct {
		description    string
		requestBody    box_v1.BoxCreateRequest
		wantCreateErr  error
		wantDeleteResp *abstract.MessageResponse
		wantDeleteErr  error
		wantGetResp    *box_v1.BoxAllInfo
		wantGetErr     error
		assertFunc     func(t *testing.T, deleteResp *abstract.MessageResponse, wantDeleteResp *abstract.MessageResponse, messages *bytes.Buffer, boxID int64)
	}{
		{
			description: "Successfully deleted box_v1",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantCreateErr:  nil,
			wantDeleteResp: &abstract.MessageResponse{Message: "Successfully Deleted Box\n"},
			wantDeleteErr:  nil,
			wantGetResp:    nil,
			wantGetErr:     status.Errorf(codes.NotFound, "Error: Box not found"),
			assertFunc: func(t *testing.T, deleteResp *abstract.MessageResponse, wantDeleteResp *abstract.MessageResponse, messages *bytes.Buffer, boxID int64) {
				assert.Equal(t, wantDeleteResp.Message, deleteResp.Message)
				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true,"weight":10.1}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/DeleteBox, Request: {"boxID":`+strconv.Itoa(int(boxID))+`}`)
				assert.Contains(t, messagesSlice[2], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {"boxID":`+strconv.Itoa(int(boxID))+`}`)
			},
		},
		{
			description: "Failed to delete Box",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
				},
			},
			wantCreateErr:  status.Errorf(codes.InvalidArgument, "reqvalidator.ValidateRequest Key: 'Request.Weight' Error:Field validation for 'Weight' failed on the 'required' tag"),
			wantDeleteResp: nil,
			wantDeleteErr:  status.Errorf(codes.NotFound, "Error: Box not found"),
			wantGetResp:    nil,
			wantGetErr:     status.Errorf(codes.NotFound, "Error: Box not found"),
			assertFunc: func(t *testing.T, deleteResp *abstract.MessageResponse, wantDeleteResp *abstract.MessageResponse, messages *bytes.Buffer, boxID int64) {
				assert.Equal(t, wantDeleteResp, deleteResp)
				message := messages.String()
				messagesSlice := strings.Split(message, "\n")
				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/DeleteBox, Request: {}`)
				assert.Contains(t, messagesSlice[2], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {}`)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			kafkaProducer, client := setup(t, tdb)
			var messages bytes.Buffer
			kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka.Brokers, &messages)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err = kafkaConsumer.ReadMessages(ctx, cfg.Kafka.Topic); err != nil {
					err = kafkaProducer.Close()
					require.NoError(t, err)
					err = kafkaConsumer.Close()
					require.NoError(t, err)
				}
			}()

			createResp, err := client.CreateBox(ctx, &tc.requestBody)
			assert.Equal(t, tc.wantCreateErr, err)

			boxID := int64(0)
			if err == nil {
				boxID, err = strconv.ParseInt(createResp.Message, 10, 64)
				require.NoError(t, err)
			}

			deleteResp, err := client.DeleteBox(ctx, &box_v1.BoxIDRequest{BoxID: boxID})
			assert.Equal(t, tc.wantDeleteErr, err)

			getResp, err := client.GetBoxByID(ctx, &box_v1.BoxIDRequest{BoxID: boxID})
			assert.Equal(t, tc.wantGetErr, err)
			assert.Equal(t, tc.wantGetResp, getResp)
			wg.Wait()

			tc.assertFunc(t, deleteResp, tc.wantDeleteResp, &messages, boxID)
			t.Cleanup(func() {
				err = tdb.DropRowByID(context.Background(), "box", boxID)
				require.NoError(t, err)
			})
		})
	}
}

// TestIntegrationBoxHandler_GetBoxByID is
func TestIntegrationBoxHandler_GetBoxByID(t *testing.T) {
	// get connection from database
	ENV := ".env.prod"
	if os.Getenv("ENV") == "testing" {
		ENV = ".env.testing"
	}

	err := godotenv.Load("./../../" + ENV)
	require.NoError(t, err)

	configPath := os.Getenv("CONFIG_PATH")
	if ENV == ".env.prod" {
		configPath = "./../." + configPath
	}

	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	tdb, err := connection.NewTDB(context.Background(), cfg.Postgres, t)
	require.NoError(t, err)

	rdb, err := connection.NewCache(context.Background(), cfg.Redis)
	require.NoError(t, err)

	defer tdb.Close()

	setup := func(t *testing.T, tdb *connection.TDB) (kafka.Producer, box_v1.BoxServiceClient) {
		kafkaProducer, err := producer.NewProducer(cfg.Kafka)
		require.NoError(t, err)
		dataStore := postgres.NewDataStore(tdb)
		cacheStore := cache.NewClientRDRepository(rdb)
		useCase := BoxUseCase.NewBoxUseCase(dataStore, cacheStore)
		handler := delivery.NewBoxHandler(useCase)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.KafkaInterceptor(kafkaProducer)))
		box_v1.RegisterBoxServiceServer(grpcServer, handler)
		listener, err := net.Listen("tcp", "localhost:9000")
		require.NoError(t, err)
		go func() {
			require.NoError(t, grpcServer.Serve(listener))
		}()
		conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		client := box_v1.NewBoxServiceClient(conn)
		t.Cleanup(func() {
			grpcServer.Stop()
			conn.Close()
			listener.Close()
		})
		return kafkaProducer, client
	}

	tests := []struct {
		description   string
		requestBody   box_v1.BoxCreateRequest
		wantCreateErr error
		wantGetResp   *box_v1.BoxAllInfo
		wantGetErr    error
		assertFunc    func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer)
	}{
		{
			description: "Successfully got box_v1 data",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantCreateErr: nil,
			wantGetResp: &box_v1.BoxAllInfo{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantGetErr: nil,
			assertFunc: func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer) {
				assert.Equal(t, wantGetResp.Box.Name, getResp.Box.Name)
				assert.Equal(t, wantGetResp.Box.Weight, getResp.Box.Weight)
				assert.Equal(t, wantGetResp.Box.Cost, getResp.Box.Cost)
				assert.Equal(t, wantGetResp.Box.IsCheck, getResp.Box.IsCheck)

				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true,"weight":10.1}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {"boxID":`+strconv.Itoa(int(getResp.ID))+`}`)
				err = tdb.DropRowByID(context.Background(), "box", getResp.ID)
				require.NoError(t, err)
			},
		},
		{
			description: "Fail",
			requestBody: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			wantCreateErr: nil,
			wantGetResp:   nil,
			wantGetErr:    status.Errorf(codes.NotFound, "Error: Box not found"),
			assertFunc: func(t *testing.T, getResp *box_v1.BoxAllInfo, wantGetResp *box_v1.BoxAllInfo, messages *bytes.Buffer) {
				assert.Equal(t, wantGetResp, getResp)
				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true,"weight":10.1}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/GetBoxByID, Request: {}`)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			kafkaProducer, client := setup(t, tdb)
			var messages bytes.Buffer
			kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka.Brokers, &messages)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err = kafkaConsumer.ReadMessages(ctx, cfg.Kafka.Topic); err != nil {
					err = kafkaProducer.Close()
					require.NoError(t, err)
					err = kafkaConsumer.Close()
					require.NoError(t, err)
				}
			}()

			createResp, err := client.CreateBox(ctx, &tc.requestBody)
			assert.Equal(t, tc.wantCreateErr, err)

			boxID, err := strconv.ParseInt(createResp.Message, 10, 64)
			require.NoError(t, err)

			if tc.description == "Fail" {
				err = tdb.DropRowByID(context.Background(), "box", boxID)
				require.NoError(t, err)
				boxID = 0
			}

			getResp, err := client.GetBoxByID(ctx, &box_v1.BoxIDRequest{BoxID: boxID})
			wg.Wait()

			tc.assertFunc(t, getResp, tc.wantGetResp, &messages)

		})
	}
}

// TestIntegrationBoxHandler_ListBoxes is
func TestIntegrationBoxHandler_ListBoxes(t *testing.T) {
	// get connection from database
	ENV := ".env.prod"
	if os.Getenv("ENV") == "testing" {
		ENV = ".env.testing"
	}

	err := godotenv.Load("./../../" + ENV)
	require.NoError(t, err)

	configPath := os.Getenv("CONFIG_PATH")
	if ENV == ".env.prod" {
		configPath = "./../." + configPath
	}

	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	tdb, err := connection.NewTDB(context.Background(), cfg.Postgres, t)
	require.NoError(t, err)

	rdb, err := connection.NewCache(context.Background(), cfg.Redis)
	require.NoError(t, err)

	defer tdb.Close()

	setup := func(t *testing.T, tdb *connection.TDB) (kafka.Producer, box_v1.BoxServiceClient) {
		kafkaProducer, err := producer.NewProducer(cfg.Kafka)
		require.NoError(t, err)
		dataStore := postgres.NewDataStore(tdb)
		cacheStore := cache.NewClientRDRepository(rdb)
		useCase := BoxUseCase.NewBoxUseCase(dataStore, cacheStore)
		handler := delivery.NewBoxHandler(useCase)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.KafkaInterceptor(kafkaProducer)))
		box_v1.RegisterBoxServiceServer(grpcServer, handler)
		listener, err := net.Listen("tcp", "localhost:9000")
		require.NoError(t, err)
		go func() {
			require.NoError(t, grpcServer.Serve(listener))
		}()
		conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		client := box_v1.NewBoxServiceClient(conn)
		t.Cleanup(func() {
			grpcServer.Stop()
			conn.Close()
			listener.Close()
		})
		err = tdb.TruncateTable(context.Background(), "box")
		require.NoError(t, err)
		return kafkaProducer, client
	}

	tests := []struct {
		description   string
		createRequest box_v1.BoxCreateRequest
		listRequest   abstract.Page
		wantCreateErr error
		wantListResp  *box_v1.BoxListResponse
		wantListErr   error
		assertFunc    func(t *testing.T, listResp *box_v1.BoxListResponse, wantListResp *box_v1.BoxListResponse, messages *bytes.Buffer)
	}{
		{
			description: "Successfully got list of boxes",
			createRequest: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
					Weight:  10.1,
				},
			},
			listRequest: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantCreateErr: nil,
			wantListResp: &box_v1.BoxListResponse{
				BoxAllInfo: []*box_v1.BoxAllInfo{{
					Box: &box_v1.Box{
						Name:    "test",
						Cost:    12.1,
						IsCheck: true,
						Weight:  10.1,
					},
				}},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 1,
					},
					TotalItems: 1,
				},
			},
			wantListErr: nil,
			assertFunc: func(t *testing.T, listResp *box_v1.BoxListResponse, wantListResp *box_v1.BoxListResponse, messages *bytes.Buffer) {
				assert.Equal(t, len(wantListResp.BoxAllInfo), len(listResp.BoxAllInfo))
				assert.Equal(t, wantListResp.BoxAllInfo[0].Box.Name, listResp.BoxAllInfo[0].Box.Name)
				assert.Equal(t, wantListResp.BoxAllInfo[0].Box.Cost, listResp.BoxAllInfo[0].Box.Cost)
				assert.Equal(t, wantListResp.BoxAllInfo[0].Box.IsCheck, listResp.BoxAllInfo[0].Box.IsCheck)
				assert.Equal(t, wantListResp.BoxAllInfo[0].Box.Weight, listResp.BoxAllInfo[0].Box.Weight)
				assert.Equal(t, wantListResp.Pagination.Page.ItemsPerPage, listResp.Pagination.Page.ItemsPerPage)
				assert.Equal(t, wantListResp.Pagination.Page.CurrentPage, listResp.Pagination.Page.CurrentPage)
				assert.Equal(t, wantListResp.Pagination.TotalItems, listResp.Pagination.TotalItems)

				message := messages.String()
				messagesSlice := strings.Split(message, "\n")

				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true,"weight":10.1}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/ListBoxes, Request: {"currentPage":1,"itemsPerPage":10}`)

				t.Cleanup(func() {
					err = tdb.DropRowByID(context.Background(), "box", listResp.BoxAllInfo[0].ID)
					require.NoError(t, err)
				})
			},
		},
		{
			description: "Empty list",
			createRequest: box_v1.BoxCreateRequest{
				Box: &box_v1.Box{
					Name:    "test",
					Cost:    12.1,
					IsCheck: true,
				},
			},
			listRequest: abstract.Page{
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			wantCreateErr: status.Errorf(codes.InvalidArgument, "reqvalidator.ValidateRequest Key: 'Request.Weight' Error:Field validation for 'Weight' failed on the 'required' tag"),
			wantListResp: &box_v1.BoxListResponse{
				BoxAllInfo: []*box_v1.BoxAllInfo{},
				Pagination: &abstract.Pagination{
					Page: &abstract.Page{
						CurrentPage:  1,
						ItemsPerPage: 0,
					},
					TotalItems: 0,
				},
			},
			wantListErr: nil,
			assertFunc: func(t *testing.T, listResp *box_v1.BoxListResponse, wantListResp *box_v1.BoxListResponse, messages *bytes.Buffer) {
				assert.Equal(t, len(wantListResp.BoxAllInfo), len(listResp.BoxAllInfo))
				assert.Equal(t, wantListResp.Pagination.Page.ItemsPerPage, listResp.Pagination.Page.ItemsPerPage)
				assert.Equal(t, wantListResp.Pagination.Page.CurrentPage, listResp.Pagination.Page.CurrentPage)
				assert.Equal(t, wantListResp.Pagination.TotalItems, listResp.Pagination.TotalItems)

				message := messages.String()
				messagesSlice := strings.Split(message, "\n")
				assert.Contains(t, messagesSlice[0], `Received Kafka Message: Method: /BoxService/CreateBox, Request: {"box":{"name":"test","cost":12.1,"isCheck":true}}`)
				assert.Contains(t, messagesSlice[1], `Received Kafka Message: Method: /BoxService/ListBoxes, Request: {"currentPage":1,"itemsPerPage":10}`)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			kafkaProducer, client := setup(t, tdb)
			var messages bytes.Buffer
			kafkaConsumer, err := consumer.NewConsumer(cfg.Kafka.Brokers, &messages)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err = kafkaConsumer.ReadMessages(ctx, cfg.Kafka.Topic); err != nil {
					err = kafkaProducer.Close()
					require.NoError(t, err)
					err = kafkaConsumer.Close()
					require.NoError(t, err)
				}
			}()

			_, err = client.CreateBox(ctx, &tc.createRequest)
			assert.Equal(t, tc.wantCreateErr, err)

			listResp, err := client.ListBoxes(ctx, &tc.listRequest)
			assert.Equal(t, tc.wantListErr, err)
			wg.Wait()

			tc.assertFunc(t, listResp, tc.wantListResp, &messages)
		})
	}
}
