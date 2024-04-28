package sample

import (
	"github.com/moataz-hamed/pb/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewKeyboard returns a new sample keyboard
func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}

func NewCpu() *pb.CPU {
	brand := randomCpuBrand()
	name := randomCpuName(brand)
	coreNumber := randomInt(2, 8)
	threadsNumber := randomInt(coreNumber, 12)
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)
	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(coreNumber),
		NumberThreads: uint32(threadsNumber),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
	return cpu
}

func NewGpu() *pb.GPU {
	brand := randomGpuBrand()
	name := randomGpuName(brand)
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)
	mem := &pb.Memory{
		Value: uint32(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: mem,
	}
	return gpu
}

func NewRam() *pb.Memory {
	ram := &pb.Memory{
		Value: uint32(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	return ram
}

func NewSSD() *pb.Storage {
	mem := &pb.Memory{
		Value: uint32(randomInt(128, 1024)),
		Unit:  pb.Memory_GIGABYTE,
	}
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: mem,
	}
	return ssd
}

func NewHDD() *pb.Storage {
	mem := &pb.Memory{
		Value: uint32(randomInt(1, 6)),
		Unit:  pb.Memory_TERABYTE,
	}
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: mem,
	}
	return hdd
}

func NewScreen() *pb.Screen {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	screen := &pb.Screen{
		SizeInch: float32(randomFloat64(13, 17)),
		Resolution: &pb.Screen_Resolution{
			Width:  uint32(width),
			Height: uint32(height),
		},
		Panel:      randomScreenPanel(),
		Nultitouch: randomBool(),
	}
	return screen
}

func NewLaptop() *pb.Laptop {
	brand := randLaptopBrand()
	name := randLaptopName(brand)
	laptop := &pb.Laptop{
		Id:       randID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCpu(),
		Ram:      NewRam(),
		Gpu:      []*pb.GPU{NewGpu()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(5.0, 10.0),
		},
		PriceUsd:    randomFloat64(500, 2000),
		ReleaseYear: uint32(randomInt(2000, 2024)),
		UpdatedAt:   timestamppb.Now(),
	}
	return laptop
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
