// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpc015

import (
	"reflect"
	"testing"
)

func Test_reads(t *testing.T) {

	input := "010142AE51520156000D0001E6A7"

	t.Run("hex 스트링 read 테스트", func(t *testing.T) {
		newS, got, err := readBytes(input, 4)
		if err != nil {
			t.Errorf("readBytes() error = %v", err)
			return
		}
		input = newS
		if !reflect.DeepEqual(got, []byte{0x01, 0x01, 0x42, 0xAE}) {
			t.Errorf("readBytes() = %v, want %v", got, []byte{0x01, 0x01, 0x42, 0xAE})
		}
	})

	t.Run("hex 스트링 readU8 테스트", func(t *testing.T) {
		newS, got, err := readU8(input)
		if err != nil {
			t.Errorf("readU8() error = %v", err)
			return
		}
		input = newS
		if !reflect.DeepEqual(got, 0x51) {
			t.Errorf("readU8() = %v, want %v", got, 0x51)
		}
	})

	t.Run("hex 스트링 readU16 테스트", func(t *testing.T) {
		newS, got, err := readU16(input)
		if err != nil {
			t.Errorf("readU16() error = %v", err)
			return
		}
		input = newS
		if !reflect.DeepEqual(got, 0x5201) {
			t.Errorf("readU16() = %v, want %v", got, 0x5201)
		}
	})

	t.Run("hex 스트링 readU32 테스트", func(t *testing.T) {
		newS, got, err := readU32(input)
		if err != nil {
			t.Errorf("readU32() error = %v", err)
			return
		}
		input = newS
		if !reflect.DeepEqual(got, 0x56000D00) {
			t.Errorf("readU32() = %v, want %v", got, 0x56000D00)
		}
	})
}
