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

	// relays are NORMALLY CLOSED (even if the relay fails, lights continue on)
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_OFF = 0x10
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_ON  = 0xEF
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_OFF  = 0x20
	RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_ON   = 0xDF

	// DAC
	RIOT_DIGITAL_INPUT_REGISTER  = 0x09
	RIOT_DIGITAL_OUTPUT_REGISTER = 0x0A

	RIOT_DIGITAL_TO_ANALOG_CONVERTER_SLAVE_ADDRESS_ONE = 0x60
	RIOT_DIGITAL_TO_ANALOG_CONVERTER_SLAVE_ADDRESS_TWO = 0x61

	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_SET   = 0x40
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_RESET = 0xBF
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_SET    = 0x80
	RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_RESET  = 0x7F

	// ADC
	RIOT_ANALOG_TO_DIGITAL_OUTPUT_REGISTER = 0x00
	RIOT_ANALOG_TO_DIGITAL_INIT_REGISTER   = 0x01

	RIOT_ANALOG_TO_DIGITAL_CONVERTER_SLAVE_ADDRESS = 0x49

	RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_ZERO  = 0x83C5
	RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_ONE   = 0x83D5
	RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_TWO   = 0x83E5
	RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_THREE = 0x83F5
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

	b.AddCommand("ResetDigitalOutputChannelZero", func(params map[string]interface{}) interface{} {
		err := b.ResetDigitalOutput(RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ZERO_RESET)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetDigitalOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_SET)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("ResetDigitalOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.ResetDigitalOutput(RIOT_GPIO_DIGITAL_OUTPUT_CHANNEL_ONE_RESET)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetRelayOutputChannelZero", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_OFF)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("ResetRelayOutputChannelZero", func(params map[string]interface{}) interface{} {
		err := b.ResetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ZERO_ON)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("SetRelayOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_OFF)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("ResetRelayOutputChannelOne", func(params map[string]interface{}) interface{} {
		err := b.ResetDigitalOutput(RIOT_GPIO_RELAY_OUTPUT_CHANNEL_ONE_ON)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("DimLuminaireUp", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalAnalogConverter(0x0F, 0xFF)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("DimLuminaireDown", func(params map[string]interface{}) interface{} {
		err := b.SetDigitalAnalogConverter(0x00, 0x00)
		return map[string]interface{}{"err": err}
	})

	b.AddCommand("ReadADCChannelZero", func(params map[string]interface{}) interface{} {
		data, err := b.ReadADC(RIOT_ANALOG_TO_DIGITAL_INIT_REGISTER, RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_ZERO)
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
	})

	b.AddCommand("ReadADCChannelOne", func(params map[string]interface{}) interface{} {
		data, err := b.ReadADC(RIOT_ANALOG_TO_DIGITAL_INIT_REGISTER, RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_ONE)
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
	})

	b.AddCommand("ReadADCChannelTwo", func(params map[string]interface{}) interface{} {
		data, err := b.ReadADC(RIOT_ANALOG_TO_DIGITAL_INIT_REGISTER, RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_TWO)
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
	})

	b.AddCommand("ReadADCChannelThree", func(params map[string]interface{}) interface{} {
		data, err := b.ReadADC(RIOT_ANALOG_TO_DIGITAL_INIT_REGISTER, RIOT_ANALOG_TO_DIGITAL_CONVERTER_INPUT_CHANNEL_THREE)
		return map[string]interface{}{"raw": fmt.Sprintf("%X", data), "digitalInput01": fmt.Sprintf("%X", data[0]&0X01), "digitalInput02": fmt.Sprintf("%X", data[0]&0X02>>1), "digitalInput03": fmt.Sprintf("%X", data[0]&0X04>>2), "digitalInput04": fmt.Sprintf("%X", data[0]&0X08>>3), "err": err}
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
	data, _ := b.ReadDigitalInput()

	// fmt.Printf("[0]-> %X %X %X %X %X\n", data, data[0]&0X01, data[0]&0X02>>1, data[0]&0X04>>2, data[0]&0X08>>3) // fmt.Printf("[1]-> %X %X %X %X %X\n", data, data[0]&0X10>>4, data[0]&0X20>>5, data[0]&0X40>>6, data[0]&0X80>>7) // fmt.Printf("[ch]-> %X %X %X\n", channel, data[0]|channel)
	b.connection.I2cWriteWord(RIOT_ADDRESS, RIOT_DIGITAL_OUTPUT_REGISTER, uint16(data[0]|channel))
	// if err := b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, data[0] | channel}); err != nil {
	// 	return
	// }
	return
}

// Digital output
func (b *RIoTDriver) ResetDigitalOutput(channel byte) (errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}
	// read current register value
	data, _ := b.ReadDigitalInput()

	fmt.Printf("[0]-> %X %X %X %X %X\n", data, data[0]&0X01, data[0]&0X02>>1, data[0]&0X04>>2, data[0]&0X08>>3)
	fmt.Printf("[1]-> %X %X %X %X %X\n", data, data[0]&0X10>>4, data[0]&0X20>>5, data[0]&0X40>>6, data[0]&0X80>>7)
	fmt.Printf("[ch]-> %X %X %X\n", channel, data[0]&channel)

	b.connection.I2cWriteWord(RIOT_ADDRESS, RIOT_DIGITAL_OUTPUT_REGISTER, uint16(data[0]&channel))
	// if err := b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, data[0] | channel}); err != nil {
	// 	return
	// }
	return
}

// Digital Analog Converter
func (b *RIoTDriver) SetDigitalAnalogConverter(value01 byte, value02 uint16) (errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}

	b.connection.I2cWriteWord(RIOT_DIGITAL_TO_ANALOG_CONVERTER_SLAVE_ADDRESS_TWO, value01, value02)
	// if err := b.connection.I2cWrite(RIOT_ADDRESS, []byte{RIOT_DIGITAL_OUTPUT_REGISTER, data[0] | channel}); err != nil {
	// 	return
	// }
	return
}

// Analog to Digital Converter
func (b *RIoTDriver) ReadADC(value01 byte, value02 uint16) (data []byte, errs []error) {
	if err := b.initializeRIoTInterfaceBoard(); err != nil {
		return
	}
	b.connection.I2cWriteWord(RIOT_ANALOG_TO_DIGITAL_CONVERTER_SLAVE_ADDRESS, value01, value02)

	data, err := b.connection.I2cReadRegister([]byte{RIOT_ANALOG_TO_DIGITAL_CONVERTER_SLAVE_ADDRESS, RIOT_ANALOG_TO_DIGITAL_OUTPUT_REGISTER}, 2) // 2 == 2 bytes == word
	fmt.Printf("data %X \n", data)
	return data, []error{err}
}
