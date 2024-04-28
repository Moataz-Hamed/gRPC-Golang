package sample

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/moataz-hamed/pb/pb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomBool() bool {
	return rand.Intn(2) == 1 // returns 0 or 1
}

func randomCpuBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	return a[rand.Intn(len(a))]
}

func randomCpuName(s string) string {
	if s == "intel" {
		return randomStringFromSet(
			"Core i3",
			"Core i5",
			"Core i7",
			"Core i9",
		)
	}

	return randomStringFromSet(
		"Ryzen 3",
		"Ryzen 5",
		"Ryzen 7",
	)
}

func randomInt(a, b int) int {
	n := a + rand.Intn(b-a+1)
	return n
}

func randomFloat64(a, b float64) float64 {
	return a + rand.Float64()*(b-a)
}

func randomGpuBrand(a ...string) string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGpuName(a string) string {
	if a == "NVIDIA" {
		return "RTX"
	}
	return "RX VEGA"
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randID() string {
	return uuid.New().String()
}

func randLaptopBrand() string {
	return randomStringFromSet("DELL", "MAC", "LENOVO")
}

func randLaptopName(a string) string {
	switch a {
	case "DELL":
		return randomStringFromSet("G series", "Inspiron", "AlienWare", "latitude")
	case "MAC":
		return randomStringFromSet("Air", "Pro")
	default:
		return randomStringFromSet("Thinkpad", "Ideapad")
	}
}
