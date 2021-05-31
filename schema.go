// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hpc015/schema implements hpc015's data structure and more.
package hpc015

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// There are two kind of command, which are:
//   - `getsetting`: To obtain setting value.
//   - `cache`: To uplaod cache data.
type RequestSchema struct {
	Cmd    string        // for request
	Flag   uint16        // for request, means timestamp
	Data   [][]byte      // for request
	Status *DeviceStatus // for cache request, means information of device
	Count  uint16        // for cache request, means number of [Data]
	Result uint8         // for response, do not use this field
}

// NewRequestSchema makes RequestSchema from raw string.
//
// Returned value is still not useful,
// covert to getsetting struct or cache struct depend on RequestSchema.Cmd
func NewRequestSchema(reqestString string) (*RequestSchema, error) {
	fields := strings.Split(reqestString, "&")

	fieldsTable := make(map[string]string, 0)
	for _, v := range fields {
		field := strings.Split(v, "=")
		fieldsTable[field[0]] = field[1]
	}

	request := &RequestSchema{
		Data: make([][]byte, 0, 1),
	}

	v, ok := fieldsTable["cmd"]
	if ok {
		request.Cmd = v
	}

	v, ok = fieldsTable["status"]
	if ok {
		status, err := NewDeviceStatus(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode status: %s", err.Error())
		}
		request.Status = status
	}

	v, ok = fieldsTable["flag"]
	if ok {
		n, err := strconv.ParseUint(v, 16, 16)
		if err == nil {
			request.Flag = uint16(n)
		}
	}

	v, ok = fieldsTable["data"]
	if ok {
		data, err := hex.DecodeString(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode data: %s", err.Error())
		}
		request.Data = append(request.Data, data)
	}

	v, ok = fieldsTable["count"]
	if ok {
		n, err := strconv.ParseUint(v, 16, 16)
		if err == nil {
			request.Count = uint16(n)
		}
	}

	v, ok = fieldsTable["result"]
	if ok {
		n, err := strconv.ParseUint(v, 16, 16)
		if err == nil {
			request.Result = uint8(n)
		}
	}

	return request, nil
}

type Configuration struct {
	CommandType     *byte
	Speed           *byte
	RecordingCycle  *byte
	UploadCycle     *byte
	FixedTimeUpload *byte
	UploadClock     *time.Time
	Model           *byte
	DisableType     *byte
	SystemTime      *time.Time
	OpenClock       *time.Time
	CloseClock      *time.Time
}

func (data *Configuration) fromRequestFormat() {

}

type GetSettingRequest struct {
	SerialNumber    []byte
	CommandType     byte
	Speed           byte
	RecordingCycle  byte
	UploadCycle     byte
	FixedTimeUpload byte
	UploadHour1     byte
	UploadMinute1   byte
	UploadHour2     byte
	UploadMinute2   byte
	UploadHour3     byte
	UploadMinute3   byte
	UploadHour4     byte
	UploadMinute4   byte
	Model           byte
	DisableType     byte
	MacAddress1     []byte
	MacAddress2     []byte
	MacAddress3     []byte
	Year            byte
	Month           byte
	Day             byte
	Hour            byte
	Minute          byte
	Secound         byte
	Week            byte
	OpenHour        byte
	OpenMinute      byte
	CloseHour       byte
	CloseMinute     byte
	Crc16           uint16
}

func NewSettingRequest(data []byte) (*GetSettingRequest, error) {
	if len(data) != 53 {
		return nil, fmt.Errorf("failed to parse GetSettingRequest: length must be 53 byte, but came %d byte", len(data))
	}

	getSetting := &GetSettingRequest{
		SerialNumber:    data[0:4],
		CommandType:     data[4],
		Speed:           data[5],
		RecordingCycle:  data[6],
		UploadCycle:     data[7],
		FixedTimeUpload: data[8],
		UploadHour1:     data[9],
		UploadMinute1:   data[10],
		UploadHour2:     data[11],
		UploadMinute2:   data[12],
		UploadHour3:     data[13],
		UploadMinute3:   data[14],
		UploadHour4:     data[15],
		UploadMinute4:   data[16],
		Model:           data[17],
		DisableType:     data[18],
		MacAddress1:     data[19:26],
		MacAddress2:     data[26:33],
		MacAddress3:     data[33:40],
		Year:            data[40],
		Month:           data[41],
		Day:             data[42],
		Hour:            data[43],
		Minute:          data[44],
		Secound:         data[45],
		Week:            data[46],
		OpenHour:        data[47],
		OpenMinute:      data[48],
		CloseHour:       data[49],
		CloseMinute:     data[50],
		// TODO: ensure endian at here
		Crc16: binary.LittleEndian.Uint16(data[51:53]),
	}

	return getSetting, nil
}

// CRC16 verifies all byte calculation before crc fields(excluding “result=”)
// 1 BYTE Hi8
// 1 BYTE Low8
func (request GetSettingRequest) calcCrc16(data []byte) (uint16, error) {
	var crc uint16 = 0xFFFF

	if len(data) > 78 {
		return 0, errors.New("leng of data must less than 78")
	}

	for j := 0; j < len(data); j++ {
		crc ^= uint16(data[j])

		for i := 0; i < 8; i++ {
			if (crc & 0x01) == 1 {
				crc >>= 1
				crc ^= 0xA001
			} else {
				crc >>= 1
			}
		}
		data[j] = uint8(crc % 0x100)
		data[j+1] = uint8(crc / 0x100)
	}

	return crc, nil
}

type GetSettingResponse struct {
	RespondingType  RespondingType
	Flag            uint16
	SerialNumber    []byte
	CommandType     byte
	Speed           byte
	RecordingCycle  byte
	UploadCycle     byte
	FixedTimeUpload byte
	UploadHour1     byte
	UploadMinute1   byte
	UploadHour2     byte
	UploadMinute2   byte
	UploadHour3     byte
	UploadMinute3   byte
	UploadHour4     byte
	UploadMinute4   byte
	Model           byte
	DisableType     byte
	MacAddress1     []byte
	MacAddress2     []byte
	MacAddress3     []byte
	Year            byte
	Month           byte
	Day             byte
	Hour            byte
	Minute          byte
	Second          byte
	Week            byte
	OpenHour        byte
	OpenMinute      byte
	CloseHour       byte
	CloseMinute     byte
	Reserved1       byte
	Reserved2       byte
	Crc16           uint16
}

func NewSettingResponse(request *GetSettingRequest, flag uint16) *GetSettingResponse {
	return &GetSettingResponse{
		RespondingType:  RespondingTypeConfirmation,
		Flag:            flag,
		SerialNumber:    []byte{0, 0, 0, 0},
		CommandType:     request.CommandType,
		Speed:           request.Speed,
		RecordingCycle:  request.RecordingCycle,
		UploadCycle:     request.UploadCycle,
		FixedTimeUpload: request.FixedTimeUpload,
		UploadHour1:     request.UploadHour1,
		UploadMinute1:   request.UploadMinute1,
		UploadHour2:     request.UploadHour2,
		UploadMinute2:   request.UploadMinute2,
		UploadHour3:     request.UploadHour3,
		UploadMinute3:   request.UploadMinute3,
		UploadHour4:     request.UploadHour4,
		UploadMinute4:   request.UploadMinute4,
		Model:           request.Model,
		DisableType:     request.DisableType,
		MacAddress1:     []byte{0, 0, 0, 0, 0, 0, 0},
		MacAddress2:     []byte{0, 0, 0, 0, 0, 0, 0},
		MacAddress3:     []byte{0, 0, 0, 0, 0, 0, 0},
		Year:            request.Year,
		Month:           request.Month,
		Day:             request.Day,
		Hour:            request.Hour,
		Minute:          request.Minute,
		Second:          request.Secound,
		Week:            0,
		OpenHour:        request.OpenHour,
		OpenMinute:      request.OpenMinute,
		CloseHour:       request.CloseHour,
		CloseMinute:     request.CloseMinute,
		Reserved1:       0,
		Reserved2:       0,
		Crc16:           request.Crc16,
	}
}

func (resp GetSettingResponse) GetConfiguration() *Configuration {
	var uploadClock = time.Date(
		0,
		0,
		0,
		int(resp.UploadHour1),
		int(resp.UploadMinute1),
		0,
		0,
		nil,
	)
	var SystemTime = time.Date(
		int(resp.Year),
		time.Month(resp.Month),
		int(resp.Day),
		int(resp.Hour),
		int(resp.Minute),
		int(resp.Second),
		0,
		nil,
	)

	var OpenClock = time.Date(
		int(resp.Year),
		time.Month(resp.Month),
		int(resp.Day),
		int(resp.Hour),
		int(resp.Minute),
		int(resp.Second),
		0,
		nil,
	)

	var CloseClock = time.Date(
		int(resp.Year),
		time.Month(resp.Month),
		int(resp.Day),
		int(resp.Hour),
		int(resp.Minute),
		int(resp.Second),
		0,
		nil,
	)

	return &Configuration{
		CommandType:     &resp.CommandType,
		Speed:           &resp.Speed,
		RecordingCycle:  &resp.RecordingCycle,
		UploadCycle:     &resp.UploadCycle,
		FixedTimeUpload: &resp.FixedTimeUpload,
		UploadClock:     &uploadClock,
		Model:           &resp.Model,
		DisableType:     &resp.DisableType,
		SystemTime:      &SystemTime,
		OpenClock:       &OpenClock,
		CloseClock:      &CloseClock,
	}
}

/// SetConfiguration apply configuration
/// if configuration is diffrent, mark RespondingType as NewParameterValue(0x04)
func (response *GetSettingResponse) SetConfiguration(cog Configuration) (bool, error) {
	original := response.GetConfiguration()

	if cog.CommandType != nil && original.CommandType != cog.CommandType {
		response.CommandType = *cog.CommandType
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.Speed != nil && original.Speed != cog.Speed {
		response.Speed = *cog.Speed
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.RecordingCycle != nil && original.RecordingCycle != cog.RecordingCycle {
		response.RecordingCycle = *cog.RecordingCycle
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.UploadCycle != nil && original.UploadCycle != cog.UploadCycle {
		response.UploadCycle = *cog.UploadCycle
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.FixedTimeUpload != nil && original.FixedTimeUpload != cog.FixedTimeUpload {
		response.FixedTimeUpload = *cog.FixedTimeUpload
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.UploadClock != nil && original.UploadClock != cog.UploadClock {
		response.UploadHour1 = byte(cog.UploadClock.Hour())
		response.UploadMinute1 = byte(cog.UploadClock.Minute())
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.Model != nil && original.Model != cog.Model {
		response.Model = *cog.Model
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.DisableType != nil && original.DisableType != cog.DisableType {
		response.DisableType = *cog.DisableType
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.SystemTime != nil && original.SystemTime != cog.SystemTime {
		response.Year = byte(cog.SystemTime.Year())
		response.Month = byte(cog.SystemTime.Month())
		response.Day = byte(cog.SystemTime.Day())
		response.Hour = byte(cog.SystemTime.Hour())
		response.Minute = byte(cog.SystemTime.Minute())
		response.Second = byte(cog.SystemTime.Second())
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.OpenClock != nil && original.OpenClock != cog.OpenClock {
		response.OpenHour = byte(cog.OpenClock.Hour())
		response.OpenMinute = byte(cog.OpenClock.Minute())
		response.RespondingType = RespondingTypeNewParameterValue
	}

	if cog.CloseClock != nil && original.CloseClock != cog.CloseClock {
		response.CloseHour = byte(cog.CloseClock.Hour())
		response.CloseMinute = byte(cog.CloseClock.Minute())
		response.RespondingType = RespondingTypeNewParameterValue
	}

	return false, nil
}

func (request GetSettingResponse) Binary() []byte {
	// buf := make([]byte, 0, 58)
	// buf[0] = byte(request.RespondingType)
	// binary.LittleEndian.PutUint16(buf[1:3], request.Flag)

	buf := bytes.NewBuffer(make([]byte, 0, 58))
	// buf.WriteByte(byte(request.RespondingType))
	// buf.WriteRune(rune(request.Flag))

	binary.Write(buf, binary.LittleEndian, request.RespondingType)
	binary.Write(buf, binary.LittleEndian, request.Flag)
	binary.Write(buf, binary.LittleEndian, request.SerialNumber)
	binary.Write(buf, binary.LittleEndian, request.CommandType)
	binary.Write(buf, binary.LittleEndian, request.Speed)
	binary.Write(buf, binary.LittleEndian, request.RecordingCycle)
	binary.Write(buf, binary.LittleEndian, request.UploadCycle)
	binary.Write(buf, binary.LittleEndian, request.FixedTimeUpload)
	binary.Write(buf, binary.LittleEndian, request.UploadHour1)
	binary.Write(buf, binary.LittleEndian, request.UploadMinute1)
	binary.Write(buf, binary.LittleEndian, request.UploadHour2)
	binary.Write(buf, binary.LittleEndian, request.UploadMinute2)
	binary.Write(buf, binary.LittleEndian, request.UploadHour3)
	binary.Write(buf, binary.LittleEndian, request.UploadMinute3)
	binary.Write(buf, binary.LittleEndian, request.UploadHour4)
	binary.Write(buf, binary.LittleEndian, request.UploadMinute4)
	binary.Write(buf, binary.LittleEndian, request.Model)
	binary.Write(buf, binary.LittleEndian, request.DisableType)
	binary.Write(buf, binary.LittleEndian, request.MacAddress1)
	binary.Write(buf, binary.LittleEndian, request.MacAddress2)
	binary.Write(buf, binary.LittleEndian, request.MacAddress3)
	binary.Write(buf, binary.LittleEndian, request.Year)
	binary.Write(buf, binary.LittleEndian, request.Month)
	binary.Write(buf, binary.LittleEndian, request.Day)
	binary.Write(buf, binary.LittleEndian, request.Hour)
	binary.Write(buf, binary.LittleEndian, request.Minute)
	binary.Write(buf, binary.LittleEndian, request.Second)
	binary.Write(buf, binary.LittleEndian, request.Week)
	binary.Write(buf, binary.LittleEndian, request.OpenHour)
	binary.Write(buf, binary.LittleEndian, request.OpenMinute)
	binary.Write(buf, binary.LittleEndian, request.CloseHour)
	binary.Write(buf, binary.LittleEndian, request.CloseMinute)
	binary.Write(buf, binary.LittleEndian, request.Reserved1)
	binary.Write(buf, binary.LittleEndian, request.Reserved2)
	binary.Write(buf, binary.LittleEndian, request.Crc16)

	return buf.Bytes()
}

// 0x00 exclude the verification hours and business hours
// 0x01 include the time of verifying the system
// 0x02 include the time of verifying the business hours
// 0x03 include the time of verifying the system and business hours
type CommandType uint8

// Operation mode
type ModelEnum byte

const (
	OpModeGridConnected           ModelEnum = 0
	OpModeStandAlone              ModelEnum = 0
	OpModeStandAloneWithoutUpload ModelEnum = 1
)

type DiableTypeEnum byte

const (
	DpModeNotDisplay       DiableTypeEnum = 0
	DpModeTotalDisplay     DiableTypeEnum = 1
	DpModebilateralDisplay DiableTypeEnum = 2
)

type RespondingType byte

const (
	RespondingTypeNewParameterValue = 0x04 // new parameter value
	RespondingTypeConfirmation      = 0x05 // parameter confirmation, after confirmation and responding, the parameter will be neglected.
)

// status=010142AE51520156000D0001E6A7
//   - 0101 Indicates the device firmware version number, version 1.1
//   - 42AE5152 indicates the device SN number, with the lower digit first, that is, The device SN number is 5251AE42
//   - 01 Device retention information
//   - 56 Represents the remaining power of the infrared transmitter battery used, The current remaining power is 86%
//   - 00 Device retention information
//   - 0D Indicates the current remaining capacity of the counter battery, the current remaining capacity is 13%
//
type DeviceStatus struct {
	Version            uint16
	SerialNumber       uint32
	Retention1         uint8
	TransmitterBattery uint8
	Retention2         uint8
	ReceiverBattery    uint8
	WTF                uint32 //not specified in manual
}

func NewDeviceStatus(data string) (*DeviceStatus, error) {
	if len(data) != 28 {
		return nil, fmt.Errorf("failed to parse GetSettingRequest: length must be 53 byte, but came %d byte", len(data))
	}

	var status = new(DeviceStatus)
	// buf := bytes.NewBufferString(data)
	// versionString, err := buf.ReadBytes(2)
	// if err != nil {
	// 	return nil, err
	// }
	// version, err := strconv.ParseUint(string(versionString), 16, 16)
	// if err != nil {
	// 	return nil, err
	// }
	// status.Version = uint16(version)

	data, version, err := readU16(data)
	if err != nil {
		return nil, err
	}
	data, serialNumber, err := readU32(data)
	if err != nil {
		return nil, err
	}
	data, retention1, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, transmitterBattery, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, retention2, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, receiverBattery, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, wtf, err := readU32(data)
	if err != nil {
		return nil, err
	}

	status.Version = version
	status.SerialNumber = serialNumber
	status.Retention1 = retention1
	status.TransmitterBattery = transmitterBattery
	status.Retention2 = retention2
	status.ReceiverBattery = receiverBattery
	status.WTF = wtf

	return status, nil
}
