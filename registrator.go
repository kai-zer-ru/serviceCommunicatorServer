package serviceCommunicatorServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DaemonData struct {
	ServerAddress string
	Daemon        struct {
		DaemonName    string        `json:"name"`
		DaemonAddress string        `json:"address"`
		Description   string        `json:"description"`
		Commands      CommandStruct `json:"commands"`
	}
}

func (d *DaemonData) Register() error {
	data, err := json.Marshal(d.Daemon)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	resp, err := http.Post(d.ServerAddress+"/registerService", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(responseData))
	return resp.Body.Close()
}
