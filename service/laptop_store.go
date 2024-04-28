package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/moataz-hamed/pb/pb"
)

var ErrAlreadyExists = errors.New("Record already exists")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

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

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("Can not copy laptop data:%w", err)
	}

	store.data[laptop.Id] = other

	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, err
	}
	return other, nil
}

func (store *InMemoryLaptopStore) Search(filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		if isQualified(filter, laptop) {
			other := &pb.Laptop{}

			err := copier.Copy(other, laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) int {
	value := int(memory.GetValue())
	unit := memory.GetUnit()

	switch unit {
	case pb.Memory_BIT:
		return int(value)
	case pb.Memory_BYTE:
		return value * 8 //value << 3 (since 8 = 2^3)
	case pb.Memory_KILOBYTE:
		return value * 8 * 1024 // value << 13
	case pb.Memory_MEGABYTE:
		return value * 8 * 1024 * 1024 // value << 23
	case pb.Memory_GIGABYTE:
		return value * 8 * 1024 * 1024 * 1024 // value << 33
	case pb.Memory_TERABYTE:
		return value * 8 * 1024 * 1024 * 1024 * 1024 // value << 43

	default:
		return 0
	}

}
