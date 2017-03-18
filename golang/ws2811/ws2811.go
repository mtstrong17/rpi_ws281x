// Copyright (c) 2015, Jacques Supcik, HEIA-FR
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the <organization> nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Interface to ws2811 chip (neopixel driver). Make sure that you have
ws2811.h and pwm.h in a GCC include path (e.g. /usr/local/include) and
libws2811.a in a GCC library path (e.g. /usr/local/lib).
See https://github.com/jgarff/rpi_ws281x for instructions
*/

package ws2811

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lws2811
#include "ws2811.go.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type Strip struct {
	ledCount   uint16
	pin        uint8
	freqHz     uint32
	dma        uint8
	invert     bool
	brightness uint8
	channel    uint8
	stripType  uint
}

func NewStrip(stripType uint, ledCount uint16, gpioPin uint8, brightness uint8, channel uint8, invert int) error {
	for i := 0; i < 2; i++ {
		C.ledstring.channel[i].gpionum = C.int(0)
		C.ledstring.channel[i].count = C.int(0)
		C.ledstring.channel[i].invert = C.int(0)
		C.ledstring.channel[i].brightness = C.uint8_t(0)
	}
	C.ledstring.channel[channel].gpionum = C.int(gpioPin)
	C.ledstring.channel[channel].count = C.int(ledCount)
	C.ledstring.channel[channel].invert = C.int(invert)
	C.ledstring.channel[channel].brightness = C.uint8_t(brightness)
	C.ledstring.channel[channel].strip_type = C.int(stripType)
	res := int(C.ws2811_init(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.init.%d", res))
	}
}

func Fini() {
	C.ws2811_fini(&C.ledstring)
}

func Render() error {
	res := int(C.ws2811_render(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.render.%d", res))
	}
}

func Wait() error {
	res := int(C.ws2811_wait(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.wait.%d", res))
	}
}

func SetLed(index int, value uint32) {
	C.ws2811_set_led(&C.ledstring, C.int(index), C.uint32_t(value))
}

func ColorRGB(red, green, blue uint32) uint32 {
	return (red << 16) | (green << 8) | blue
}
func ColorRGBW(red, green, blue, white uint32) uint32 {
	return (white << 24) | (red << 16) | (green << 8) | blue
}

func Clear() {
	C.ws2811_clear(&C.ledstring)
}

func SetBitmap(a []uint32) {
	C.ws2811_set_bitmap(&C.ledstring, unsafe.Pointer(&a[0]), C.int(len(a)*4))
}
