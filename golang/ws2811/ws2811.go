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
)

const SK6812_STRIP_RGBW = 0x18100800
const SK6812_STRIP_RBGW = 0x18100008
const SK6812_STRIP_GRBW = 0x18081000
const SK6812_STRIP_GBRW = 0x18080010
const SK6812_STRIP_BRGW = 0x18001008
const SK6812_STRIP_BGRW = 0x18000810
const SK6812_SHIFT_WMASK = 0xf0000000

// 3 color R, G and B ordering
const WS2811_STRIP_RGB = 0x00100800
const WS2811_STRIP_RBG = 0x00100008
const WS2811_STRIP_GRB = 0x00081000
const WS2811_STRIP_GBR = 0x00080010
const WS2811_STRIP_BRG = 0x00001008
const WS2811_STRIP_BGR = 0x00000810

type Strip struct {
	ledstring C.ws2811_t
	LedCount  int
}

func NewStrip(stripType uint, ledCount uint16, gpioPin uint8, brightness uint8, channel uint8, invert int) (Strip, error) {
	var strip Strip
	ledstring := C.ws2811_t_create()
	for i := 0; i < 2; i++ {
		ledstring.channel[i].gpionum = C.int(0)
		ledstring.channel[i].count = C.int(0)
		ledstring.channel[i].invert = C.int(0)
		ledstring.channel[i].brightness = C.uint8_t(0)
	}
	ledstring.channel[channel].gpionum = C.int(gpioPin)
	ledstring.channel[channel].count = C.int(ledCount)
	ledstring.channel[channel].invert = C.int(invert)
	ledstring.channel[channel].brightness = C.uint8_t(brightness)
	ledstring.channel[channel].strip_type = C.int(stripType)
	res := int(C.ws2811_init(&ledstring))
	if res == 0 {
		strip = Strip{ledstring, int(ledCount)}
		return strip, nil
	} else {
		return strip, errors.New(fmt.Sprintf("Error ws2811.init.%d", res))
	}
}

func (s *Strip) Fini() {
	C.ws2811_fini(&s.ledstring)
}

func (s *Strip) Render() error {
	res := int(C.ws2811_render(&s.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.render.%d", res))
	}
}

func (s *Strip) NumPixels() int {
	return s.LedCount
}

// func Wait() error {
// 	res := int(C.ws2811_wait(&s.ledstring))
// 	if res == 0 {
// 		return nil
// 	} else {
// 		return errors.New(fmt.Sprintf("Error ws2811.wait.%d", res))
// 	}
// }

func (s *Strip) SetLed(index int, value uint32) {
	C.ws2811_set_led(&s.ledstring, C.int(index), C.uint32_t(value))
}

func ColorRGB(red, green, blue uint32) uint32 {
	return (red << 16) | (green << 8) | blue
}

func ColorRGBW(red, green, blue, white uint32) uint32 {
	return (white << 24) | (red << 16) | (green << 8) | blue
}

func (s *Strip) SetStrip(color uint32) {
	for i := 0; i < s.LedCount; i++ {
		s.SetLed(i, color)
	}
}

func ShiftColor(color uint32, goal uint32, step uint32) uint32 {
	w := (color & 0xff000000) >> 24
	r := (color & 0x00ff0000) >> 16
	g := (color & 0x0000ff00) >> 8
	b := (color & 0x000000ff)

	w2 := (goal & 0xff000000) >> 24
	r2 := (goal & 0x00ff0000) >> 16
	g2 := (goal & 0x0000ff00) >> 8
	b2 := (goal & 0x000000ff)

	bitShift := func(bit1, bit2 uint32) uint32 {
		if bit1 < bit2 {
			if (bit1 + step) > bit2 {
				return bit2
			}
			return bit1 + step
		}
		if bit1 > bit2 {
			if (bit1 - step) < bit2 {
				return bit2
			}
			return bit1 - bit2
		}
		return bit1
	}

	return bitShift(w, w2) | bitShift(r, r2) | bitShift(g, g2) | bitShift(b, b2)
}

// func Clear() {
// 	C.ws2811_clear(&s.ledstring)
// }
//
// func SetBitmap(a []uint32) {
// 	C.ws2811_set_bitmap(&s.ledstring, unsafe.Pointer(&a[0]), C.int(len(a)*4))
// }
