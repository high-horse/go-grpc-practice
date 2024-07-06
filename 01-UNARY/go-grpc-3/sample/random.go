package sample

import (
	"grpc-3/pb"
	"math/rand"

	"github.com/google/uuid"
)

func randomBool() bool{
	return rand.Intn(2) == 1 
}

func randomKeyboardLayout() pb.Keyboard_Layout{
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTZ
	case 2:
		return pb.Keyboard_AZERTY
	default:
		return pb.Keyboard_QWERTY
	}
}

func randomCPUBrand() string{
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomCPUName(brand string) string{
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-8700K",
			"Core i5-8265U",
			"Core i3-1005G1",
		)
	}
	return randomStringFromSet(
		"Ryzen 7 2700",
		"Ryzen 5 3600",
		"Ryzen 3 3200G",
	)
}


func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")	
}

func randomGPUName(brand string) string{
	if brand == "Nvidia" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX 1070",
		)
	}
	return randomStringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700",
		"RX Vega 64",
	)
}

func randomPanel() pb.Screen_Panel{
	if randomBool() {
		return pb.Screen_OLED
	}
	return pb.Screen_IPS
}

func randomScreenResolution() *pb.Screen_Resolution{
	height := randomInt(1080, 4320)
	width := height * 16 / 9;

	return &pb.Screen_Resolution{
		Width: uint32(height),
		Height: uint32(width),
	}
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":	
		return randomStringFromSet("Latitude", "XPS", "Vostro")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1")
	}
}

func randomInt(min int, max int) int{
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomID() string{
	return uuid.New().String()
}