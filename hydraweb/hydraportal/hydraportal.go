package hydraportal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	hydraConfigurator "github.com/pienaahj/hydra/hydraconfigurator"
	"github.com/pienaahj/hydra/hydradblayer"
	hydratestapi "github.com/pienaahj/hydra/hydraweb/hydrarestapi"
	"github.com/spf13/viper"
)

// create the template
var hydraWebTemplate *template.Template

/*
	func findDirFiles() {
		files, err := os.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
*/
func Run() error {
	var err error
	// findDirFiles()
	//  cash the templates
	hydraWebTemplate, err = template.ParseFiles("./hydraweb/hydraportal/cover/Crew/crew.html", "./hydraweb/hydraportal/cover/about/about.html")
	if err != nil {
		return err
	}
	// make the deployment folder configurable
	conf := struct {
		Filespath string `json:"filespath"`
	}{}
	err = hydraConfigurator.GetConfiguration(hydraConfigurator.JSON, &conf, "./hydraweb/portalconfig.json")
	if err != nil {
		return err
	}

	fmt.Println("Configuration obtained...")
	hydratestapi.InitializeAPIHandlers()
	log.Println(http.Dir(conf.Filespath))
	fs := http.FileServer(http.Dir(conf.Filespath))
	//  creates a static web page
	http.Handle("/", fs)
	http.HandleFunc("/Crew/", crewhandler)
	http.HandleFunc("/about/", abouthandler)
	return http.ListenAndServe(":8061", nil)
}
func loadcredetials() (string, error) {
	/*
		MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
		MYSQL_DATABASE: ${DB_NAME}
		MYSQL_USER: ${DB_USER}
		MYSQL_PASSWORD: ${DB_PASSWORD}
	*/
	var connectionString string
	var connectionStringU string
	var connectionStringP string
	viper.AutomaticEnv()
	if connectionStringT, ok := viper.Get("DB_USER").(string); !ok {
		log.Println("Cannot get user env variable")
		return "", fmt.Errorf("env error")
	} else {
		connectionStringU = connectionStringT
		// fmt.Printf("viper : %s = %s \n", "Connection string", connectionString)
	}
	if connectionStringT, ok := viper.Get("DB_PASSWORD").(string); !ok {
		log.Println("Cannot get passw env variable")
		return "", fmt.Errorf("env error")
	} else {
		connectionStringP = connectionStringT
		// fmt.Printf("viper : %s = %s \n", "Connection string", connectionString)
	}

	connectionString = connectionStringU + ":" + connectionStringP
	return connectionString, nil
}

// BuildConnectionString builds a connection string
func buildConnectionString(connStr string) string {
	return connStr + "@tcp(127.0.0.1:3306)/Hydra?parseTime=true"
}

func crewhandler(w http.ResponseWriter, r *http.Request) {
	connStr, err := loadcredetials()
	if err != nil {
		log.Println("Error loading credentials")
		return
	}
	fmt.Println("Credentials loaded...")

	connStr = buildConnectionString(connStr)
	dblayer, err := hydradblayer.ConnectDatabase("mysql", connStr)
	if err != nil {
		log.Println("Error while connecting to db")
		return
	}
	all, err := dblayer.AllMembers()
	if err != nil {
		return
	}

	fmt.Println(all)
	err = hydraWebTemplate.ExecuteTemplate(w, "crew.html", all)
	if err != nil {
		log.Println(err)
	}
}

func abouthandler(w http.ResponseWriter, r *http.Request) {
	about := struct {
		Msg string `json:"message"`
	}{}
	err := hydraConfigurator.GetConfiguration(hydraConfigurator.JSON, &about, "./hydraweb/about.json")
	if err != nil {
		return
	}
	err = hydraWebTemplate.ExecuteTemplate(w, "about.html", about)
	if err != nil {
		log.Println(err)
	}
}
