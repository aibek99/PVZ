package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"Homework-1/internal/model/cli"
	"Homework-1/pkg/errlst"
)

// Storage is
type Storage struct {
	store     *os.File
	storeName string
	data      map[string]cli.PVZ
	mtx       sync.RWMutex
}

// New is
func New(storeName string) (*Storage, error) {
	storeName = filepath.Clean(storeName)
	file, err := os.OpenFile(storeName, os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	store := &Storage{file, storeName, make(map[string]cli.PVZ), sync.RWMutex{}}

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
func (s *Storage) Create(value *cli.PVZ) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.data[value.Name]
	if ok {
		return errlst.ErrPVZAlreadyExists
	}
	s.data[value.Name] = cli.PVZ{
		Name:    value.Name,
		Address: value.Address,
		Contact: value.Contact,
	}
	return nil
}

// Find is
func (s *Storage) Find(key string) (cli.PVZ, error) {
	s.mtx.RLock()
	value, ok := s.data[key]
	s.mtx.RUnlock()
	if !ok {
		return cli.PVZ{}, errlst.ErrPVZNotFound
	}
	return value, nil
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
