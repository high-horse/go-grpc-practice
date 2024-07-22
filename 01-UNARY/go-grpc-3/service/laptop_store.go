package service

import (
	"errors"
	"fmt"
	"grpc-3/pb"
	"sync"

	"github.com/jinzhu/copier"
)

// Define error
var ErrAlreadyExists = errors.New("record already exists")

// Interface to LaptopStore
type LaptopStore interface {
	Save (laptop *pb.Laptop) error
}

// InMemoryLaptopStore stores laptops in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

// NewInMemoryLaptopStore creates new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}


func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %v", err)
	}

	store.data[other.Id] = other
	return nil
}