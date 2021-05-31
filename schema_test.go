// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpc015

import (
	"encoding/hex"
	"log"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

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
				CommandType:     0x3,
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
				Model:           0x0,
				DisableType:     0x2,
				MacAddress1:     []byte{0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress2:     []byte{0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress3:     []byte{0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA},
				Year:            0xF,
				Month:           0x1,
				Day:             0x1,
				Hour:            0x0,
				Minute:          0x2,
				Secound:         0x2E,
				Week:            0x1,
				OpenHour:        0x0,
				OpenMinute:      0x0,
				CloseHour:       0x17,
				CloseMinute:     0x3B,
				Crc16:           0xC080,
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

func TestNewSettingResponse(t *testing.T) {
	type args struct {
		request *GetSettingRequest
		flag    uint16
	}
	tests := []struct {
		name string
		args args
		want *GetSettingResponse
	}{
		{
			name: "",
			args: args{
				&GetSettingRequest{
					SerialNumber:    []byte{0xD, 0x3B, 0xB3, 0x82},
					CommandType:     0x3,
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
					Model:           0x0,
					DisableType:     0x2,
					MacAddress1:     []byte{0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
					MacAddress2:     []byte{0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
					MacAddress3:     []byte{0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA},
					Year:            0xF,
					Month:           0x1,
					Day:             0x1,
					Hour:            0x0,
					Minute:          0x2,
					Secound:         0x2E,
					Week:            0x1,
					OpenHour:        0x0,
					OpenMinute:      0x0,
					CloseHour:       0x17,
					CloseMinute:     0x3B,
					Crc16:           0xC080,
				},
				0x022E,
			},
			want: &GetSettingResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSettingResponse(tt.args.request, tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettingResponse() = %v, want %v", got, tt.want)
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
				Flag:            0x12AB,
				SerialNumber:    []byte{0xD, 0x3B, 0xB3, 0x82},
				CommandType:     0x3,
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
				Model:           0x0,
				DisableType:     0x2,
				MacAddress1:     []byte{0x8, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress2:     []byte{0xA, 0x5D, 0xDD, 0x5A, 0x75, 0xCB, 0xDC},
				MacAddress3:     []byte{0x90, 0x9F, 0x33, 0x17, 0x3C, 0xE4, 0xDA},
				Year:            0xF,
				Month:           0x1,
				Day:             0x1,
				Hour:            0x0,
				Minute:          0x2,
				Secound:         0x2E,
				Week:            0x1,
				OpenHour:        0x0,
				OpenMinute:      0x0,
				CloseHour:       0x17,
				CloseMinute:     0x3B,
				Crc16:           0xC080,
			},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := GetSettingResponse{
				RespondingType:  tt.fields.RespondingType,
				Flag:            tt.fields.Flag,
				SerialNumber:    tt.fields.SerialNumber,
				CommandType:     tt.fields.CommandType,
				Speed:           tt.fields.Speed,
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
				Model:           tt.fields.Model,
				DisableType:     tt.fields.DisableType,
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
				Carge:          0x00,
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

func TestCalcCrc16(t *testing.T) {
	type TestCase struct {
		input  []byte
		output uint16
	}
	var tests = []TestCase{
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x2, 0x1, 0x0, 0x0, 0x17, 0x3b},
			0xECE4,
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x23, 0x1, 0x0, 0x0, 0x17, 0x3b},
			0x5DE2,
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0xa, 0x1, 0x0, 0x0, 0x17, 0x3b},
			0xA4E5,
		},
		{
			[]uint8{0xd, 0x3b, 0xb3, 0x82, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x8, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0xa, 0x5d, 0xdd, 0x5a, 0x75, 0xcb, 0xdc, 0x90, 0x9f, 0x33, 0x17, 0x3c, 0xe4, 0xda, 0xf, 0x1, 0x1, 0x0, 0x0, 0x8, 0x1, 0x0, 0x0, 0x17, 0x3b},
			0x46E4,
		},
	}

	for _, tc := range tests {
		res, err := calcCrc16(tc.input)
		if err != nil {
			t.FailNow()
		}
		if res != tc.output {
			t.FailNow()
		}
	}
}

// use to make test case
func TestHexStringToGoLiteral(t *testing.T) {
	var inputs = []string{
		"0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F0101000002010000173BECE4",
		"0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F0101000023010000173B5DE2",
		"0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F010100000A010000173BA4E5",
		"0D3BB382030000000000000000000000000002085DDD5A75CBDC0A5DDD5A75CBDC909F33173CE4DA0F0101000008010000173B46E4",
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
