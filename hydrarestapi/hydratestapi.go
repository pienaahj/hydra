package hydratestapi

import (
	"log"
	"net/http"

	hydraConfigurator "github.com/pienaahj/hydra/hydraconfigurator"
)

type DBlayerconfig struct {
	DB   string `json:"database"`
	Conn string `json:"connectionstring"`
}

func InitializeAPIHandlers() error {
	conf := new(DBlayerconfig)
	err := hydraConfigurator.GetConfiguration(hydraConfigurator.JSON, conf, "/hydraweb/apiconfig.json")
	if err != nil {
		log.Println("Error decoding JSON, err")
		return err
	}
	h := NewhydraCrewReqHandler()
	err = h.connect(conf.DB, conf.Conn)
	if err != nil {
		log.Println("Error connecting to db ", err)
		return err
	}
	http.HandleFunc("/hydracrew/", h.handleHydraCrewRequests)
	// http.handle("/hydracrew/", h) if using ServeHTTP
	return nil
}

func RunAPI() error {
	if err := InitializeAPIHandlers(); err != nil {
		return err
	}
	return http.ListenAndServe(":8061", nil)
}
