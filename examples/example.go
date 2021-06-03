// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/cmsong-shina/hpc015"
)

// Same path as you configured on 192.168.8.1(device AP mode) -> SET NET -> SERVER
// And you MUST set path, not only hostname(from my own experience)
const (
	server_host  = ":8888"
	handler_path = "/cs"
	count_path   = handler_path + "/count"
)

// variable for count
var (
	counter = hpc015.Counter(0)
)

// run http server
func main() {
	log.Println("- server is running on:", server_host+handler_path)

	http.HandleFunc(handler_path, hpc015Handler) // handle hpc015
	http.HandleFunc(count_path, count_handler)   // handle set/get count

	log.Fatal(http.ListenAndServe(server_host, nil))
}

// handler
func hpc015Handler(w http.ResponseWriter, req *http.Request) {
	bin, _ := ioutil.ReadAll(req.Body)

	fmt.Println()
	log.Println("> request from:", req.RemoteAddr, string(bin))
	requestSchema, err := hpc015.NewRequestSchema(string(bin))
	if err != nil {
		log.Println("! failed to parse RequestSchema:", err.Error())
		return
	}

	switch requestSchema.Cmd {
	case "getsetting":
		// getsetting has one data field
		setReq, err := hpc015.NewSettingRequest(requestSchema.Data[0])
		if err != nil {
			log.Println("! failed to parse SettingRequest:", err.Error())
			return
		}

		// new response based on request
		setResp := setReq.Response(requestSchema.Flag)

		// (optional) get current configuration
		conf := setResp.GetConfiguration()
		fmt.Printf("- old TimeVerifyMode: %v\n", conf.TimeVerifyMode)
		fmt.Printf("- old Speed: %v\n", conf.Speed)
		fmt.Printf("- old RecordingCycle: %v\n", conf.RecordingCycle)
		fmt.Printf("- old UploadCycle: %v\n", conf.UploadCycle)
		fmt.Printf("- old EnableFixedTimeUpload: %v\n", conf.EnableFixedTimeUpload)
		fmt.Printf("- old UploadClock: %v\n", conf.UploadClock)
		fmt.Printf("- old NetworkType: %v\n", conf.NetworkType)
		fmt.Printf("- old DisplayType: %v\n", conf.DisplayType)
		fmt.Printf("- old SystemTime: %v\n", conf.SystemTime)
		fmt.Printf("- old OpenClock: %v\n", conf.OpenClock)
		fmt.Printf("- old CloseClock: %v\n", conf.CloseClock)

		{
			// (optional) apply new configuration
			//
			// DO NOT change configuration every time.
			// When you modify configuration, device send request to confirmation,
			// and if you change system time(for example) again, device send confirmation again. It is loop.
			//
			// in this case, we apply configuration when systemtime difference more than 10 minutes.

			duration := conf.SystemTime.Sub(time.Now())
			if math.Abs(duration.Minutes()) > 5 {
				conf.SystemTime = time.Now()
			}
			setResp.SetConfiguration(*conf)
		}

		// response
		bin, err := setResp.Binary()
		if err != nil {
			log.Println("! failed to convert binary:", err.Error())
			return
		}
		resp := fmt.Sprintf("result=%X", bin)
		log.Println("< response with:", resp)
		w.Write([]byte(resp))
		return

	case "cache":
		// device will send cache request when they got respose about getsetting correctly

		// create cache request
		cacheReq, err := hpc015.NewCacheRequest(requestSchema)
		if err != nil {
			log.Println("! failed to parse CacheRequest:", err.Error())
			return
		}

		// create cache response
		cacheResp := cacheReq.Response(hpc015.OK, requestSchema.Flag, obtainCog())

		// send cache response
		bin, err := cacheResp.Binary()
		if err != nil {
			log.Println("! failed to convert binary:", err.Error())
			return
		}
		resp := fmt.Sprintf("result=%X", bin)
		_, err = w.Write([]byte(resp))
		if err != nil {
			log.Println("! failed to send cache response:", err.Error())
			return
		}
		log.Println("< response with:", resp)

		// if send successfully, process event
		for _, data := range cacheReq.Data {
			counter.Count(data)
		}
		fmt.Println("--------current:", counter.Get())
		return
	}
}

func count_handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	switch req.Method {
	case http.MethodGet:
		w.Write([]byte(strconv.FormatInt(int64(counter.Get()), 10)))

	case http.MethodPost:
		bin, _ := ioutil.ReadAll(req.Body)

		i, err := strconv.ParseInt(string(bin), 10, 64)
		if err == nil {
			counter.Set(int(i))
		}
		log.Println("> request set to:", i)
	}
}

// Implement your own configuration provider.
// When you write Clock, mind not to set Year/Month/Day as 0.
func obtainCog() hpc015.Configuration {
	return hpc015.Configuration{
		TimeVerifyMode:        hpc015.Exclude,
		Speed:                 hpc015.Low,
		RecordingCycle:        0,
		UploadCycle:           0,
		EnableFixedTimeUpload: 0,
		NetworkType:           hpc015.Online,
		DisplayType:           hpc015.All,
		SystemTime:            time.Now(),
		OpenClock:             time.Date(1, 1, 1, 0, 0, 0, 0, time.Local),
		CloseClock:            time.Date(1, 1, 1, 23, 59, 0, 0, time.Local),
	}
}
