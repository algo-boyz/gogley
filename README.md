# gogley - 3D Printed Biomimetic Mechatronic Hand Go Controller
TinyGo version maintains the same functionality as the Arduino code but with several key changes:
Sorta mirror of [Bionic Hand Delta 1.1](https://willcogley.notion.site/Will-Cogley-Project-Archive-75a4864d73ab4361ab26cabaadaec33a?p=3e7718a58fc34e5ab0736f6c523bee1e&pm=c)

## Direct I2C Communication:

Uses TinyGo's machine package for I2C communication
Implements direct register access to the PCA9685 boards
Maintains the same I2C addresses (0x40 and 0x41)


## Serial Communication:

Replaces Uduino with simple serial communication
Expects commands in a comma-separated format
Example command format: "command1,250,300,350,400,250,300,350,400"


## PWM Control:

Maintains the same servo frequency (50Hz)
Uses the same oscillator frequency (27MHz)
Implements direct PWM control through register manipulation
