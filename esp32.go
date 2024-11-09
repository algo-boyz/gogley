package main

import (
	"machine"
	"strconv"
	"strings"
	"time"
)

// PCA9685 registers
const (
	PCA9685_MODE1      = 0x00
	PCA9685_PRESCALE   = 0xFE
	PCA9685_LED0_ON_L  = 0x06
	MODE1_RESTART      = 0x80
	MODE1_AI           = 0x20
	MODE1_SLEEP        = 0x10
	SERVO_FREQ         = 50
	OSCILLATOR_FREQ    = 27000000
)

// Define I2C addresses for the two PCA9685 boards
const (
	PCA9685_ADDR1 = 0x40
	PCA9685_ADDR2 = 0x41
)

var (
	i2c    machine.I2C
	uart   machine.UART
	buffer = make([]byte, 64)
)

// initPCA9685 initializes a PCA9685 board at the given address
func initPCA9685(addr uint8) {
	// Reset
	writeByte(addr, PCA9685_MODE1, MODE1_RESTART)
	time.Sleep(10 * time.Millisecond)

	// Set to sleep mode
	mode1 := readByte(addr, PCA9685_MODE1)
	mode1 = (mode1 & ^uint8(MODE1_RESTART)) | MODE1_SLEEP
	writeByte(addr, PCA9685_MODE1, mode1)

	// Set PWM frequency
	prescale := uint8(OSCILLATOR_FREQ/(4096*SERVO_FREQ) - 1)
	writeByte(addr, PCA9685_PRESCALE, prescale)

	// Wake up
	mode1 = readByte(addr, PCA9685_MODE1)
	mode1 = mode1&^MODE1_SLEEP | MODE1_AI
	writeByte(addr, PCA9685_MODE1, mode1)
	time.Sleep(5 * time.Millisecond)
}

// setPWM sets the PWM value for a channel
func setPWM(addr uint8, channel uint8, value uint16) {
	reg := PCA9685_LED0_ON_L + (channel * 4)
	writeByte(addr, reg, 0)
	writeByte(addr, reg+1, 0)
	writeByte(addr, reg+2, uint8(value))
	writeByte(addr, reg+3, uint8(value>>8))
}

func writeByte(addr uint8, reg uint8, value uint8) {
	buf := []byte{reg, value}
	i2c.Tx(addr, buf, nil)
}

func readByte(addr uint8, reg uint8) uint8 {
	buf := make([]byte, 1)
	i2c.Tx(addr, []byte{reg}, buf)
	return buf[0]
}

func processCommand(cmd string) {
	parts := strings.Split(cmd, ",")
	if len(parts) < 3 {
		return
	}

	switch parts[0] {
	case "command1":
		if len(parts) >= 9 {
			// Parse values for first set of servos
			prox02, _ := strconv.Atoi(parts[1])
			med02, _ := strconv.Atoi(parts[2])
			dist02, _ := strconv.Atoi(parts[3])
			lat02, _ := strconv.Atoi(parts[4])
			prox03, _ := strconv.Atoi(parts[5])
			med03, _ := strconv.Atoi(parts[6])
			dist03, _ := strconv.Atoi(parts[7])
			lat03, _ := strconv.Atoi(parts[8])

			// Set PWM values for first PCA9685
			setPWM(PCA9685_ADDR1, 2, uint16(prox02))
			setPWM(PCA9685_ADDR1, 1, uint16(med02))
			setPWM(PCA9685_ADDR1, 0, uint16(dist02))
			setPWM(PCA9685_ADDR2, 3, uint16(lat02))

			setPWM(PCA9685_ADDR1, 4, uint16(prox03))
			setPWM(PCA9685_ADDR1, 5, uint16(med03))
			setPWM(PCA9685_ADDR1, 7, uint16(dist03))
			setPWM(PCA9685_ADDR2, 4, uint16(lat03))
		}
	case "command2":
		if len(parts) >= 9 {
			// Parse values for second set of servos
			prox04, _ := strconv.Atoi(parts[1])
			med04, _ := strconv.Atoi(parts[2])
			dist04, _ := strconv.Atoi(parts[3])
			lat04, _ := strconv.Atoi(parts[4])
			prox05, _ := strconv.Atoi(parts[5])
			med05, _ := strconv.Atoi(parts[6])
			dist05, _ := strconv.Atoi(parts[7])
			lat05, _ := strconv.Atoi(parts[8])

			// Set PWM values for second set
			setPWM(PCA9685_ADDR1, 8, uint16(prox04))
			setPWM(PCA9685_ADDR1, 9, uint16(med04))
			setPWM(PCA9685_ADDR1, 10, uint16(dist04))
			setPWM(PCA9685_ADDR2, 7, uint16(lat04))

			setPWM(PCA9685_ADDR1, 12, uint16(prox05))
			setPWM(PCA9685_ADDR1, 13, uint16(med05))
			setPWM(PCA9685_ADDR1, 14, uint16(dist05))
			setPWM(PCA9685_ADDR2, 8, uint16(lat05))
		}
	case "command3":
		if len(parts) >= 4 {
			// Parse values for thumb servos
			prox01, _ := strconv.Atoi(parts[1])
			dist01, _ := strconv.Atoi(parts[2])
			lat01, _ := strconv.Atoi(parts[3])

			// Set PWM values for thumb
			setPWM(PCA9685_ADDR2, 0, uint16(prox01))
			setPWM(PCA9685_ADDR2, 1, uint16(dist01))
			setPWM(PCA9685_ADDR2, 2, uint16(lat01))
		}
	}
}

func main() {
	// Initialize I2C
	i2c = machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 400000, // 400kHz
		SDA:       machine.I2C0_SDA_PIN,
		SCL:       machine.I2C0_SCL_PIN,
	})

	// Initialize UART
	uart = machine.UART0
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.UART0_TX_PIN,
		RX:       machine.UART0_RX_PIN,
	})

	// Initialize both PCA9685 boards
	initPCA9685(PCA9685_ADDR1)
	initPCA9685(PCA9685_ADDR2)

	var cmdBuffer string
	for {
		if uart.Buffered() > 0 {
			byte, _ := uart.ReadByte()
			
			if byte == '\n' {
				processCommand(cmdBuffer)
				cmdBuffer = ""
			} else {
				cmdBuffer += string(byte)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
