// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

// variable for count
var (
	counter = hpc015.Counter(2)
)

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

// run http server
func main() {
	log.Println("- server is running on:", server_host+handler_path)
	http.HandleFunc(handler_path, hpc015Handler)
	log.Fatal(http.ListenAndServe(server_host, nil))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	log.Println("income!")

	reqBuf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("! failed to read request:", err.Error())
		return
	}
	fmt.Println("- request: " + string(reqBuf))

	resp, err := http.Post("http://192.168.0.52:8900/dataport", "", bytes.NewBuffer(reqBuf))
	if err != nil {
		log.Println("! failed to bypass:", err.Error())
		return
	}

	respBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("! failed to read body:", err.Error())
		return
	}

	fmt.Println("- response: " + string(respBuf))
	w.Write(respBuf)

	/* 요청 분석 */
	func() {
		requestSchema, err := hpc015.NewRequestSchema(string(reqBuf))
		if err != nil {
			log.Println("\t! failed to parse RequestSchema:", err.Error())
			return
		}

		switch requestSchema.Cmd {
		case "getsetting":
			/* 설정 요청 */
			setReq, err := hpc015.NewSettingRequest(requestSchema.Data[0])
			if err != nil {
				log.Println("\t! failed to parse SettingRequest:", err.Error())
				return
			}

			// new response based on request
			setResp := setReq.Response(requestSchema.Flag)

			// (optional) get current configuration
			conf := setResp.GetConfiguration()
			fmt.Printf("- current TimeVerifyMode: %v\n", conf.TimeVerifyMode)
			fmt.Printf("- current Speed: %v\n", conf.Speed)
			fmt.Printf("- current RecordingCycle: %v\n", conf.RecordingCycle)
			fmt.Printf("- current UploadCycle: %v\n", conf.UploadCycle)
			fmt.Printf("- current EnableFixedTimeUpload: %v\n", conf.EnableFixedTimeUpload)
			fmt.Printf("- current UploadClock: %v\n", conf.UploadClock)
			fmt.Printf("- current NetworkType: %v\n", conf.NetworkType)
			fmt.Printf("- current DisplayType: %v\n", conf.DisplayType)
			fmt.Printf("- current SystemTime: %v\n", conf.SystemTime)
			fmt.Printf("- current OpenClock: %v\n", conf.OpenClock)
			fmt.Printf("- current CloseClock: %v\n", conf.CloseClock)

		case "cache":
			/* 캐시 요청 */
			cacheReq, err := hpc015.NewCacheRequest(requestSchema)
			if err != nil {
				log.Println("! failed to parse CacheRequest:", err.Error())
				return
			}

			//
			cacheResp := cacheReq.Response(hpc015.OK, requestSchema.Flag, obtainCog())

			bin, err := cacheResp.Binary()
			if err != nil {
				log.Println("! failed to convert binary:", err.Error())
				return
			}
			_ = bin
		}
	}()
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
