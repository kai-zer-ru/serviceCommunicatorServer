package serviceCommunicatorServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DaemonData struct {
	ServerAddress string
	Daemon        struct {
		DaemonName    string `json:"name"`
		DaemonAddress string `json:"address"`
		Description   string `json:"description"`
	}
}

func (d *DaemonData) Register() {
	go func() {
		time.Sleep(5 * time.Second)
		data, err := json.Marshal(d.Daemon)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		resp, err := http.Post(d.ServerAddress+"/registerService", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
			return
		}
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(responseData))
		_ = resp.Body.Close()
		return
	}()
}
