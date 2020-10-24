package controllers

import (
	"log"
	"net/http"

	"github.com/johnpili/notinaut/models"
	"github.com/patrickmn/go-cache"
	"github.com/psi-incontrol/go-sprocket/sprocket"

	"github.com/go-zoo/bone"

	"github.com/psi-incontrol/go-sprocket/page"
)

// PageController ...
type PageController struct{}

// RequestMapping ...
func (z *PageController) RequestMapping(router *bone.Mux) {
	router.GetFunc("/", z.IndexHandler)
	router.PostFunc("/", z.TriggerHandler)
}

// IndexHandler ...
func (z *PageController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	page := page.New()
	page.Title = "Notinaut"
	renderPage(w, r, page, "base.html", "index.html")
}

// TriggerHandler ...
func (z *PageController) TriggerHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := z.getIPDetails(r)
	_, found := ipCache.Get(ipInfo.IP)
	if !found {
		log.Printf("Adding into IP cache control %s\n", ipInfo.IP)
		ipCache.Set(ipInfo.IP, "", cache.DefaultExpiration)
		_, err := serialPort.Write([]byte("run\n"))
		if err != nil {
			log.Println(err)
		}
	}
	sprocket.RespondOkayJSON(w, "")
}

func (z *PageController) getIPDetails(r *http.Request) models.IPInfo {
	ip := ""
	if len(configuration.Extraction.HeaderKey) > 0 {
		ip = r.Header.Get(configuration.Extraction.HeaderKey) // Extract IP from header because we are using reverse proxy example X-Real-Ip
		if ip == "" {
			ip = r.RemoteAddr
		}
	}

	if len(ip) == 0 { // Fallback
		ip = extractIPAddress(r.RemoteAddr)
	}

	ipInfo := models.IPInfo{
		IP:        ip,
		UserAgent: r.Header.Get("User-Agent"),
	}

	if configuration.Extraction.DebugHeader {
		log.Print(r.Header)
	}

	return ipInfo
}

func extractIPAddress(ip string) string {
	if len(ip) > 0 {
		for i := len(ip); i >= 0; i-- {
			offset := len(ip)
			if (i + 1) <= len(ip) {
				offset = i + 1
			}
			if ip[i:offset] == ":" {
				return ip[:i]
			}
		}
	}
	return ip
}
