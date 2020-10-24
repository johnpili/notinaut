package controllers

import (
	"log"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/sessions"
	"github.com/johnpili/notinaut/models"
	"github.com/patrickmn/go-cache"
	"github.com/psi-incontrol/go-sprocket/page"
	"github.com/psi-incontrol/go-sprocket/sprocket"
	"github.com/rs/xid"
	"github.com/tarm/serial"
)

var (
	cookieStore   *sessions.CookieStore
	viewBox       *rice.Box
	staticBox     *rice.Box
	configuration *models.Config
	serialPort    *serial.Port
	ipCache       *cache.Cache
)

//New ...
func New(vb *rice.Box, sb *rice.Box, store *sessions.CookieStore, config *models.Config) *Hub {
	viewBox = vb
	staticBox = sb
	cookieStore = store
	configuration = config

	expiryDuration := time.Duration(config.IPCacheControl.ExpirySec) * time.Second
	purgeDuration := time.Duration(config.IPCacheControl.PurgeSec) * time.Second

	ipCache = cache.New(expiryDuration, purgeDuration)

	setupSerial()
	hub := new(Hub)
	hub.LoadControllers()
	return hub
}

//Hub ...
type Hub struct {
	Controllers []interface{}
}

func renderPage(w http.ResponseWriter, r *http.Request, vm interface{}, filenames ...string) {
	page := vm.(*page.Page)

	if page.Data == nil {
		page.SetData(make(map[string]interface{}))
	}

	if page.ErrorMessages == nil {
		page.ResetErrors("")
	}

	x, err := sprocket.GetTemplates(viewBox, filenames)
	err = x.Execute(w, page)
	if err != nil {
		log.Panic(err.Error())
	}
}
func generateGUID() string {
	guid := xid.New()
	guid.Time().Add(100)
	return guid.String()[12:20]
}

func setupSerial() {
	c := &serial.Config{Name: configuration.System.SerialName, Baud: configuration.System.SerialBaud}
	var err error
	serialPort, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
}
