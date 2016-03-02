package sysfs

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	I2C_SLAVE = 0x0703
	I2C_SMBUS = 0x0720

	/* i2c_smbus_xfer read or write markers */
	I2C_SMBUS_WRITE = 0
	I2C_SMBUS_READ  = 1

	/* SMBus transaction types (size parameter in the above functions)
	   Note: these no longer correspond to the (arbitrary) PIIX4 internal codes! */
	I2C_SMBUS_QUICK            = 0
	I2C_SMBUS_BYTE             = 1
	I2C_SMBUS_BYTE_DATA        = 2
	I2C_SMBUS_WORD_DATA        = 3
	I2C_SMBUS_PROC_CALL        = 4
	I2C_SMBUS_BLOCK_DATA       = 5
	I2C_SMBUS_I2C_BLOCK_BROKEN = 6
	I2C_SMBUS_BLOCK_PROC_CALL  = 7
	I2C_SMBUS_I2C_BLOCK_DATA   = 8

	// Adapter functionality
	I2C_FUNCS                       = 0x0705
	I2C_FUNC_SMBUS_READ_BLOCK_DATA  = 0x01000000
	I2C_FUNC_SMBUS_WRITE_BLOCK_DATA = 0x02000000
)

type i2cSmbusIoctlData struct {
	readWrite byte
	command   byte
	size      uint32
	data      uintptr
}

type I2cDevice interface {
	io.ReadWriteCloser
	SetAddress(int) error
	ReadRegister([]byte, []byte) (int, error)
	WriteWord(uint8, uint16) (int, error)
}

type i2cDevice struct {
	file  File
	funcs uint64 // adapter functionality mask
}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location and device address
func NewI2cDevice(location string, address int) (d *i2cDevice, err error) {
	d = &i2cDevice{}

	if d.file, err = OpenFile(location, os.O_RDWR, os.ModeExclusive); err != nil {
		return
	}
	if err = d.queryFunctionality(); err != nil {
		return
	}

	err = d.SetAddress(address)

	return
}

func (d *i2cDevice) queryFunctionality() (err error) {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_FUNCS,
		uintptr(unsafe.Pointer(&d.funcs)),
	)

	if errno != 0 {
		err = fmt.Errorf("Querying functionality failed with syscall.Errno %v", errno)
	}
	fmt.Printf("Functionality: 0x%x\n", d.funcs)
	return
}

func (d *i2cDevice) SetAddress(address int) (err error) {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SLAVE,
		uintptr(byte(address)),
	)

	if errno != 0 {
		err = fmt.Errorf("Setting address failed with syscall.Errno %v", errno)
	}

	return
}

func (d *i2cDevice) Close() (err error) {
	return d.file.Close()
}

func (d *i2cDevice) Read(b []byte) (n int, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
		// Adapter doesn't support SMBus block read
		return d.file.Read(b)
	}

	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_READ,
		command:   0,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return n, fmt.Errorf("Read failed with syscall.Errno %v", errno)
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) ReadRegister(address []byte, b []byte) (n int, err error) {
	// if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
	// 	// Adapter doesn't support SMBus block read
	// 	return d.file.Read(b)
	// }

	register := byte(address[0])

	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_READ,
		command:   register,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return n, fmt.Errorf("Read failed with syscall.Errno %v", errno)
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) ReadWordRegister(address []byte, b []byte) (n int, err error) {
	if err = d.queryFunctionality(); err != nil {
		fmt.Printf("QueryFunctionality: %v \n", err)
	}

	// if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
	// 	// Adapter doesn't support SMBus block read
	// 	return d.file.Read(b)
	// }

	register := byte(address[0])

	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_READ,
		command:   register,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return n, fmt.Errorf("Read failed with syscall.Errno %v", errno)
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) Write(b []byte) (n int, err error) {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_BLOCK_DATA == 0 {
		// Adapter doesn't support SMBus block write
		return d.file.Write(b)
	}

	command := byte(b[0])
	buf := b[1:]

	data := make([]byte, len(buf)+1)
	data[0] = byte(len(buf))

	copy(data[1:], buf)

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_WRITE,
		command:   command,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		err = fmt.Errorf("Write failed with syscall.Errno %v", errno)
	}

	return len(b), err
}

func (d *i2cDevice) WriteWord(register uint8, value uint16) (n int, err error) {
	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_WRITE,
		command:   register,
		size:      I2C_SMBUS_WORD_DATA,
		data:      uintptr(unsafe.Pointer(&value)),
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		err = fmt.Errorf("Write failed with syscall.Errno %v", errno)
	}

	return 0, err
}
