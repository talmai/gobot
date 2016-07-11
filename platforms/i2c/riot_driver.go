package i2c

import (
	"fmt"
	"github.com/talmai/gobot"
)

var _ gobot.Driver = (*RIoTDriver)(nil)

const (
	RIOT_ADDRESS = 0x20

	RIOT_INITIALIZATION_ADDRESS_01 = 0x00
	RIOT_INITIALIZATION_ADDRESS_02 = 0x01

	RIOT_DIGITAL_INPUT_REGISTER  = 0x09
	RIOT_DIGITAL_OUTPUT_REGISTER = 0x0A

	// relays are NORMALLY CLOSED (even if the relay fails, lights continue on)
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_OFF = 0x10
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_ON  = 0xEF
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_OFF  = 0x20
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_ON   = 0xDF

	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_SET   = 0x40
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_RESET = 0xBF
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_SET    = 0x80
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_RESET  = 0x7F
)

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

	b.AddCommand("ReadDigitalInput", func(params map[string]interface{}) interface{} {
		data, err := b.ReadDigitalInput()
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
	})

	b.AddCommand("SetDigitalOutputChannelZero", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_SET)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetDigitalOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_SET)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetRelayOutputChannelZero", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_OFF)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetRelayOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_OFF)
		return map[string]interface{}{"err": err}
	})

	return b
}

func (b *RIoTDriver) Name() string                 { return b.name }
func (b *RIoTDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start writes start bytes
func (b *RIoTDriver) Start() (errs []error) {
	if !b.initialized {
		if err := b.connection.I2cStart(RIOT_ADDRESS); err != nil {
			return []error{err}
		}
		b.initialized = true
	}
	return
}

// Halt returns true if device is halted successfully
func (b *RIoTDriver) Halt() (errs []error) { return }

// initializes RIoT board
func (b *RIoTDriver) initializeRIoTInterfaceBoard() (errs []error) {
	if !b.digitalIoInitialized {
		// Digital I/O initialization:
		// i2c.writeByteSync(0x20, 0x00, 0x0F);    i2c.writeByteSync(0x20, 0x01, 0x00);
		if err := b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_INITIALIZATION_ADDRESS_01, 0x0F}); err != nil {
			return []error{err}
		}
		if err := b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_INITIALIZATION_ADDRESS_02, 0x00}); err != nil {
			return []error{err}
		}
		b.digitalIoInitialized = true
	}
	return nil
}

// Digital input
func (b *RIoTDriver) ReadDigitalInput() (data []byte, errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}
	// The lower four bits of “input” corresponding to digital input channel 0-3
	data, err := b.connection.I2cReadRegister([]byte{RIOT_ADDRESS, RIOT_DIGITAL_INPUT_REGISTER}, 1)
	return data, []error{err}
}

// Digital output
func (b *RIoTDriver) SetDigitalOutput(channel byte) (errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}
	// read current register value
	if data, err := b.ReadDigitalInput(); err != nil {
		return []error{err}
	}
	if err = b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, data | channel}); err != nil {
		return []error{err}
	}
	return
}

// Digital output
func (b *RIoTDriver) ResetDigitalOutput(channel byte) (errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}
	// read current register value
	data, err := b.ReadDigitalInput()
	if err = b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, channel}); err != nil {
		return []error{err}
	}
	if err = b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, channel}); err != nil {
		return []error{err}
	}
	return
}
