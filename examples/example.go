// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cmsong-shina/hpc015"
)

// Same path as you configured on 192.168.8.1(device AP mode) -> SET NET -> SERVER
// And you MUST set path, not only hostname(from my own experience)
const (
	server_host  = ":8888"
	handler_path = "/cs"
)

// set your own configuration
func obtainCog() hpc015.Configuration {
	now := time.Now()

	return hpc015.Configuration{
		// CommandType
		// Speed
		// RecordingCycle
		// UploadCycle
		// FixedTimeUpload
		// UploadClock
		// Model
		// DisableType
		SystemTime: &now,
		// OpenClock
		// CloseClock
	}
}

// run http server
func main() {
	log.Println("server is running on:", server_host+handler_path)
	http.HandleFunc(handler_path, hpc015Handler)
	log.Fatal(http.ListenAndServe(server_host, nil))
}

// handler
func hpc015Handler(w http.ResponseWriter, req *http.Request) {
	bin, _ := ioutil.ReadAll(req.Body)

	requestSchema, _ := hpc015.NewRequestSchema(string(bin))
	log.Println("> request from:", req.RemoteAddr, string(bin))

	switch requestSchema.Cmd {
	case "getsetting":
		// getsetting has one data field
		setReq, _ := hpc015.NewSettingRequest(requestSchema.Data[0])

		// new response based on request
		setResp := setReq.Response(requestSchema.Flag)

		// (optional) get current configuration
		_ = setResp.GetConfiguration()

		// (optional) apply new configuration
		//
		// DO NOT change configuration every time.
		// When you modify configuration, device send request to confirmation,
		// and if you change system time(for example) again, device send confirmation again. It is loop.
		//   setResp.SetConfiguration(obtainCog())

		// response
		bin, err := setResp.Binary()
		if err != nil {
			log.Println("failed to convert binary:", err.Error())
			return
		}
		resp := fmt.Sprintf("result=%X", bin)
		log.Println("< response with:", resp)
		w.Write([]byte(resp))
		return

	case "cache":
		// device will send cache request when they got respose about getsetting correctly

		// we can use Status within cache command
		_ = requestSchema.Status

		// cache can hold multiple data fields
	}
}
