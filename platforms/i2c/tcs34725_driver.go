package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*Tcs34725Driver)(nil)

// TCS34725 RGB sensor definitions
const (
	TCS34725_ADDR = 0x29

	TCS34725_COMMAND   = 0x80
	TCS34725_TYPE_AUTO = 0x20

	TCS34725_ENABLE_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x00
	TCS34725_ATIME_RW  = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x01
	TCS34725_WTIME_RW  = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x03

	TCS34725_AILTL_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x04
	TCS34725_AILTH_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x05
	TCS34725_AIHTL_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x06
	TCS34725_AIHTH_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x07

	TCS34725_CONFIG_RW  = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x0D
	TCS34725_CONTROL_RW = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x0F

	TCS34725_ID_R     = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x12
	TCS34725_STATUS_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x13

	TCS34725_CDATAL_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x14
	TCS34725_CDATAH_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x15
	TCS34725_RDATAL_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x16
	TCS34725_RDATAH_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x17
	TCS34725_GDATAL_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x18
	TCS34725_GDATAH_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x19
	TCS34725_BDATAL_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x1A
	TCS34725_BDATAH_R = TCS34725_COMMAND | TCS34725_TYPE_AUTO | 0x1B

	TSL2591_ADDRESS    = 0x29
	TSL2591_COMMAND    = 0x80
	TSL2591_NORMAL_OP  = 0x20
	TSL2591_ENABLE_RW  = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x00
	TSL2591_CONFIG_RW  = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x01
	TSL2591_C0_DATA_LR = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x14
	TSL2591_C1_DATA_LR = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x16
)

type Tcs34725Driver struct {
	name       string
	connection I2cExtended
	gobot.Commander
	initialized bool
}

func NewTcs34725Driver(i I2cExtended, name string) *Tcs34725Driver {
	b := &Tcs34725Driver{
		name:        name,
		connection:  i,
		Commander:   gobot.NewCommander(),
		initialized: false,
	}

	b.AddCommand("ColorSensorInput", func(params map[string]interface{}) interface{} {
		redData, greenData, blueData, err := b.ReadSensor()
		return map[string]interface{}{"Red": fmt.Sprintf("%X", redData), "Green": fmt.Sprintf("%X", greenData), "Blue": fmt.Sprintf("%X", ColorSensorInput), "err": err}
	})

	return b
}

func (b *Tcs34725Driver) Name() string                 { return b.name }
func (b *Tcs34725Driver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start writes start bytes
func (b *Tcs34725Driver) Start() (errs []error) {
	if !b.initialized {
		// Enable sensor
		b.connection.I2cWriteWord(TCS34725_ADDR, TCS34725_ENABLE_RW, 0x03)
		// Set integration time to be 700 ms
		b.connection.I2cWriteWord(TCS34725_ADDR, TCS34725_ATIME_RW, 0x00)
		// Set the gain to be 1x
		b.connection.I2cWriteWord(TCS34725_ADDR, TCS34725_CONTROL_RW, 0x00)

		b.initialized = true
	}
	return
}

// Halt returns true if device is halted successfully
func (b *Tcs34725Driver) Halt() (errs []error) { return }

// Digital input
func (b *Tcs34725Driver) ReadSensor() (local []byte, object []byte, errs []error) {

	// Read red data
	redData, errRed := b.connection.I2cReadRegister([]byte{TCS34725_ADDR, TCS34725_RDATAL_R}, 2)
	fmt.Printf("redData %X \n", ambIR)

	// Read green data
	greenData, errGreen := b.connection.I2cReadRegister([]byte{TCS34725_ADDR, TCS34725_RDATAL_R}, 2)
	fmt.Printf("greenData %X \n", iR)

	// Read blue data
	blueData, errBlue := b.connection.I2cReadRegister([]byte{TCS34725_ADDR, TCS34725_RDATAL_R}, 2)
	fmt.Printf("blueData %X \n", ambIR)

	return redData, greenData, blueData, []error{errRed, errGreen, errBlue}
}

/*

def init():
    # Enable sensor
    bus.write_byte_data(TCS34725_ADDR, TCS34725_ENABLE_RW, 0x03)
    # Set integration time to be 700 ms
    bus.write_byte_data(TCS34725_ADDR, TCS34725_ATIME_RW, 0x00)
    # Set the gain to be 1x
    bus.write_byte_data(TCS34725_ADDR, TCS34725_CONTROL_RW, 0x00)

def read_lux():
    # Read red data
    r = bus.read_word_data(TCS34725_ADDR, TCS34725_RDATAL_R)
    # Read green data
    g = bus.read_word_data(TCS34725_ADDR, TCS34725_GDATAL_R)
    # Read blue data
    b = bus.read_word_data(TCS34725_ADDR, TCS34725_BDATAL_R)
    # Calculated lux value based on RGB data
    lux = (-0.32466 * r) + (1.57837 * g) + (-0.73191 * b)

    # print "red is: " + r + "\nblue is : " + b + "\ngreen is: " + g + "\nlux is: " + lux

    if lux < 0:
        return 0
    else:
        return lux

*/
