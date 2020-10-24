package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-zoo/bone"
	"github.com/gorilla/sessions"
	"github.com/johnpili/notinaut/controllers"
	"github.com/johnpili/notinaut/models"
	"gopkg.in/yaml.v2"

	socketio "github.com/googollee/go-socket.io"
)

// Configurations / Settings
var (
	configuration  models.Config
	cookieStore    *sessions.CookieStore
	socketIOServer *socketio.Server
)

// This will handle the loading of config.yml
func loadConfiguration(c string) {
	f, err := os.Open(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	pid := os.Getpid()
	err := ioutil.WriteFile("application.pid", []byte(strconv.Itoa(pid)), 0666)
	if err != nil {
		log.Fatal(err)
	}

	var configLocation string
	flag.StringVar(&configLocation, "config", "config.yml", "Set the location of configuration file")
	flag.Parse()

	log.Println("------------------------------------------------------")
	log.Println("| Notinaut                                           |")
	log.Println("| Author: John Pili                                  |")
	log.Println("------------------------------------------------------")
	loadConfiguration(configLocation)

	envCookieKey := os.Getenv("ENV_HTTP_PROBE_COOKIE_KEY")
	if len(envCookieKey) > 0 {
		configuration.System.CookieKey = envCookieKey
	}

	if len(configuration.System.CookieKey) <= 0 {
		log.Fatalln("Missing cookie_key, please set the key value in the config.yml")
	}

	cookieKey := configuration.System.CookieKey
	cookieStore = sessions.NewCookieStore([]byte(cookieKey))

	viewBox := rice.MustFindBox("views")
	staticBox := rice.MustFindBox("static")

	controllersHub := controllers.New(viewBox, nil, cookieStore, &configuration)

	//#region SINGLE BINARY
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	//#endregion

	router := bone.New()
	router.Get("/static/", staticFileServer)
	controllersHub.BindRequestMapping(router)

	// CODE FROM https://medium.com/@mossila/running-go-behind-iis-ce1a610116df
	port := strconv.Itoa(configuration.HTTP.Port)
	if os.Getenv("ASPNETCORE_PORT") != "" { // get enviroment variable that set by ACNM
		port = os.Getenv("ASPNETCORE_PORT")
	}

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	if configuration.HTTP.IsTLS {
		log.Printf("Server running at https://localhost:%s/\n", port)
		log.Fatal(httpServer.ListenAndServeTLS(configuration.HTTP.ServerCert, configuration.HTTP.ServerKey))
		return
	}
	log.Printf("Server running at http://localhost:%s/\n", port)
	log.Fatal(httpServer.ListenAndServe())
}
