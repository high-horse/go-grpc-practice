package sample

import (
	"grpc-3/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *pb.Keyboard{
	keyboard := &pb.Keyboard{
		Layout: randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

func NewCPU() *pb.CPU{	
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	minGhx := randomFloat64(2.5, 3.5)
	maxGhx := randomFloat64(minGhx, 5.0)

	cpu := &pb.CPU{
		Brand: brand,
		Name: name,
		NumberCores: uint32(randomInt(2, 8)),
		NumberThreads: uint32(randomInt(2, 8)),
		MinGhz: minGhx,
		MaxGhz: maxGhx,
	}
	return cpu
}

func NewGPU() *pb.GPU{
	brand := randomGPUBrand()

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit: pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand: brand,
		Name: randomGPUName(brand),
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
	return gpu
}

func NewRam() *pb.Memory{
	memory := &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit: pb.Memory_GIGABYTE,
	}
	return memory
}

func NewSSD() *pb.Storage{
	storage := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit: pb.Memory_GIGABYTE,
		},
	}
	return storage
}

func NewHDD() *pb.Storage{
	storage := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit: pb.Memory_TERABYTE,
		},
	}
	return storage
}

func NewScreen() *pb.Screen{

	screen := &pb.Screen{
		SizeInch: float32(randomFloat64(13, 17)),
		Resolution: randomScreenResolution(),
		Panel: randomPanel(),
		Multitouch: randomBool(),
	}
	return screen
}

func NewLaptop() *pb.Laptop{
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &pb.Laptop{
		Id: randomID(),
		Brand: brand,
		Name: name,
		Cpu: NewCPU(),
		Ram: NewRam(),
		Gpus: []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen: NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd: randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2024)),
		UpdatedAt: timestamppb.Now(),
	}
	return laptop
}