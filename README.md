# gogley - 3D Printed Biomimetic Mechatronic Hand Go Controller
TinyGo version of Will's [Bionic Hand Delta 1.1](https://willcogley.notion.site/Will-Cogley-Project-Archive-75a4864d73ab4361ab26cabaadaec33a?p=3e7718a58fc34e5ab0736f6c523bee1e&pm=c) port of the Arduino code with few changes:

## Direct I2C Communication:

- Uses TinyGo's machine package for I2C communication
- Implements direct register access on the PCA9685 boards
- Maintains I2C addresses (0x40 and 0x41)

## Serial Communication:

- Replaces Uduino with serial communication
- Expects commands in csv format: "command1,250,300,350,400,250,300,350,400"

## PWM Control:

- Maintains same servo frequency of 50Hz
- Uses same oscillator frequency of 27MHz
- Direct PWM control through register instructions
