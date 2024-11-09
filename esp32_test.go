package main

import (
	"machine"
	"time"
)

// TestMode constants
const (
	TEST_I2C = iota
	TEST_SINGLE_SERVO
	TEST_FINGER
	TEST_ALL_FINGERS
	TEST_SEQUENCE
)

// Test configuration
const (
	TEST_MODE        = TEST_I2C  // Change this to run different tests
	DEBUG_LED_PIN    = machine.LED
	TEST_SERVO_CHANNEL = 0       // For single servo test
	TEST_FINGER_INDEX = 0        // For finger test (0-4)
)

// Define servo positions for testing
const (
	SERVO_MIN = 250
	SERVO_MID = 325
	SERVO_MAX = 400
)

func runI2CTest() {
	println("Starting I2C test...")
	
	// Test reading MODE1 register from both PCA9685 boards
	mode1_board1 := readByte(PCA9685_ADDR1, PCA9685_MODE1)
	mode1_board2 := readByte(PCA9685_ADDR2, PCA9685_MODE1)
	
	println("Board 1 MODE1:", mode1_board1)
	println("Board 2 MODE1:", mode1_board2)
	
	// Blink LED if both boards respond
	if mode1_board1 != 0 && mode1_board2 != 0 {
		blinkLED(3, 100) // 3 quick blinks = success
	} else {
		blinkLED(1, 1000) // 1 long blink = failure
	}
}

func runSingleServoTest() {
	println("Testing single servo on channel", TEST_SERVO_CHANNEL)
	
	// Sweep from min to max
	for pos := SERVO_MIN; pos <= SERVO_MAX; pos += 5 {
		setPWM(PCA9685_ADDR1, uint8(TEST_SERVO_CHANNEL), uint16(pos))
		time.Sleep(50 * time.Millisecond)
	}
	
	// Return to middle position
	setPWM(PCA9685_ADDR1, uint8(TEST_SERVO_CHANNEL), uint16(SERVO_MID))
}

func runFingerTest(fingerIndex int) {
	println("Testing finger", fingerIndex)
	
	// Define finger servo channels
	type Finger struct {
		board                    uint8
		proximal, medial, distal uint8
	}
	
	fingers := []Finger{
		{PCA9685_ADDR2, 0, 1, 2},    // Thumb
		{PCA9685_ADDR1, 2, 1, 0},    // Index
		{PCA9685_ADDR1, 4, 5, 7},    // Middle
		{PCA9685_ADDR1, 8, 9, 10},   // Ring
		{PCA9685_ADDR1, 12, 13, 14}, // Pinky
	}
	
	if fingerIndex < 0 || fingerIndex >= len(fingers) {
		println("Invalid finger index")
		return
	}
	
	finger := fingers[fingerIndex]
	
	// Test sequence: open -> close -> open
	positions := []uint16{SERVO_MIN, SERVO_MAX, SERVO_MIN}
	
	for _, pos := range positions {
		// Move all joints
		setPWM(finger.board, finger.proximal, pos)
		setPWM(finger.board, finger.medial, pos)
		setPWM(finger.board, finger.distal, pos)
		time.Sleep(500 * time.Millisecond)
	}
}

func runAllFingersTest() {
	println("Testing all fingers")
	
	// Test each finger in sequence
	for i := 0; i < 5; i++ {
		runFingerTest(i)
		time.Sleep(1 * time.Second)
	}
}

func runSequenceTest() {
	println("Running test sequence")
	
	// Make a fist
	println("Making fist...")
	processCommand("command1," + 
		"400,400,400,400," + // Index finger
		"400,400,400,400")   // Middle finger
	processCommand("command2," + 
		"400,400,400,400," + // Ring finger
		"400,400,400,400")   // Pinky finger
	processCommand("command3,400,400,400") // Thumb
	time.Sleep(2 * time.Second)
	
	// Open hand
	println("Opening hand...")
	processCommand("command1," + 
		"250,250,250,250," + // Index finger
		"250,250,250,250")   // Middle finger
	processCommand("command2," + 
		"250,250,250,250," + // Ring finger
		"250,250,250,250")   // Pinky finger
	processCommand("command3,250,250,250") // Thumb
	time.Sleep(2 * time.Second)
}

func blinkLED(times int, duration time.Duration) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	for i := 0; i < times; i++ {
		led.High()
		time.Sleep(duration)
		led.Low()
		time.Sleep(duration)
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
	
	println("Starting test mode:", TEST_MODE)
	
	// Run selected test
	switch TEST_MODE {
	case TEST_I2C:
		runI2CTest()
	case TEST_SINGLE_SERVO:
		runSingleServoTest()
	case TEST_FINGER:
		runFingerTest(TEST_FINGER_INDEX)
	case TEST_ALL_FINGERS:
		runAllFingersTest()
	case TEST_SEQUENCE:
		runSequenceTest()
	}
	
	// Main loop - keep program running
	for {
		time.Sleep(1 * time.Second)
	}
}
