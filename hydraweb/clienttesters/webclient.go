package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// url := "https://eojrcxous026ckj.m.pipedream.net"
	/*
		resp, err := http.Get(url)
		inspectResponse(resp, err)

		//  create a standard json data piece
		data, err := json.Marshal(struct {
			X int
			Y float32
		}{X: 4, Y: 3.8})
		if err != nil {
			log.Fatal("Error occured while marshalling json ", err)
		}
		resp, err = http.Post(url, "application/json", bytes.NewReader(data))
		inspectResponse(resp, err)

		client := http.Client{
			Timeout: 3 * time.Minute,
		}
		client.Get(url)

		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("x-testheader", "learning go header")
		req.Header.Set("User-Agent", "Go learning Http/1.1")
		resp, err := client.Do(req)
		inspectResponse(resp, err)
	*/

	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	v := struct {
		IP string `json:"ip,omitempty"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v.IP)
}

// create a generic way to inspect a response
func inspectResponse(resp *http.Response, err error) {
	if err != nil {
		log.Fatal("Error occured while marshalling json ", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error occured while trying to read response body ", err)
	}
	log.Println(string(b))
}
