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

var (
	EnableDebugMessage = true
)

// RequestSchema basic form of request, it just hold bunch of raw data.
//
// Returned value is still not useful,
// covert to getsetting struct or cache struct depend on RequestSchema.Cmd
//
// There are two kind of commands, which are:
//   - `getsetting`: To obtain setting value
//   - `cache`: To upload cache data
//
// `Flag` always treated as BigEndian
type RequestSchema struct {
	Cmd    string   // for any request
	Flag   uint16   // for any request, means timestamp, bigendian
	Data   [][]byte // for any request
	Status string   // for cache request, means information of device
	Count  uint16   // for cache request, means number of [Data]
}

// NewRequestSchema makes RequestSchema from raw request string
func NewRequestSchema(reqestString string) (*RequestSchema, error) {
	if len(reqestString) < 30 {
		return nil, errors.New("length of request must be longer than 30")
	}

	fields := strings.Split(reqestString, "&")

	fieldsTable := make(map[string]string, 0)
	for _, v := range fields {
		field := strings.Split(v, "=")
		fieldsTable[field[0]] = field[1]
	}

	request := &RequestSchema{
		Data: make([][]byte, 0, 1),
	}

	// parse message
	for _, field := range fields {
		var k, v string

		{
			s := strings.Split(field, "=")
			if len(s) != 2 {
				continue
			}
			k = s[0]
			v = s[1]
		}

		switch k {
		case "cmd":
			request.Cmd = v

		case "status":
			request.Status = v

		case "flag":
			b, err := hex.DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode flag: %s", err.Error())
			}
			flag := binary.BigEndian.Uint16(b)
			request.Flag = flag

		case "data":
			data, err := hex.DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode data: %s", err.Error())
			}
			request.Data = append(request.Data, data)

		case "count":
			count, err := strconv.ParseUint(v, 16, 16)
			if err != nil {
				return nil, fmt.Errorf("failed to decode count: %s", err.Error())
			}
			request.Count = uint16(count)
		}
	}

	// validate schema, [Cmd, Data] can not be empty
	if request.Cmd == "" || len(request.Data) == 0 {
		return nil, errors.New("Cmd/Data field can not be empty")
	}

	return request, nil
}

// Configuration represent hpc015's configuration.
//
// SystemTime, OpenClock, CloseClock are mandatory.
//
// About Recording and Uploading,
// Recording means within business hour, timestamp interval of data,
// Uploading meas within business hour, specify the uploading time period via WIFI.
//
type Configuration struct {
	TimeVerifyMode        TimeVerifyMode
	Speed                 Speed
	RecordingCycle        byte // 1 to 225 min, 0 is real-time
	UploadCycle           byte // 1 to 225 min, 0 is real-time
	EnableFixedTimeUpload byte // TODO: Need to implement this feature
	UploadClock           time.Time
	NetworkType           NetworkType
	DisplayType           DisplayType
	SystemTime            time.Time
	OpenClock             time.Time
	CloseClock            time.Time
}

// Default configuration
func Default() *Configuration {
	return &Configuration{
		TimeVerifyMode:        Both,
		Speed:                 Low,
		RecordingCycle:        0,
		UploadCycle:           0,
		EnableFixedTimeUpload: 0,
		NetworkType:           Online,
		DisplayType:           Unidirectinal,
		SystemTime:            time.Now(),
		OpenClock:             time.Date(1, 1, 1, 0, 0, 0, 0, time.Local),
		CloseClock:            time.Date(1, 1, 1, 23, 59, 0, 0, time.Local),
	}
}

type GetSettingRequest struct {
	SerialNumber    []byte
	TimeVerifyMode  TimeVerifyMode
	Speed           Speed
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
	NetworkType     NetworkType
	DisplayType     DisplayType
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
	Crc16           uint16
}

// NewSettingRequest makes new GetSettingRequest instance.
//
//   - Length of [data] must be 53
//   - This function vaild CRC16
func NewSettingRequest(data []byte) (*GetSettingRequest, error) {
	if len(data) != 53 {
		return nil, fmt.Errorf("length must be 53 byte, but came %d byte", len(data))
	}

	crc, err := calcCrc16(data[:51])
	if err != nil {
		return nil, errors.New("failed to verify crc:" + err.Error())
	}

	incomeCrc := binary.BigEndian.Uint16(data[51:53])

	if crc != incomeCrc {
		return nil, errors.New("failed to verify crc: incorrect crc")
	}

	getSetting := &GetSettingRequest{
		SerialNumber:    data[0:4],
		TimeVerifyMode:  TimeVerifyMode(data[4]),
		Speed:           Speed(data[5]),
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
		NetworkType:     NetworkType(data[17]),
		DisplayType:     DisplayType(data[18]),
		MacAddress1:     data[19:26],
		MacAddress2:     data[26:33],
		MacAddress3:     data[33:40],
		Year:            data[40],
		Month:           data[41],
		Day:             data[42],
		Hour:            data[43],
		Minute:          data[44],
		Second:          data[45],
		Week:            data[46],
		OpenHour:        data[47],
		OpenMinute:      data[48],
		CloseHour:       data[49],
		CloseMinute:     data[50],
		Crc16:           crc,
	}

	return getSetting, nil
}

// Response generate response about request
//   - need to provider `flag`
//   - see also: `GetSettingResponse`
func (request GetSettingRequest) Response(flag uint16) *GetSettingResponse {
	return &GetSettingResponse{
		RespondingType:  Confirmation,
		Flag:            reverseU16(flag),
		SerialNumber:    []byte{0, 0, 0, 0},
		TimeVerifyMode:  request.TimeVerifyMode,
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
		NetworkType:     request.NetworkType,
		DisplayType:     request.DisplayType,
		MacAddress1:     []byte{0, 0, 0, 0, 0, 0, 0},
		MacAddress2:     []byte{0, 0, 0, 0, 0, 0, 0},
		MacAddress3:     []byte{0, 0, 0, 0, 0, 0, 0},
		Year:            request.Year,
		Month:           request.Month,
		Day:             request.Day,
		Hour:            request.Hour,
		Minute:          request.Minute,
		Second:          request.Second,
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

type GetSettingResponse struct {
	RespondingType  RespondingType
	Flag            uint16
	SerialNumber    []byte
	TimeVerifyMode  TimeVerifyMode
	Speed           Speed
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
	NetworkType     NetworkType
	DisplayType     DisplayType
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

func (resp GetSettingResponse) GetConfiguration() *Configuration {
	var uploadClock = time.Date(
		1,
		1,
		1,
		int(resp.UploadHour1),
		int(resp.UploadMinute1),
		0,
		0,
		time.Local,
	)
	var SystemTime = time.Date(
		int(resp.Year)+2000,
		time.Month(resp.Month),
		int(resp.Day),
		int(resp.Hour),
		int(resp.Minute),
		int(resp.Second),
		0,
		time.Local,
	)

	var OpenClock = time.Date(
		1,
		1,
		1,
		int(resp.OpenHour),
		int(resp.OpenMinute),
		0,
		0,
		time.Local,
	)

	var CloseClock = time.Date(
		1,
		1,
		1,
		int(resp.CloseHour),
		int(resp.CloseMinute),
		0,
		0,
		time.Local,
	)

	return &Configuration{
		TimeVerifyMode:        resp.TimeVerifyMode,
		Speed:                 resp.Speed,
		RecordingCycle:        resp.RecordingCycle,
		UploadCycle:           resp.UploadCycle,
		EnableFixedTimeUpload: resp.FixedTimeUpload,
		UploadClock:           uploadClock,
		NetworkType:           resp.NetworkType,
		DisplayType:           resp.DisplayType,
		SystemTime:            SystemTime,
		OpenClock:             OpenClock,
		CloseClock:            CloseClock,
	}
}

// SetConfiguration apply configuration
// If configuration is diffrent, mark RespondingType as NewParameterValue(0x04)
// It still not applid crc
func (response *GetSettingResponse) SetConfiguration(cog Configuration) (bool, error) {
	original := response.GetConfiguration()

	if original.TimeVerifyMode != cog.TimeVerifyMode {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: TimeVerifyMode: %v -> %v\n", response.TimeVerifyMode, cog.TimeVerifyMode)
		}
		response.TimeVerifyMode = cog.TimeVerifyMode
		response.RespondingType = NewParameterValue
	}

	if original.Speed != cog.Speed {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: Speed: %v -> %v\n", response.Speed, cog.Speed)
		}
		response.Speed = cog.Speed
		response.RespondingType = NewParameterValue
	}

	if original.RecordingCycle != cog.RecordingCycle {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: RecordingCycle: %v -> %v\n", response.RecordingCycle, cog.RecordingCycle)
		}
		response.RecordingCycle = cog.RecordingCycle
		response.RespondingType = NewParameterValue
	}

	if original.UploadCycle != cog.UploadCycle {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: UploadCycle: %v -> %v\n", response.UploadCycle, cog.UploadCycle)
		}
		response.UploadCycle = cog.UploadCycle
		response.RespondingType = NewParameterValue
	}

	if original.EnableFixedTimeUpload != cog.EnableFixedTimeUpload {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: EnableFixedTimeUpload: %v -> %v\n", response.FixedTimeUpload, cog.EnableFixedTimeUpload)
		}
		response.FixedTimeUpload = cog.EnableFixedTimeUpload
		response.RespondingType = NewParameterValue
	}

	// TODO: 업로드 시간 1~4
	// if original.UploadClock != cog.UploadClock {
	// 	response.UploadHour1 = byte(cog.UploadClock.Hour())
	// 	response.UploadMinute1 = byte(cog.UploadClock.Minute())
	// 	response.RespondingType = RespondingTypeNewParameterValue
	// }

	if original.NetworkType != cog.NetworkType {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: NetworkType: %v -> %v\n", response.NetworkType, cog.NetworkType)
		}
		response.NetworkType = cog.NetworkType
		response.RespondingType = NewParameterValue
	}

	if original.DisplayType != cog.DisplayType {
		if EnableDebugMessage {
			fmt.Printf("- configuration changed: DisplayType: %v -> %v\n", response.DisplayType, cog.DisplayType)
		}
		response.DisplayType = cog.DisplayType
		response.RespondingType = NewParameterValue
	}

	if !equalTime(original.SystemTime, cog.SystemTime) {
		if EnableDebugMessage {
			fmt.Printf(
				"- configuration changed: SystemTime: %v-%v-%v %v:%v:%v -> %v-%v-%v %v:%v:%v\n",
				response.Year,
				response.Month,
				response.Day,
				response.Hour,
				response.Minute,
				response.Second,
				cog.SystemTime.Year()%2000,
				cog.SystemTime.Month(),
				cog.SystemTime.Day(),
				cog.SystemTime.Hour(),
				cog.SystemTime.Minute(),
				cog.SystemTime.Second(),
			)
		}
		response.Year = byte(cog.SystemTime.Year() % 2000)
		response.Month = byte(cog.SystemTime.Month())
		response.Day = byte(cog.SystemTime.Day())
		response.Hour = byte(cog.SystemTime.Hour())
		response.Minute = byte(cog.SystemTime.Minute())
		response.Second = byte(cog.SystemTime.Second())
		response.RespondingType = NewParameterValue
	}

	if !equalClockOmitSec(original.OpenClock, cog.OpenClock) {
		if EnableDebugMessage {
			fmt.Printf(
				"- configuration changed: OpenClock: %v:%v -> %v:%v\n",
				response.OpenHour,
				response.OpenMinute,
				cog.OpenClock.Hour(),
				cog.OpenClock.Minute(),
			)
		}
		response.OpenHour = byte(cog.OpenClock.Hour())
		response.OpenMinute = byte(cog.OpenClock.Minute())
		response.RespondingType = NewParameterValue
	}

	if !equalClockOmitSec(original.CloseClock, cog.CloseClock) {
		if EnableDebugMessage {
			fmt.Printf(
				"- configuration changed: CloseClock: %v:%v -> %v:%v\n",
				response.CloseHour,
				response.CloseMinute,
				cog.CloseClock.Hour(),
				cog.CloseClock.Minute(),
			)
		}
		response.CloseHour = byte(cog.CloseClock.Hour())
		response.CloseMinute = byte(cog.CloseClock.Minute())
		response.RespondingType = NewParameterValue
	}

	return false, nil
}

// Binary generate response represneted by binary
//
// Response encode to device with hexencode and result tag.
// For example:
//   resp := fmt.Sprintf("result=%X", bin)
func (response GetSettingResponse) Binary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 58))
	binary.Write(buf, binary.BigEndian, response.RespondingType)
	binary.Write(buf, binary.BigEndian, response.Flag)
	binary.Write(buf, binary.BigEndian, response.SerialNumber)
	binary.Write(buf, binary.BigEndian, response.TimeVerifyMode)
	binary.Write(buf, binary.BigEndian, response.Speed)
	binary.Write(buf, binary.BigEndian, response.RecordingCycle)
	binary.Write(buf, binary.BigEndian, response.UploadCycle)
	binary.Write(buf, binary.BigEndian, response.FixedTimeUpload)
	binary.Write(buf, binary.BigEndian, response.UploadHour1)
	binary.Write(buf, binary.BigEndian, response.UploadMinute1)
	binary.Write(buf, binary.BigEndian, response.UploadHour2)
	binary.Write(buf, binary.BigEndian, response.UploadMinute2)
	binary.Write(buf, binary.BigEndian, response.UploadHour3)
	binary.Write(buf, binary.BigEndian, response.UploadMinute3)
	binary.Write(buf, binary.BigEndian, response.UploadHour4)
	binary.Write(buf, binary.BigEndian, response.UploadMinute4)
	binary.Write(buf, binary.BigEndian, response.NetworkType)
	binary.Write(buf, binary.BigEndian, response.DisplayType)
	binary.Write(buf, binary.BigEndian, response.MacAddress1)
	binary.Write(buf, binary.BigEndian, response.MacAddress2)
	binary.Write(buf, binary.BigEndian, response.MacAddress3)
	binary.Write(buf, binary.BigEndian, response.Year)
	binary.Write(buf, binary.BigEndian, response.Month)
	binary.Write(buf, binary.BigEndian, response.Day)
	binary.Write(buf, binary.BigEndian, response.Hour)
	binary.Write(buf, binary.BigEndian, response.Minute)
	binary.Write(buf, binary.BigEndian, response.Second)
	binary.Write(buf, binary.BigEndian, response.Week)
	binary.Write(buf, binary.BigEndian, response.OpenHour)
	binary.Write(buf, binary.BigEndian, response.OpenMinute)
	binary.Write(buf, binary.BigEndian, response.CloseHour)
	binary.Write(buf, binary.BigEndian, response.CloseMinute)
	binary.Write(buf, binary.BigEndian, response.Reserved1)
	binary.Write(buf, binary.BigEndian, response.Reserved2)

	// eval crc
	crc, err := calcCrc16(buf.Bytes())
	if err != nil {
		return nil, err
	}
	binary.Write(buf, binary.BigEndian, crc)

	if buf.Len() != 58 {
		return nil, errors.New("response length must be 58")
	}

	return buf.Bytes(), err
}

// DeviceStatus represent datus of device
//
// It contain Version, SerialNumber, Focus, Battery...
//
// status=010142AE51520156000D0001E6A7
//   - 0101 Indicates the device firmware version number, version 1.1
//   - 42AE5152 indicates the device SN number, with the lower digit first, that is, The device SN number is 5251AE42
//   - 01 Device retention information
//   - 56 Represents the remaining power of the infrared transmitter battery used, The current remaining power is 86%
//   - 00 Device retention information
//   - 0D Indicates the current remaining capacity of the counter battery, the current remaining capacity is 13%
//
type DeviceStatus struct {
	Version        uint16 // TODO: BigEndian or LittleEndian. Not commented in manual.
	SerialNumber   []byte // Little Endian
	Focus          Focus
	TransmitterBAT byte
	Reserved_1     byte // TODO: WTF
	CounterBAT     byte
	Charge         Charge
	Reserved_2     byte   // TODO: WTF
	Crc16          uint16 // BigEndian
}

func NewDeviceStatus(data string) (*DeviceStatus, error) {
	if len(data) != 28 {
		return nil, fmt.Errorf("failed to parse GetSettingRequest: length must be 53 byte, but came %d byte", len(data))
	}

	var status = new(DeviceStatus)

	data, version, err := readU16(data)
	if err != nil {
		return nil, err
	}
	data, serialNumber, err := readBytes(data, 4)
	if err != nil {
		return nil, err
	}
	data, focus, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, transmitterBattery, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, retention1, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, receiverBattery, err := readU8(data)
	if err != nil {
		return nil, err
	}
	data, charge, err := readU8(data)
	if err != nil {
		return nil, err
	}

	data, retention2, err := readU8(data)
	if err != nil {
		return nil, err
	}

	data, crc, err := readU16(data)
	if err != nil {
		return nil, err
	}

	status.Version = version
	status.SerialNumber = serialNumber
	status.Focus = Focus(focus)
	status.TransmitterBAT = transmitterBattery
	status.Reserved_1 = retention1
	status.CounterBAT = receiverBattery
	status.Charge = Charge(charge)
	status.Reserved_2 = retention2
	status.Crc16 = crc

	return status, nil
}

// calcCrc16 verifies all byte calculation before crc fields(excluding “result=”)
//   - length of `data` must not logner than 78.
//   - 1 byte high 8, 1 byte low 8
//   - result of this function always treated as BigEndian
func calcCrc16(data []byte) (uint16, error) {
	var crc uint16 = 0xFFFF

	if len(data) > 78 {
		return 0, errors.New("length of data must less than 78")
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
	}

	crc = (crc % 0x100) | ((crc / 0x100) << 8)

	return crc, nil
}

type CacheData struct {
	Year    byte
	Month   byte
	Day     byte
	Hour    byte
	Minute  byte
	Secound byte
	Focus   Focus
	DxIn    uint32 // LittleEndian
	DxOut   uint32 // LittleEndian
	Crc16   uint16 // BigEndian
}

func NewCacheData(data []byte) (*CacheData, error) {
	if len(data) != 17 {
		return nil, fmt.Errorf("failed to parse CacheData: length must be 17 byte, but came %d byte", len(data))
	}

	crc, err := calcCrc16(data[:15])
	if err != nil {
		return nil, errors.New("failed to verify crc:" + err.Error())
	}

	if crc != binary.BigEndian.Uint16(data[15:17]) {
		return nil, errors.New("failed to parse CacheData: incorrect crc")
	}

	return &CacheData{
		data[0],
		data[1],
		data[2],
		data[3],
		data[4],
		data[5],
		Focus(data[6]),
		binary.LittleEndian.Uint32(data[7:11]),
		binary.LittleEndian.Uint32(data[11:15]),
		crc,
	}, nil
}

// There are more fields, such as Tend and temp,
// but no description on manual.
type CacheRequest struct {
	Status *DeviceStatus
	Data   []*CacheData
}

func NewCacheRequest(requestSchema *RequestSchema) (*CacheRequest, error) {
	var request = new(CacheRequest)

	if int(requestSchema.Count) != len(requestSchema.Data) {
		return nil, errors.New("failed to parse cache request: count and length of data is not same")
	}

	status, err := NewDeviceStatus(requestSchema.Status)
	if err != nil {
		return nil, errors.New("failed to parse cache request:" + err.Error())
	}

	request.Status = status
	for _, data := range requestSchema.Data {
		cData, err := NewCacheData(data)
		if err != nil {
			return nil, err
		}
		request.Data = append(request.Data, cData)
	}

	return request, nil
}

func (request *CacheRequest) Response(answerType AnswerType, flag uint16, configured Configuration) *CacheResponse {
	return &CacheResponse{
		AnswerType:     answerType,
		Flag:           reverseU16(flag),
		TimeVerifyMode: 0,
		Year:           byte(configured.SystemTime.Year() % 2000),
		Month:          byte(configured.SystemTime.Month()),
		Day:            byte(configured.SystemTime.Day()),
		Hour:           byte(configured.SystemTime.Hour()),
		Minute:         byte(configured.SystemTime.Minute()),
		Second:         byte(configured.SystemTime.Second()),
		Week:           0,
		OpenHour:       byte(configured.OpenClock.Hour()),
		OpenMinute:     byte(configured.OpenClock.Minute()),
		CloseHour:      byte(configured.CloseClock.Hour()),
		CloseMinute:    byte(configured.CloseClock.Minute()),
	}
}

type CacheResponse struct {
	AnswerType     AnswerType
	Flag           uint16
	TimeVerifyMode TimeVerifyMode
	Year           byte
	Month          byte
	Day            byte
	Hour           byte
	Minute         byte
	Second         byte
	Week           byte
	OpenHour       byte
	OpenMinute     byte
	CloseHour      byte
	CloseMinute    byte
	Crc16          uint16
}

func (response *CacheResponse) Binary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 58))
	binary.Write(buf, binary.BigEndian, response.AnswerType)
	binary.Write(buf, binary.BigEndian, response.Flag)
	binary.Write(buf, binary.BigEndian, response.TimeVerifyMode)
	binary.Write(buf, binary.BigEndian, response.Year)
	binary.Write(buf, binary.BigEndian, response.Month)
	binary.Write(buf, binary.BigEndian, response.Day)
	binary.Write(buf, binary.BigEndian, response.Hour)
	binary.Write(buf, binary.BigEndian, response.Minute)
	binary.Write(buf, binary.BigEndian, response.Second)
	binary.Write(buf, binary.BigEndian, response.Week)
	binary.Write(buf, binary.BigEndian, response.OpenHour)
	binary.Write(buf, binary.BigEndian, response.OpenMinute)
	binary.Write(buf, binary.BigEndian, response.CloseHour)
	binary.Write(buf, binary.BigEndian, response.CloseMinute)

	// eval crc
	crc, err := calcCrc16(buf.Bytes())
	if err != nil {
		return nil, err
	}
	binary.Write(buf, binary.BigEndian, crc)

	return buf.Bytes(), err
}
