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

// same path as you configured on
// 192.168.8.1 -> SET NET -> SERVER
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
		settingRequest, _ := hpc015.NewSettingRequest(requestSchema.Data[0])

		// new response based on request
		settingResponse := hpc015.NewSettingResponse(settingRequest, requestSchema.Flag)

		// get current configuration
		_ = settingResponse.GetConfiguration()

		// apply new configuration
		settingResponse.SetConfiguration(obtainCog())

		// response
		bin, err := settingResponse.Binary()
		if err != nil {
			log.Println("failed to convert binary:", err.Error())
			return
		}
		resp := fmt.Sprintf("result=%X", bin)
		log.Println("< response with:", resp)
		w.Write([]byte(resp))
		return

	case "cache":
		// we can use Status within cache command
		_ = requestSchema.Status

		// cache can hold multiple data fields
	}
}
