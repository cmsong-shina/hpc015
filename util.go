// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hpc015/util implements string to data convert.
package hpc015

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func readU8(s string) (string, uint8, error) {
	if len(s) < 1 {
		return s, 0, fmt.Errorf("failed to read uint8: length must be 1 byte, but came %d byte", len(s))
	}

	n, err := strconv.ParseUint(s[:2], 16, 1*8)
	if err != nil {
		return s, 0, fmt.Errorf("")
	}
	return s[2:], uint8(n), nil
}

func readU16(s string) (string, uint16, error) {
	if len(s) < 2 {
		return s, 0, fmt.Errorf("failed to read uint16: length must be 2 byte, but came %d byte", len(s))
	}

	n, err := strconv.ParseUint(s[:4], 16, 2*8)
	if err != nil {
		return s, 0, fmt.Errorf("failed to read uint16: %s", err.Error())
	}
	return s[4:], uint16(n), nil
}

func readU32(s string) (string, uint32, error) {
	if len(s) < 4 {
		return s, 0, fmt.Errorf("failed to read uint32: length must be 4 byte, but came %d byte", len(s))
	}

	n, err := strconv.ParseUint(s[:8], 16, 4*8)
	if err != nil {
		return s, 0, fmt.Errorf("failed to read uint32: %s", err.Error())
	}
	return s[8:], uint32(n), nil
}

func readU64(s string) (string, uint64, error) {
	if len(s) < 8 {
		return s, 0, fmt.Errorf("failed to read uint64: length must be 8 byte, but came %d byte", len(s))
	}

	n, err := strconv.ParseUint(s[:16], 16, 8*8)
	if err != nil {
		return s, 0, fmt.Errorf("failed to read uint64: %s", err.Error())
	}
	return s[16:], n, nil
}

func readBytes(s string, length int) (string, []byte, error) {
	if len(s) < length {
		return s, nil, fmt.Errorf("failed to read bytes: length must be %d byte, but came %d byte", length, len(s))
	}

	data, err := hex.DecodeString(s[:length*2])
	if err != nil {
		return s, nil, fmt.Errorf("failed to read bytes: %s", err.Error())
	}
	return s[length*2:], data, nil
}
