package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*Tmp007Driver)(nil)

// TMP007 thermopile definitions
const (
	TMP007_ADDRESS            = 0x40
	TMP007_LOCAL_TEMPERATURE  = 0x01
	TMP007_OBJECT_TEMPERATURE = 0x03
)

type Tmp007Driver struct {
	name       string
	connection I2cExtended
	gobot.Commander
	initialized bool
}

func NewTmp007Driver(i I2cExtended, name string) *Tmp007Driver {
	b := &Tmp007Driver{
		name:        name,
		connection:  i,
		Commander:   gobot.NewCommander(),
		initialized: false,
	}

	b.AddCommand("ThermopileInput", func(params map[string]interface{}) interface{} {
		local, object, err := b.ReadSensor()
		return map[string]interface{}{"local": fmt.Sprintf("%X", local), "object": fmt.Sprintf("%X", object), "err": err}
	})

	return b
}

func (b *Tmp007Driver) Name() string                 { return b.name }
func (b *Tmp007Driver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start writes start bytes
func (b *Tmp007Driver) Start() (errs []error) {
	b.initialized = true
	return
}

// Halt returns true if device is halted successfully
func (b *Tmp007Driver) Halt() (errs []error) { return }

// Digital input
func (b *Tmp007Driver) ReadSensor() (local []byte, object []byte, errs []error) {
	localTemp, errLocal := b.connection.I2cReadRegister([]byte{TMP007_ADDRESS, TMP007_LOCAL_TEMPERATURE}, 2)
	//	v := fmt.Fscanf(b.swap16(localTemp), "%f", &v)
	//	fmt.Printf("The local temperature is:  [%X] %X degree C\n", localTemp, v/128.0)
	// fmt.Printf("local %X \n", localTemp)

	objectTemp, errObject := b.connection.I2cReadRegister([]byte{TMP007_ADDRESS, TMP007_OBJECT_TEMPERATURE}, 2)
	// v = fmt.Fscanf(b.swap16(objectTemp), "%f", &v)
	// fmt.Printf("The object temperature is:  [%X] %X degree C\n", objectTemp, v/128.0)
	// fmt.Printf("object %X \n", objectTemp)

	return localTemp, objectTemp, []error{errLocal, errObject}
}

/*
// Swap bytes to change endianness
func (b *Tmp007Driver) swap16(val []byte) (data []byte) {
	return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}
*/

/*

// Read TMP007 thermopile sensor
function readTempopile() {
    var localTemp = i2c.readWordSync(TMP007_ADDR, TMP007_LOCAL_TEMP);
    localTemp = (swap16(localTemp))/128.0;
    console.log("The local temperature is:  " + localTemp.toString() + " degree C");

    var objectTemp = i2c.readWordSync(TMP007_ADDR, TMP007_OBJ_TEMP);
    objectTemp = (swap16(objectTemp))/128.0;
    console.log("The object temperature is: " + objectTemp.toString() + " degree C");
}

*/
