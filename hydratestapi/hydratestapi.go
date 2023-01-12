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

func InitializeAPIHandlers() {
	conf := new(DBlayerconfig)
	err := hydraConfigurator.GetConfiguration(hydraConfigurator.JSON, conf, "/hydraweb/apiconfig.json")
	if err != nil {
		log.Fatal("Error decoding JSON, err")
	}
	h := NewhydraCrewReqHandler()
	err = h.connect(conf.DB, conf.Conn)
	if err != nil {
		log.Fatal("Error connecting to db ", err)
	}
	http.HandleFunc("/hydracrew/", h.handleHydraCrewRequests)
	// http.handle("/hydracrew/", h) if using ServeHTTP
}

func RunAPI() {
	InitializeAPIHandlers()
	http.ListenAndServe(":8061", nil)
}
