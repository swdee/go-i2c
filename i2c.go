// Package i2c provides low level control over the Linux i2c bus.
//
// Before usage you should load the i2c-dev kernel module
//
//	sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// Linux systems contain several buses.
package i2c

import (
	"os"
	"syscall"
	"unsafe"
)

// Options represents a connection to I2C-device.
type Options struct {
	addr uint8
	dev  string
	rc   *os.File
}

// i2c_msg struct represents an I2C message
type i2c_msg struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}

// i2c_rdwr_ioctl_data struct for I2C_RDWR ioctl operation
type i2c_rdwr_ioctl_data struct {
	msgs  uintptr
	nmsgs uint32
}

// New opens a connection for I2C-device.
// SMBus (System Management Bus) protocol over I2C
// supported as well: you should preliminary specify
// register address to read from, either write register
// together with the data in case of write operations.
func New(addr uint8, dev string) (*Options, error) {

	i2c := &Options{
		addr: addr,
		dev:  dev,
	}

	f, err := os.OpenFile(dev, os.O_RDWR, 0600)

	if err != nil {
		return i2c, err
	}

	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return i2c, err
	}

	i2c.rc = f
	return i2c, nil
}

// GetAddr return device occupied address in the bus.
func (o *Options) GetAddr() uint8 {
	return o.addr
}

// GetDev return full device name.
func (o *Options) GetDev() string {
	return o.dev
}

// READ SECTION

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (o *Options) ReadBytes(buf []byte) (int, error) {

	n, err := o.rc.Read(buf)

	if err != nil {
		return n, err
	}

	return n, nil
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
func (o *Options) ReadRegBytes(reg byte, n int) ([]byte, int, error) {

	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return nil, 0, err
	}

	buf := make([]byte, n)
	c, err := o.ReadBytes(buf)

	if err != nil {
		return nil, 0, err
	}

	return buf, c, nil
}

// ReadRegU8 reads byte from I2C-device register specified in reg.
func (o *Options) ReadRegU8(reg byte) (byte, error) {

	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 1)

	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16BE(reg byte) (uint16, error) {

	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 2)

	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	w := uint16(buf[0])<<8 + uint16(buf[1])

	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16LE(reg byte) (uint16, error) {

	w, err := o.ReadRegU16BE(reg)

	if err != nil {
		return 0, err
	}

	// exchange bytes
	w = (w&0xFF)<<8 + w>>8

	return w, nil
}

// ReadRegS16BE reads signed big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16BE(reg byte) (int16, error) {

	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 2)

	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	w := int16(buf[0])<<8 + int16(buf[1])

	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16LE(reg byte) (int16, error) {

	w, err := o.ReadRegS16BE(reg)

	if err != nil {
		return 0, err
	}

	// exchange bytes
	w = (w&0xFF)<<8 + w>>8

	return w, nil
}

// ReadRegU32BE reads unsigned big endian word (32 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU32BE(reg byte) (uint32, error) {

	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 4)

	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	w := uint32(buf[0])<<24 | uint32(buf[1])<<16 |
		uint32(buf[2])<<8 | uint32(buf[3])

	return w, nil
}

// WRITE SECTION

// WriteBytes send bytes to the remote I2C-device. The interpretation of
// the message is implementation-dependent.
func (o *Options) WriteBytes(buf []byte) (int, error) {
	return o.rc.Write(buf)
}

// WriteRegBytes send bytes to the remote I2C-device starting from reg address.
func (o *Options) WriteRegBytes(reg byte, buf []byte) (int, error) {
	b := append([]byte{reg}, buf...)
	return o.WriteBytes(b)
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
func (o *Options) WriteRegU8(reg byte, value byte) error {

	buf := []byte{reg, value}

	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// WriteRegU16BE writes unsigned big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16BE(reg byte, value uint16) error {

	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}

	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return o.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16BE(reg byte, value int16) error {

	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}

	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return o.WriteRegS16BE(reg, w)
}

// WriteRegU24BE writes unsigned big endian word (24 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU24BE(reg byte, value uint32) error {

	buf := []byte{reg, byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}

	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// WriteRegU32BE writes unsigned big endian word (32 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU32BE(reg byte, value uint32) error {

	buf := []byte{reg, byte(value >> 24 & 0xFF), byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}

	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// WriteThenReadBytes sends two I2C messages, the first to write some bytes then
// the second to read them.  This function allows us to perform a Write then Read
// without a I2C Stop condition occurring between the two messages which
// happens if WriteBytes() then ReadBytes() functions were called individually.
func (o *Options) WriteThenReadBytes(writeBuf, readBuf []byte) (int, int, error) {

	msgs := []i2c_msg{
		{
			addr:  uint16(o.addr),
			flags: 0,
			len:   uint16(len(writeBuf)),
			buf:   uintptr(unsafe.Pointer(&writeBuf[0])),
		},
		{
			addr:  uint16(o.addr),
			flags: I2C_M_RD,
			len:   uint16(len(readBuf)),
			buf:   uintptr(unsafe.Pointer(&readBuf[0])),
		},
	}

	data := i2c_rdwr_ioctl_data{
		msgs:  uintptr(unsafe.Pointer(&msgs[0])),
		nmsgs: uint32(len(msgs)),
	}

	if err := ioctl(o.rc.Fd(), I2C_RDWR, uintptr(unsafe.Pointer(&data))); err != nil {
		return 0, 0, err
	}

	return len(writeBuf), len(readBuf), nil
}

// Close I2C-connection.
func (o *Options) Close() error {
	return o.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) error {

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, arg); err != 0 {

		return err
	}

	return nil
}
