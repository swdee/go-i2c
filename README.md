# I2C-bus interaction for peripheral sensors on Single Board Computers

[![Go Report Card](https://goreportcard.com/badge/github.com/swdee/go-i2c)](https://goreportcard.com/report/github.com/swdee/go-i2c)
[![GoDoc](https://pkg.go.dev/badge/github.com/swdee/go-i2c)](https://godoc.org/github.com/swdee/go-i2c)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This library written in [Go programming language](https://golang.org/) intended to activate and interact with the I2C bus by reading and writing data.


## Usage


```go
func main() {
  // create new connection to I2C dev on /dev/i2c-0 with address 0x27
  i2cDevice, _ := i2c.New(0x27, "/dev/i2c-0")

  // free I2C connection on exit
  defer i2cDevice.Close()

  // write to command 0x1 the value of 0xF3
  if _, err := i2cDevice.WriteBytes([]byte{0x1, 0xF3}); err != nil {
    return err
  }
}
```
Note: Error handling has been skipped for brevity.




## Fork

Forked from https://github.com/googolgl/go-i2c as this library is unresponsive to [Pull requests](https://github.com/googolgl/go-i2c/pull/3).  Changes in this fork include:

* Added `ReadRegU32BE()` and `WriteRegBytes()` functions.
* Added `WriteThenReadBytes()` function.
* Clean up of code: removal of Logrus debug logging and CGO.


## Tutorial

In [repositories](https://github.com/d2r2?tab=repositories) contain quite a lot projects, which use i2c library as a starting point to interact with various peripheral devices and sensors for use on embedded Linux devices. All these libraries start with a standard call to open I2C-connection to specific bus line and address, than pass i2c instance to device.

You will find here the list of all devices and sensors supported by me, that reference this library:

- [Liquid-crystal display driven by Hitachi HD44780 IC](https://github.com/d2r2/go-hd44780).
- [BMP180/BMP280/BME280 temperature and pressure sensors](https://github.com/d2r2/go-bsbmp).
- [DHT12/AM2320 humidity and temperature sensors](https://github.com/d2r2/go-aosong).
- [Si7021 relative humidity and temperature sensor](https://github.com/d2r2/go-si7021).
- [SHT3x humidity and temperature sensor](https://github.com/d2r2/go-sht3x).
- [VL53L0X time-of-flight ranging sensor](https://github.com/d2r2/go-vl53l0x).
- [BH1750 ambient light sensor](https://github.com/d2r2/go-bh1750).
- [MPL3115A2 pressure and temperature sensor](https://github.com/d2r2/go-mpl3115a2).
- [PCA9685 16-Channel 12-Bit PWM Driver](https://github.com/googolgl/go-pca9685).
- [MCP23017 16-Bit I/O Expander with Serial Interface Driver](https://github.com/googolgl/go-mcp23017).


## Getting help

GoDoc [documentation](https://godoc.org/github.com/swdee/go-i2c)


## Troubleshooting


#### How to enable I2C bus

Various SBC's (Single Board Computers) has vendor specific methods for activating the I2C
bus on the GPIO Pins.  On the Raspberry Pi you may need to activate using the `raspi-config`
utility.  On Radxa Rock Pi devices uses `rsetup` to activate through Overlays.

After activating the I2C bus a reboot is usually required for it to start working. 


#### How to find I2C bus allocation and device address

Use the `i2cdetect` utility to check your i2c bus is active.
```
$ i2cdetect -l

i2c-10  i2c             ddc                                     I2C adapter
i2c-1   i2c             rk3x-i2c                                I2C adapter
i2c-6   i2c             rk3x-i2c                                I2C adapter
i2c-4   i2c             rk3x-i2c                                I2C adapter
i2c-11  i2c             ddc                                     I2C adapter
i2c-0   i2c             rk3x-i2c                                I2C adapter
i2c-9   i2c             fde50000.dp                             I2C adapter
i2c-7   i2c             rk3x-i2c                                I2C adapter

```

To discover is a sensor is active on the bus, scan the appropriate bus it is connected to.
```
$ i2cdetect -y 0

     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:                         -- -- -- -- -- -- -- -- 
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
60: 60 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
70: -- -- -- -- -- -- -- --           
```

Above we can see that the VNCL4040 is connected at address `0x60`.


## License

Go-i2c is licensed under MIT License.
