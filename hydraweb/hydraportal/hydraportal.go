package hydraportal

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"sync"

	hydraConfigurator "github.com/pienaahj/hydra/hydraconfigurator"
	"github.com/pienaahj/hydra/hydradblayer"
	"github.com/pienaahj/hydra/hydradblayer/passwordvault"
	hydratestapi "github.com/pienaahj/hydra/hydraweb/hydrarestapi"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
)

// create the template
var hydraWebTemplate *template.Template
var historylog = struct {
	logs          []string // chat history logs
	*sync.RWMutex          // more than one person will be using this
}{RWMutex: new(sync.RWMutex)}

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
		Filespath string   `json:"filespath"`
		Templates []string `json:"templates"`
	}{}
	err = hydraConfigurator.GetConfiguration(hydraConfigurator.JSON, &conf, "./hydraweb/portalconfig.json")
	if err != nil {
		return err
	}

	hydraWebTemplate, err = template.ParseFiles(conf.Templates...)
	if err != nil {
		log.Println(err)
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
	http.HandleFunc("/chat/", chathandler)
	http.Handle("/chatRoom/", websocket.Handler(chatWS))
	//  serve the TLS as it is a blocking call use a goroutine
	go func() {
		err = http.ListenAndServeTLS(":8062", "cert.pem", "key.pem", nil)
		log.Println(err)
	}()
	return http.ListenAndServe(":8061", nil)
}
func chatWS(ws *websocket.Conn) {
	chatserverIP := "127.0.0.1:2100"
	conn, err := net.Dial("tcp", chatserverIP)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	//  populate the message list in the browser
	historylog.RLock() // lock the logs to read from it
	for _, log := range historylog.logs {
		err := websocket.Message.Send(ws, log)
		if err != nil {
			historylog.RUnlock()
			return
		}
	}
	historylog.RUnlock() // unlock the logs

	//  don't block the other code execution
	go func() {
		//  use scanner to send chat messages
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			message := scanner.Text()
			err := websocket.Message.Send(ws, message)
			if err != nil {
				//  in production add recovery and reconnection logic here
				return
			}
		}
	}()

	//  handle messages received via the websocket
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			return
		}
		//  write the message to the logs
		_, err = conn.Write([]byte(message))
		if err != nil {
			historylog.Lock()
			//  if the logs are more than 20 trim them
			if len(historylog.logs) > 20 {
				historylog.logs = historylog.logs[1:]
			}
			historylog.logs = append(historylog.logs, message)
			historylog.Unlock()
		}
	}
}

func chathandler(w http.ResponseWriter, r *http.Request) {
	nameStruct := struct {
		Name string
	}{}
	r.ParseForm()
	if len(r.Form) == 0 {
		if cookie, err := r.Cookie("usernames"); err == nil {
			hydraWebTemplate.ExecuteTemplate(w, "login.html", nil)
			return
		} else {
			nameStruct.Name = cookie.Value
			hydraWebTemplate.ExecuteTemplate(w, "chat.html", nameStruct)
			return
		}
	}

	if r.Method == "POST" {
		var user, pass string
		if v, ok := r.Form["username"]; ok && len(v) > 0 {
			user = v[0]
		}

		if v, ok := r.Form["password"]; ok && len(v) > 0 {
			pass = v[0]
		}
		// user := r.Form["username"][0]
		// user := r.Form["password"][0]

		if !verifyPassword(user, pass) {
			hydraWebTemplate.ExecuteTemplate(w, "login.html", nil)
			return
		}

		nameStruct.Name = user
		if _, ok := r.Form["rememberme"]; ok {
			cookie := http.Cookie{Name: "username", Value: user}
			http.SetCookie(w, &cookie)
		}
	}
	hydraWebTemplate.ExecuteTemplate(w, "chat.html", nameStruct)
}

func verifyPassword(username, pass string) bool {
	db, err := passwordvault.ConnectPasswordVault()
	if err != nil {
		return false
	}
	defer db.Close()
	data, err := passwordvault.GetPasswordBytes(db, username)
	if err != nil {
		return false
	}
	hashedPass := md5.Sum([]byte(pass))
	//  [:] converts [i] to []
	return bytes.Equal(hashedPass[:], data)
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
