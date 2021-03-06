// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hpc015/util implements string to data convert.
package hpc015

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func reverseU16(flag uint16) uint16 {
	return ((flag & 0xFF) << 8) | ((flag & 0xFF00) >> 8)
}

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

// euqalDate compare tow time, by their year, month, secounds day
// if they are same, return true, else return false
func euqalDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// equalClock compare tow time, by their hour, minute, secounds only
// if they are same, return true, else return false
func equalClock(t1, t2 time.Time) bool {
	h1, m1, s1 := t1.Clock()
	h2, m2, s2 := t2.Clock()

	return h1 == h2 && m1 == m2 && s1 == s2
}

// equalClockOmitSec compare tow time, by their hour, minute only
// if they are same, return true, else return false
func equalClockOmitSec(t1, t2 time.Time) bool {
	h1, m1, _ := t1.Clock()
	h2, m2, _ := t2.Clock()

	return h1 == h2 && m1 == m2
}

// equalTime compare tow time,
// by their year, month, secound, hour, minute, secounds only
// means, ignore under millisecond
// if they are same, return true, else return false
func equalTime(t1, t2 time.Time) bool {
	return (euqalDate(t1, t2) == true) && (equalClock(t1, t2) == true)
}

// equalTime compare tow time,
// by their year, month, secound, hour, minute only
// means, ignore under millisecond
// if they are same, return true, else return false
func equalTimeOmitSec(t1, t2 time.Time) bool {
	return (euqalDate(t1, t2) == true) && (equalClockOmitSec(t1, t2) == true)
}
