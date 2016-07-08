package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*Tsl2591Driver)(nil)

// TSL2591 ALS definitions
const (
	TSL2591_ADDRESS    = 0x29
	TSL2591_COMMAND    = 0x80
	TSL2591_NORMAL_OP  = 0x20
	TSL2591_ENABLE_RW  = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x00
	TSL2591_CONFIG_RW  = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x01
	TSL2591_C0_DATA_LR = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x14
	TSL2591_C1_DATA_LR = TSL2591_COMMAND | TSL2591_NORMAL_OP | 0x16
)

type Tsl2591Driver struct {
	name       string
	connection I2cExtended
	gobot.Commander
	initialized bool
}

func NewTsl2591Driver(i I2cExtended, name string) *Tsl2591Driver {
	b := &Tsl2591Driver{
		name:        name,
		connection:  i,
		Commander:   gobot.NewCommander(),
		initialized: false,
	}

	b.AddCommand("LuxInput", func(params map[string]interface{}) interface{} {
		ambIR, iR, err := b.ReadSensor()
		return map[string]interface{}{"AmbientIR": fmt.Sprintf("%X", ambIR), "IR": fmt.Sprintf("%X", iR), "err": err}
	})

	return b
}

func (b *Tsl2591Driver) Name() string                 { return b.name }
func (b *Tsl2591Driver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start writes start bytes
func (b *Tsl2591Driver) Start() (errs []error) {
	if !b.initialized {
		b.connection.I2cWriteWord(TSL2591_ADDRESS, TSL2591_ENABLE_RW, 0x03)
		// TSL2591_CONFIG_RW = 0x11: Medium gain mode, 200ms integration time
		b.connection.I2cWriteWord(TSL2591_ADDRESS, TSL2591_CONFIG_RW, 0x11)

		b.initialized = true
	}
	return
}

// Halt returns true if device is halted successfully
func (b *Tsl2591Driver) Halt() (errs []error) { return }

// Digital input
func (b *Tsl2591Driver) ReadSensor() (local []byte, object []byte, errs []error) {
	ambIR, errAmbIR := b.connection.I2cReadRegister([]byte{TSL2591_ADDRESS, TSL2591_C0_DATA_LR}, 2)
	fmt.Printf("ambIR %X \n", ambIR)

	iR, errIR := b.connection.I2cReadRegister([]byte{TSL2591_ADDRESS, TSL2591_C1_DATA_LR}, 2)
	fmt.Printf("ambIR %X \n", iR)

	return ambIR, iR, []error{errAmbIR, errIR}
}

/*

// Read TSL2591 ALS
function readAls() {
    var dataCh0 = i2c.readWordSync(TSL2591_ADDR, TSL2591_C0DATAL_R);
    console.log("Ambient light and IR");
    console.log(dataCh0);
    //localTemp = (swap16(localTemp))/128.0;

    var dataCh1 = i2c.readWordSync(TSL2591_ADDR, TSL2591_C1DATAL_R);
    console.log("IR only");
    console.log(dataCh1);
    //localTemp = (swap16(localTemp))/128.0;

    // CPL = Count per lux = integration time / gain
    cpl = 200/25.0;
    // Calculate lux
    lux = (dataCh0 - 2*dataCh1)/cpl;
    console.log("The detected ambient light is "+ lux + " lx");
}

*/
