package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*Tsl2591Driver)(nil)

// TSL2591 ALS definitions
const tsl2591Address = 0x29
const tsl2591Command = 0x80
const tsl2591NormalOp = 0x20
const tsl2591EnableRW = tsl2591Command | tsl2591NormalOp | 0x00
const tsl2591ConfigRW = tsl2591Command | tsl2591NormalOp | 0x01
const tsl2591C0DataLR = tsl2591Command | tsl2591NormalOp | 0x14
const tsl2591C1DataLR = tsl2591Command | tsl2591NormalOp | 0x16

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

	b.AddCommand("LuxSensor", func(params map[string]interface{}) interface{} {
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
		b.connection.I2cWriteWord(tsl2591Address, tsl2591EnableRW, 0x03)
		// TSL2591_CONFIG_RW = 0x11: Medium gain mode, 200ms integration time
		b.connection.I2cWriteWord(tsl2591Address, tsl2591ConfigRW, 0x11)

		b.initialized = true
	}
	return
}

// Halt returns true if device is halted successfully
func (b *Tsl2591Driver) Halt() (errs []error) { return }

// Digital input
func (b *Tsl2591Driver) ReadSensor() (local []byte, object []byte, errs []error) {
	ambIR, errAmbIR := b.connection.I2cReadRegister([]byte{tsl2591Address, tsl2591C0DataLR}, 2)
	fmt.Printf("ambIR %X \n", ambIR)

	iR, errIR := b.connection.I2cReadRegister([]byte{tsl2591Address, tsl2591C1DataLR}, 2)
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
