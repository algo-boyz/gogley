package main

import (
	"machine"
	"time"
)

// Test modes
const (
	MODE_POSES = iota
	MODE_JOINT_LIMITS
)

// Configuration
const (
	TEST_MODE = MODE_POSES  // Change between MODE_POSES and MODE_JOINT_LIMITS
	DELAY_MS  = 500        // Delay between movements
)

// Preset positions for all servos
type HandPosition struct {
	name string
	positions map[string]uint16
}

// Create known good positions
var (
	extendedPosition = HandPosition{
		name: "extended",
		positions: map[string]uint16{
			"02_dist": 280, "02_med": 300, "02_prox": 300, "02_lat": 305,
			"03_prox": 290, "03_med": 300, "03_lat": 310, "03_dist": 290,
			"04_prox": 290, "04_med": 300, "04_dist": 280, "04_lat": 305,
			"05_prox": 290, "05_med": 290, "05_dist": 280, "05_lat": 310,
		},
	}
	
	spreadPosition = HandPosition{
		name: "spread",
		positions: map[string]uint16{
			"02_dist": 280, "02_med": 300, "02_prox": 300, "02_lat": 325,
			"03_prox": 290, "03_med": 300, "03_lat": 320, "03_dist": 290,
			"04_prox": 290, "04_med": 300, "04_dist": 280, "04_lat": 280,
			"05_prox": 290, "05_med": 290, "05_dist": 280, "05_lat": 250,
		},
	}
	
	fistPosition = HandPosition{
		name: "fist",
		positions: map[string]uint16{
			"02_dist": 380, "02_med": 180, "02_prox": 370, "02_lat": 305,
			"03_prox": 200, "03_med": 380, "03_lat": 310, "03_dist": 220,
			"04_prox": 200, "04_med": 190, "04_dist": 380, "04_lat": 305,
			"05_prox": 380, "05_med": 200, "05_dist": 380, "05_lat": 310,
		},
	}
)

// Servo mapping to boards and channels
type ServoMap struct {
	board   uint8
	channel uint8
}

var servoMapping = map[string]ServoMap{
	// Index finger (02)
	"02_dist":  {PCA9685_ADDR1, 0},
	"02_med":   {PCA9685_ADDR1, 1},
	"02_prox":  {PCA9685_ADDR1, 2},
	"02_lat":   {PCA9685_ADDR2, 3},
	
	// Middle finger (03)
	"03_prox":  {PCA9685_ADDR1, 4},
	"03_med":   {PCA9685_ADDR1, 5},
	"03_dist":  {PCA9685_ADDR1, 7},
	"03_lat":   {PCA9685_ADDR2, 4},
	
	// Ring finger (04)
	"04_prox":  {PCA9685_ADDR1, 8},
	"04_med":   {PCA9685_ADDR1, 9},
	"04_dist":  {PCA9685_ADDR1, 10},
	"04_lat":   {PCA9685_ADDR2, 7},
	
	// Pinky (05)
	"05_prox":  {PCA9685_ADDR1, 12},
	"05_med":   {PCA9685_ADDR1, 13},
	"05_dist":  {PCA9685_ADDR1, 14},
	"05_lat":   {PCA9685_ADDR2, 8},
}

func applyPosition(pos HandPosition) {
	println("Applying position:", pos.name)
	
	for servoName, value := range pos.positions {
		if mapping, exists := servoMapping[servoName]; exists {
			setPWM(mapping.board, mapping.channel, value)
		}
	}
}

func runPosesTest() {
	println("Running poses test cycle")
	
	for {
		// Extended position
		println("Moving to extended position")
		applyPosition(extendedPosition)
		// Add thumb position for extended
		setPWM(PCA9685_ADDR2, 0, 270) // thumb prox
		setPWM(PCA9685_ADDR2, 1, 270) // thumb dist
		setPWM(PCA9685_ADDR2, 2, 400) // thumb lat
		time.Sleep(time.Duration(DELAY_MS) * time.Millisecond)
		
		// Fist position
		println("Moving to fist position")
		applyPosition(fistPosition)
		// Add thumb position for fist
		setPWM(PCA9685_ADDR2, 2, 260) // thumb lat
		time.Sleep(time.Duration(DELAY_MS) * time.Millisecond)
	}
}

func runJointLimitsTest() {
	println("Running joint limits test")
	
	// Test lateral movements
	lateralServos := []struct {
		name    string
		board   uint8
		channel uint8
	}{
		{"02_lat", PCA9685_ADDR1, 3},
		{"03_lat", PCA9685_ADDR1, 6},
		{"04_lat", PCA9685_ADDR1, 11},
		{"05_lat", PCA9685_ADDR1, 15},
	}
	
	for {
		// Test extension limits
		println("Testing extension limits")
		for _, servo := range lateralServos {
			println("Testing", servo.name, "extension")
			setPWM(servo.board, servo.channel, 325)
		}
		time.Sleep(time.Duration(DELAY_MS) * time.Millisecond)
		
		// Test flexion limits
		println("Testing flexion limits")
		for _, servo := range lateralServos {
			println("Testing", servo.name, "flexion")
			setPWM(servo.board, servo.channel, 250)
		}
		time.Sleep(time.Duration(DELAY_MS) * time.Millisecond)
	}
}

func main() {
	// Initialize I2C
	i2c = machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 400000,
		SDA:       machine.I2C0_SDA_PIN,
		SCL:       machine.I2C0_SCL_PIN,
	})

	// Initialize both PCA9685 boards
	initPCA9685(PCA9685_ADDR1)
	initPCA9685(PCA9685_ADDR2)
	
	println("Starting servo limits test in mode:", TEST_MODE)
	
	switch TEST_MODE {
	case MODE_POSES:
		runPosesTest()
	case MODE_JOINT_LIMITS:
		runJointLimitsTest()
	}
}
