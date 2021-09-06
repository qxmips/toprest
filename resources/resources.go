package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
)

type Resource struct {
	Totalram  uint64 `json:"Totalram"`
	Freeram   uint64 `json:"Freeram"`
	Totalswap uint64 `json:"Totalswap"`
	Freeswap  uint64 `json:"Freeswap"`
}

func GetResources(rw http.ResponseWriter, r *http.Request) {
	info := syscall.Sysinfo_t{}
	err := syscall.Sysinfo(&info)
	if err != nil {
		fmt.Println("Error:", err)
	}

	p := &Resource{Totalram: info.Totalram / 1024 / 1024,
		Freeram:   info.Freeram / 1024 / 1024,
		Totalswap: info.Totalswap / 1024 / 1024,
		Freeswap:  info.Freeswap / 1024 / 1024}

	rw.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(rw)
	e.Encode(p)

}
