package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"Homework-1/internal/model/cli"
)

// Storage is
type Storage struct {
	store     *os.File
	storeName string
	data      map[cli.OrderID]cli.Order
	mtx       sync.RWMutex
}

// New is
func New(storeName string) (*Storage, error) {
	// #nosec G304
	file, err := os.OpenFile(storeName, os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	store := &Storage{file, storeName, make(map[cli.OrderID]cli.Order), sync.RWMutex{}}

	err = store.ReadAll()
	if err != nil {
		log.Fatalf("[main] storage.ReadAll: %v", err)
	}

	return store, nil
}

// Close is
func (s *Storage) Close() error {
	if s.store == nil {
		return nil
	}

	if err := s.Save(); err != nil {
		return fmt.Errorf("storage.Save: %w", err)
	}

	return s.store.Close()
}

// Save is
func (s *Storage) Save() error {
	rawBytes, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = os.WriteFile(s.storeName, rawBytes, 0600)
	if err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}
	return nil
}

// Create is
func (s *Storage) Create(newOrder *cli.Order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.data[newOrder.ID]
	if ok {
		return errors.New("order with provided ID exists")
	}
	s.data[newOrder.ID] = cli.Order{
		ID:         newOrder.ID,
		UserID:     newOrder.UserID,
		ExpireAt:   newOrder.ExpireAt,
		IsDeleted:  newOrder.IsDeleted,
		IsReturned: newOrder.IsReturned,
		IsIssued:   newOrder.IsIssued,
		IsAccepted: newOrder.IsAccepted,
		ReceivedAt: newOrder.ReceivedAt,
		IssuedAt:   newOrder.IssuedAt,
	}

	return nil
}

// Find is
func (s *Storage) Find(orderID cli.OrderID) (cli.Order, error) {
	s.mtx.RLock()
	value, ok := s.data[orderID]
	s.mtx.RUnlock()
	if !ok {
		return cli.Order{}, errors.New("order was not found")
	}
	return value, nil
}

// Update is
func (s *Storage) Update(updatedOrder cli.Order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.data[updatedOrder.ID]
	if !ok {
		return errors.New("order with provided ID doesn't exist")
	}
	s.data[updatedOrder.ID] = cli.Order{
		ID:         updatedOrder.ID,
		UserID:     updatedOrder.UserID,
		ExpireAt:   updatedOrder.ExpireAt,
		IsDeleted:  updatedOrder.IsDeleted,
		IsReturned: updatedOrder.IsReturned,
		IsIssued:   updatedOrder.IsIssued,
		IsAccepted: updatedOrder.IsAccepted,
		ReceivedAt: updatedOrder.ReceivedAt,
		IssuedAt:   updatedOrder.IssuedAt,
	}
	return nil
}

// IssueAll is
func (s *Storage) IssueAll(orderIDs []cli.OrderID) error {
	for _, val := range orderIDs {
		s.mtx.Lock()
		orderValue, ok := s.data[val]
		if !ok {
			s.mtx.Unlock()
			continue
		}
		orderValue.IsIssued = true
		orderValue.IsDeleted = true
		orderValue.IssuedAt = time.Now()
		s.data[val] = orderValue
		s.mtx.Unlock()
	}
	return nil
}

// List is
func (s *Storage) List() ([]cli.Order, error) {
	activeOrders := make([]cli.Order, 0, len(s.data))
	s.mtx.RLock()
	for _, val := range s.data {
		if !val.IsAccepted && !val.IsIssued && !val.IsReturned && !val.IsDeleted && time.Now().Before(val.ExpireAt) {
			activeOrders = append(activeOrders, val)
		}
	}
	s.mtx.RUnlock()
	return activeOrders, nil
}

// ReadAll is
func (s *Storage) ReadAll() error {
	_, err := s.store.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("os.File.Seek: %w", err)
	}

	reader := bufio.NewReader(s.store)

	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	if len(rawBytes) > 0 {
		err = json.Unmarshal(rawBytes, &s.data)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	return nil
}

// ListReturns is
func (s *Storage) ListReturns(pageSize int) []cli.Order {
	listOfReturns := make([]cli.Order, 0, pageSize)
	s.mtx.RLock()
	for _, val := range s.data {
		if val.IsAccepted {
			listOfReturns = append(listOfReturns, val)
		}
	}
	s.mtx.RUnlock()
	return listOfReturns
}
