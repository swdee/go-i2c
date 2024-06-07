//go:build linux && cgo
// +build linux,cgo

package i2c

// #include <linux/i2c-dev.h>
import "C"

// Get I2C_SLAVE constant value from
// Linux OS I2C declaration file.
const (
	I2C_SLAVE = C.I2C_SLAVE
	I2C_RDWR  = C.I2C_RDWR
	I2C_M_RD  = C.I2C_M_RD
)
