package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*RIoTDriver)(nil)

const rIoTAddress = 0x20

type RIoTDriver struct {
	name       string
	connection I2cExtended
	gobot.Commander
	initialized          bool
	digitalIoInitialized bool
}

func NewRIoTDriver(i I2cExtended, name string) *RIoTDriver {
	b := &RIoTDriver{
		name:                 name,
		connection:           i,
		Commander:            gobot.NewCommander(),
		initialized:          false,
		digitalIoInitialized: false,
	}

	b.AddCommand("DigitalInput", func(params map[string]interface{}) interface{} {
		data, err := b.DigitalInput()
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
	})

	return b
}

func (b *RIoTDriver) Name() string                 { return b.name }
func (b *RIoTDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start writes start bytes
func (b *RIoTDriver) Start() (errs []error) {
	if !b.initialized {
		if err := b.connection.I2cStart(rIoTAddress); err != nil {
			return []error{err}
		}
		b.initialized = true
	}
	return
}

// Halt returns true if device is halted successfully
func (b *RIoTDriver) Halt() (errs []error) { return }

// Digital input
func (b *RIoTDriver) DigitalInput() (data []byte, errs []error) {
	if !b.digitalIoInitialized {
		// Digital I/O initialization:
		// i2c.writeByteSync(0x20, 0x00, 0x0F);    i2c.writeByteSync(0x20, 0x01, 0x00);
		if err := b.connection.I2cWrite(rIoTAddress, []byte{0x00, 0x0F}); err != nil {
			return nil, []error{err}
		}
		if err := b.connection.I2cWrite(rIoTAddress, []byte{0x01, 0x00}); err != nil {
			return nil, []error{err}
		}
		b.digitalIoInitialized = true
	}
	// The lower four bits of “input” corresponding to digital input channel 0-3
	data, err := b.connection.I2cReadRegister([]byte{rIoTAddress, 0x09}, 1)
	return data, []error{err}
}
