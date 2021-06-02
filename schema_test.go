// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpc015

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestIntegrateSetting(t *testing.T) {

	var input = "cmd=getsetting&flag=022E&data=0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F010100022E010000173B80C0"
	requestSchema, err := NewRequestSchema(input)
	if err != nil {
		t.Errorf("NewRequestSchema() error = %v", err)
		t.FailNow()
		return
	}

	setReq, err := NewSettingRequest(requestSchema.Data[0])
	if err != nil {
		t.Errorf("NewSettingRequest() error = %v", err)
		t.FailNow()
	}

	_ = setReq

}

func TestNewRequest(t *testing.T) {
	type args struct {
		reqestString string
	}
	tests := []struct {
		name    string
		args    args
		want    *RequestSchema
		wantErr bool
	}{
		{
			name: "NewRequest(1)",
			args: args{
				"cmd=getsetting&flag=022E&data=0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F010100022E010000173B80C0",
			},
			want: &RequestSchema{
				Cmd:  "getsetting",
				Flag: 0x022E,
				Data: [][]byte{{0xD, 0x3B, 0xB3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC, 0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC, 0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA, 0xF, 0x1, 0x1, 0x0, 0x2, 0x2E, 0x1, 0x0, 0x0, 0x17, 0x3B, 0x80, 0xC0}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRequestSchema(tt.args.reqestString)
			pretty.Println("got: ", got)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSettingRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *GetSettingRequest
		wantErr bool
	}{
		{
			name: "NewGetSettingRequest(1)",
			args: args{
				data: []byte{
					0xD, 0x3B, 0xB3, 0x82,
					0x3,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x0,
					0x2,
					0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC,
					0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC,
					0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA,
					0xF,
					0x1,
					0x1,
					0x0,
					0x2,
					0x2E,
					0x1,
					0x0,
					0x0,
					0x17,
					0x3B,
					0x80, 0xC0,
				},
			},
			want: &GetSettingRequest{
				SerialNumber:    []byte{0xD, 0x3B, 0xB3, 0x82},
				TimeVerifyMode:  0x3,
				Speed:           0x0,
				RecordingCycle:  0x0,
				UploadCycle:     0x0,
				FixedTimeUpload: 0x0,
				UploadHour1:     0x0,
				UploadMinute1:   0x0,
				UploadHour2:     0x0,
				UploadMinute2:   0x0,
				UploadHour3:     0x0,
				UploadMinute3:   0x0,
				UploadHour4:     0x0,
				UploadMinute4:   0x0,
				NetworkType:     0x0,
				DisplayType:     0x2,
				MacAddress1:     []byte{0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress2:     []byte{0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress3:     []byte{0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA},
				Year:            0xF,
				Month:           0x1,
				Day:             0x1,
				Hour:            0x0,
				Minute:          0x2,
				Second:          0x2E,
				Week:            0x1,
				OpenHour:        0x0,
				OpenMinute:      0x0,
				CloseHour:       0x17,
				CloseMinute:     0x3B,
				Crc16:           binary.BigEndian.Uint16([]byte{0x80, 0xC0}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSettingRequest(tt.args.data)
			pretty.Println("got: ", got)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSettingRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettingRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSettingResponse_toBindary(t *testing.T) {
	type fields struct {
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
		NetworkType     NetworkType
		DisplayMode     byte
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
		Reserved1       byte
		Reserved2       byte
		Crc16           uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "",
			fields: fields{
				RespondingType:  0x04,
				Flag:            binary.BigEndian.Uint16([]byte{0, 0}),
				SerialNumber:    []byte{0, 0, 0, 0},
				CommandType:     0x03,
				Speed:           0x00,
				RecordingCycle:  0x0A,
				UploadCycle:     0x78,
				FixedTimeUpload: 0x00,
				UploadHour1:     0x00,
				UploadMinute1:   0x00,
				UploadHour2:     0x00,
				UploadMinute2:   0x00,
				UploadHour3:     0x00,
				UploadMinute3:   0x00,
				UploadHour4:     0x00,
				UploadMinute4:   0x00,
				NetworkType:     0x00,
				DisplayMode:     0x02,
				MacAddress1:     []byte{0, 0, 0, 0, 0, 0, 0},
				MacAddress2:     []byte{0, 0, 0, 0, 0, 0, 0},
				MacAddress3:     []byte{0, 0, 0, 0, 0, 0, 0},
				Year:            0x11,
				Month:           0x03,
				Day:             0x05,
				Hour:            0x11,
				Minute:          0x11,
				Secound:         0x26,
				Week:            0x00,
				OpenHour:        0x0A,
				OpenMinute:      0x00,
				CloseHour:       0x14,
				CloseMinute:     0x1E,
			},
			want: []uint8{0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0xa, 0x78, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x3, 0x5, 0x11, 0x11, 0x26, 0x0, 0xa, 0x0, 0x14, 0x1e, 0x0, 0x0, 0x93, 0x51},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := GetSettingResponse{
				RespondingType:  tt.fields.RespondingType,
				Flag:            tt.fields.Flag,
				SerialNumber:    tt.fields.SerialNumber,
				TimeVerifyMode:  TimeVerifyMode(tt.fields.CommandType),
				Speed:           Speed(tt.fields.Speed),
				RecordingCycle:  tt.fields.RecordingCycle,
				UploadCycle:     tt.fields.UploadCycle,
				FixedTimeUpload: tt.fields.FixedTimeUpload,
				UploadHour1:     tt.fields.UploadHour1,
				UploadMinute1:   tt.fields.UploadMinute1,
				UploadHour2:     tt.fields.UploadHour2,
				UploadMinute2:   tt.fields.UploadMinute2,
				UploadHour3:     tt.fields.UploadHour3,
				UploadMinute3:   tt.fields.UploadMinute3,
				UploadHour4:     tt.fields.UploadHour4,
				UploadMinute4:   tt.fields.UploadMinute4,
				NetworkType:     tt.fields.NetworkType,
				DisplayType:     DisplayType(tt.fields.DisplayMode),
				MacAddress1:     tt.fields.MacAddress1,
				MacAddress2:     tt.fields.MacAddress2,
				MacAddress3:     tt.fields.MacAddress3,
				Year:            tt.fields.Year,
				Month:           tt.fields.Month,
				Day:             tt.fields.Day,
				Hour:            tt.fields.Hour,
				Minute:          tt.fields.Minute,
				Second:          tt.fields.Secound,
				Week:            tt.fields.Week,
				OpenHour:        tt.fields.OpenHour,
				OpenMinute:      tt.fields.OpenMinute,
				CloseHour:       tt.fields.CloseHour,
				CloseMinute:     tt.fields.CloseMinute,
				Reserved1:       tt.fields.Reserved1,
				Reserved2:       tt.fields.Reserved2,
				Crc16:           tt.fields.Crc16,
			}
			if got, err := request.Binary(); !reflect.DeepEqual(got, tt.want) {
				if err != nil {
					t.Errorf("GetSettingResponse.toBindary() = %v, want %v", got, tt.want)
				}
				t.Errorf("GetSettingResponse.toBindary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeviceStatus(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *DeviceStatus
		wantErr bool
	}{
		{
			name: "TestNewDeviceStatus",
			args: args{"010142AE51520156000D0001E6A7"},
			want: &DeviceStatus{
				Version:        0x0101,
				SerialNumber:   0x42AE5152,
				Focus:          0x01,
				TransmitterBAT: 0x56,
				Reserved_1:     0x00,
				CounterBAT:     0x0D,
				Charge:         0x00,
				Reserved_2:     0x01,
				Crc16:          0xE6A7,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeviceStatus(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeviceStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeviceStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

// freeze
func TestCalcCrc16(t *testing.T) {
	type TestCase struct {
		input  []byte
		output uint16
	}
	var tests = []TestCase{
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x2, 0x1, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]byte{0xEC, 0xE4}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x23, 0x1, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0x5D, 0xE2}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0xa, 0x1, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0xA4, 0xE5}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x8, 0x1, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0x46, 0xE4}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x2, 0x1, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0xec, 0xe4}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x1e, 0x28, 0x0, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0xe7, 0x20}),
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0x15, 0x6, 0x2, 0xe, 0x1, 0x2b, 0x0, 0x0, 0x0, 0x17, 0x3b},
			binary.BigEndian.Uint16([]uint8{0x6b, 0x46}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x2a, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0xe9, 0x7e}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x2c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0xc6, 0x5e}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x2d, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0x33, 0xcf}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x2e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0xc, 0xff}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x30, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0x5c, 0x5f}),
		},
		{
			[]uint8{0x15, 0x5, 0xd, 0xd, 0x33, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0xa9, 0xce}),
		},

		{
			[]uint8{0x5, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf, 0x1, 0x1, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x17, 0x3b, 0x0, 0x0},
			binary.BigEndian.Uint16([]uint8{0x1c, 0xe4}),
		},
	}

	for _, tc := range tests {
		res, err := calcCrc16(tc.input)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if res != tc.output {
			t.Error("not equal")
			t.FailNow()
		}
	}
}

// use to make test case
func TestHexStringToGoLiteral(t *testing.T) {
	var inputs = []string{

		"050200000000000300000000000000000000000000000000000000000000000000000000000000000000000F0101000002000000173B00001CE4",
	}

	for _, tt := range inputs {
		t.Run("", func(t *testing.T) {
			val, err := hex.DecodeString(tt)
			if err != nil {
				t.Fatal("what happened")
			} else {
				log.Println(pretty.Sprint(val))
			}
		})
	}
}

func TestParsingSettingRequests(t *testing.T) {
	var inputs = []string{
		//                             |       c s r u f u u u u u u u u   d                                                         o o c c
		//                             |       m p c p i h m h m h m h m m s                                                         p p l l
		//                             |serial|d|d|c|c|x|1|1|2|2|3|3|4|4|d|p|mac_addr1    |mac_addr2    |mac_addr3    |Y|M|D|H|M|S|W|h|m|h|m|crc|
		"cmd=getsetting&flag=0002&data=0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F0101000002010000173BECE4",
		"cmd=getsetting&flag=1E28&data=0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F0101001E28000000173BE720",
		"cmd=getsetting&flag=012B&data=0D3BB382030000000000000000000000000000085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA1506020E012B000000173B6B46",
	}

	for _, input := range inputs {
		fmt.Println()
		schema, err := NewRequestSchema(input)
		if err != nil {
			t.Fatal("NewRequestSchema() failed:", err)
		}
		_ = schema

		setReq, err := NewSettingRequest(schema.Data[0])
		if err != nil {
			log.Println("\t! failed to parse SettingRequest:", err.Error())
			return
		}

		// new response based on request
		setResp := setReq.Response(schema.Flag)

		// (optional) get current configuration
		conf := setResp.GetConfiguration()
		fmt.Printf("- current systemtime: %v\n", *&conf.SystemTime)
		fmt.Printf("- current recording cycle: %d\n", conf.RecordingCycle)
		fmt.Printf("- current uploading cycle: %d\n", conf.UploadCycle)
		fmt.Printf("- current EnableFixedTimeUpload: %d\n", conf.EnableFixedTimeUpload)
		fmt.Printf("- current CloseClock: %v\n", conf.CloseClock)
		fmt.Printf("- current OpenClock: %v\n", conf.OpenClock)
	}

}
