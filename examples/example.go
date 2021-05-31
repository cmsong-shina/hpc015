// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cmsong-shina/hpc015"
)

// same path as you configured on
// 192.168.8.1 -> SET NET -> SERVER
const (
	port = ":8888"
	path = "/cs"
)

// set your own configuration
func obtainCog() hpc015.Configuration {
	return hpc015.Configuration{
		// CommandType
		// Speed
		// RecordingCycle
		// UploadCycle
		// FixedTimeUpload
		// UploadClock
		// Model
		// DisableType
		// SystemTime
		// OpenClock
		// CloseClock
	}
}

// run http server
func main() {
	http.HandleFunc(path, hpc015Handler)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handler
func hpc015Handler(w http.ResponseWriter, req *http.Request) {
	bin, _ := ioutil.ReadAll(req.Body)

	requestSchema, _ := hpc015.NewRequestSchema(string(bin))

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
		w.Write(settingResponse.Binary())
		return

	case "cache":
		// we can use Status within cache command
		_ = requestSchema.Status

		// cache can hold multiple data fields
	}
}
