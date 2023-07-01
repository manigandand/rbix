package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/v3/host"
)

func main() {
	port := ":8888"
	http.HandleFunc("/", handler)
	log.Println("rbi listening on ", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostInfo, err := host.Info()
	if err != nil {
		log.Println("[host.Info() Error] %v", err.Error())
	}
	log.Println("host info:", hostInfo)

	res := map[string]interface{}{
		"data": map[string]interface{}{
			"host":      r.Host,
			"host_info": hostInfo,
			"message":   "you're in safe zone with rbi",
		},
		"meta": map[string]interface{}{
			"status":      "ok",
			"status_code": http.StatusOK,
		},
	}

	resp, err := json.Marshal(&res)
	if err != nil {
		log.Println("error:", err.Error())
		return
	}

	w.Header().Set("Host", r.Host)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Fatal(err.Error())
	}
}
