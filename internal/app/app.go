package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"Homework-1/internal/app/command"
	file2 "Homework-1/internal/storage/order/file"
	"Homework-1/internal/storage/pvz/file"
)

const helpConst = "help"
const lineConst = "----------------------------------------------------------------------------\n"

// CLI is
type CLI struct {
	Commands   map[string]command.Commander
	orderStore *file2.Storage
	pvzStore   *file.Storage
}

type pvzChan struct {
	ctx           context.Context
	wgLogger      sync.WaitGroup
	wgWorker      sync.WaitGroup
	logMessage    chan string
	errMessage    chan error
	createPVZArgs chan []string
	getPVZArgs    chan []string
}

// New is
func New(orderStore *file2.Storage, pvzStore *file.Storage) (*CLI, error) {
	cli := &CLI{
		Commands:   make(map[string]command.Commander),
		orderStore: orderStore,
		pvzStore:   pvzStore,
	}
	return cli, nil
}

// AddCommand is
func (c *CLI) AddCommand(cmd command.Commander) {
	c.Commands[cmd.Name()] = cmd
}

// OrderRun is
func (c *CLI) OrderRun(args []string) (*string, error) {
	if len(args) < 2 {
		return nil, errors.New("app.Run: arguments weren't provided")
	}
	operation := args[1]

	if operation == helpConst {
		var msg string
		msg += "\n" + lineConst
		for _, val := range c.Commands {
			msg = msg + fmt.Sprintf("%v\n", val.Description())
			_, err := val.CheckFormat(args, true)
			if err != nil {
				return nil, fmt.Errorf("%s.CheckFormat: %w", operation, err)
			}
			msg = msg + val.Help()
			msg = msg + lineConst
		}
		return &msg, nil
	}
	cmd, ok := c.Commands[operation]
	if !ok {
		return nil, fmt.Errorf("invalid command: %s", operation)
	}

	callHelp := len(args) > 2 && strings.HasSuffix(args[2], helpConst)

	arguments, err := cmd.CheckFormat(args[2:], false)
	if err != nil {
		return nil, fmt.Errorf("%s.CheckFormat: %w", operation, err)
	}

	if callHelp {
		help := cmd.Help()
		return &help, nil
	}

	var msg *string
	msg, err = cmd.Do(arguments)
	if err != nil {
		return nil, fmt.Errorf("%s.Do: %w", operation, err)
	}
	return msg, nil
}

// PVZRun is
func (c *CLI) PVZRun(ctx context.Context) error {
	communicate := pvzChan{
		ctx:           ctx,
		wgLogger:      sync.WaitGroup{},
		wgWorker:      sync.WaitGroup{},
		logMessage:    make(chan string),
		errMessage:    make(chan error),
		createPVZArgs: make(chan []string),
		getPVZArgs:    make(chan []string),
	}

	communicate.wgLogger.Add(1)
	go func() {
		defer communicate.wgLogger.Done()
		c.logger(&communicate)
	}()

	communicate.wgLogger.Add(1)
	go func() {
		defer communicate.wgLogger.Done()
		c.logErrors(&communicate)
	}()

	communicate.wgWorker.Add(1)
	go func() {
		defer communicate.wgWorker.Done()
		c.readInput(&communicate)
	}()

	communicate.wgWorker.Add(1)
	go func() {
		defer communicate.wgWorker.Done()
		c.createPVZ(&communicate)
	}()

	communicate.wgWorker.Add(1)
	go func() {
		defer communicate.wgWorker.Done()
		c.getPVZ(&communicate)
	}()

	<-communicate.ctx.Done()
	communicate.wgWorker.Wait()
	close(communicate.logMessage)
	close(communicate.errMessage)
	communicate.wgLogger.Wait()
	return nil
}

func (c *CLI) logger(communicate *pvzChan) {
	log.Println("[app][logger] Monitoring Log Messages: started")
	for msg := range communicate.logMessage {
		log.Printf("[app][logger] Monitoring Log: %s\n", msg)
	}
}

func (c *CLI) logErrors(communicate *pvzChan) {
	log.Println("[app][logErrors] Monitoring Error Messages: started")
	for err := range communicate.errMessage {
		log.Printf("[app][logErrors] Monitoring Error: %v\n", err)
	}
}

func (c *CLI) readInput(communicate *pvzChan) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	communicate.logMessage <- "goroutine: app.ReadInput: started"

	inputChan := make(chan string)
	go input(communicate, inputChan, ticker)

	for {
		select {
		case <-communicate.ctx.Done():
			close(communicate.createPVZArgs)
			close(communicate.getPVZArgs)
			return
		case input := <-inputChan:
			input = strings.TrimSpace(input)

			args := []string{"path"}
			args = append(args, strings.Split(input, " ")...)
			operation := args[1]

			communicate.logMessage <- fmt.Sprintf("goroutine: app.ReadInput: %s", operation)

			_, ok := c.Commands[operation]

			switch {
			case operation == "create_pvz":
				communicate.createPVZArgs <- args
			case operation == "get_pvz":
				communicate.getPVZArgs <- args
			default:
				if ok || operation == helpConst {
					//var msg *string
					msg, err := c.OrderRun(args)
					if err != nil {
						communicate.errMessage <- fmt.Errorf("app.ReadInput: c.OrderRun: %v: %w", operation, err)
					} else {
						communicate.logMessage <- fmt.Sprintf("%v\n", *msg)
					}
				} else {
					communicate.errMessage <- fmt.Errorf("app.ReadInput: unknown command")
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func input(communicate *pvzChan, inputChan chan string, ticker *time.Ticker) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		<-ticker.C
		if scanner.Scan() {
			input := scanner.Text()
			select {
			case inputChan <- input:
			case <-communicate.ctx.Done():
				close(inputChan)
				return
			}
		} else if err := scanner.Err(); err != nil {
			communicate.errMessage <- fmt.Errorf("scanner.Err: %v", err)
			close(inputChan)
			return
		} else {
			close(inputChan)
			return
		}
	}
}

func (c *CLI) createPVZ(communicate *pvzChan) {

	cmdCreatePvz := c.Commands["create_pvz"]
	for args := range communicate.createPVZArgs {
		communicate.logMessage <- fmt.Sprintf("goroutine: app.CreatePvz: %s", strings.Join(args[1:], " "))
		callHelp := len(args) > 2 && strings.HasSuffix(args[2], "help")

		argument, err := cmdCreatePvz.CheckFormat(args[2:], false)
		if err != nil {
			communicate.errMessage <- fmt.Errorf("goroutine: app.CreatePvz: %w", err)
			continue
		}

		if callHelp {
			help := cmdCreatePvz.Help()
			communicate.logMessage <- help
			continue
		}

		if argument == nil {
			communicate.logMessage <- "goroutine: app.CreatePvz: help called"
		} else {
			var msg *string
			msg, err = cmdCreatePvz.Do(argument)
			if err != nil {
				communicate.errMessage <- fmt.Errorf("goroutine: app.CreatePvz: %w", err)
				continue
			}
			communicate.logMessage <- fmt.Sprintf("goroutine: app.CreatePvz: %v: %s", *msg, strings.Join(args[1:], " "))
		}
	}
}

func (c *CLI) getPVZ(communicate *pvzChan) {
	cmdGetPvz := c.Commands["get_pvz"]
	for args := range communicate.getPVZArgs {
		communicate.logMessage <- fmt.Sprintf("goroutine: app.GetPvz: %s", strings.Join(args[1:], " "))
		callHelp := len(args) > 2 && strings.HasSuffix(args[2], "help")

		argument, err := cmdGetPvz.CheckFormat(args[2:], false)
		if err != nil {
			communicate.errMessage <- fmt.Errorf("goroutine: app.GetPvz: %w", err)
			continue
		}

		if callHelp {
			help := cmdGetPvz.Help()
			communicate.logMessage <- help
			continue
		}

		if argument == nil {
			communicate.logMessage <- "goroutine: app.GetPvz: help called"
		} else {
			var msg *string
			msg, err = cmdGetPvz.Do(argument)
			if err != nil {
				communicate.errMessage <- fmt.Errorf("goroutine: app.GetPvz: %w", err)
				continue
			}
			communicate.logMessage <- fmt.Sprintf("goroutine: app.GetPvz: %s\nPVZ info:%v", strings.Join(args, " "), *msg)
		}
	}
}
