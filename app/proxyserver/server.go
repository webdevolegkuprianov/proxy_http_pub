package proxyserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/webdevolegkuprianov/proxy_http/app/model"

	logger "github.com/webdevolegkuprianov/proxy_http/app/logger"
)

//errors
var (
	errIp = errors.New("error request")
)

//server configure
type server struct {
	router *mux.Router
	config *model.Config
}

func newServer(config *model.Config) *server {
	s := &server{
		router: mux.NewRouter(),
		config: config,
	}
	s.configureRouter()
	return s
}

//write http error
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

//write http response
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

//configure proxy
func (s *server) configureRouter() {

	whiteIp := s.router.PathPrefix("/").Subrouter()
	whiteIp.Use(s.middleWare)

	//general
	//open
	whiteIp.HandleFunc("/authentication", s.handleReverseProxy("/authentication")).Methods("POST")
	//private
	//booking, forms submit
	whiteIp.HandleFunc("/auth/requestbooking", s.handleReverseProxy("/auth/requestbooking")).Methods("POST")
	whiteIp.HandleFunc("/auth/requestform", s.handleReverseProxy("/auth/requestform")).Methods("POST")
	//gaz crm
	whiteIp.HandleFunc("/auth/requestleadget", s.handleReverseProxy("/auth/requestleadget")).Methods("POST")
	whiteIp.HandleFunc("/auth/requestworklist", s.handleReverseProxy("/auth/requestworklist")).Methods("POST")
	whiteIp.HandleFunc("/auth/requeststatus", s.handleReverseProxy("/auth/requeststatus")).Methods("POST")
	//stock
	whiteIp.HandleFunc("/auth/getdatastocks", s.handleReverseProxy("/auth/getdatastocks")).Methods("GET")
	//prices
	whiteIp.HandleFunc("/auth/getbasicmodelsprice", s.handleReverseProxy("/auth/getbasicmodelsprice")).Methods("GET")
	whiteIp.HandleFunc("/auth/getoptionsprice", s.handleReverseProxy("/auth/getoptionsprice")).Methods("GET")
	whiteIp.HandleFunc("/auth/getgeneralprice", s.handleReverseProxy("/auth/getgeneralprice")).Methods("GET")
	//sprav models
	whiteIp.HandleFunc("/auth/getsprav", s.handleReverseProxy("/auth/getsprav")).Methods("GET")
	//options
	whiteIp.HandleFunc("/auth/getoptionsdata", s.handleReverseProxy("/auth/getoptionsdata")).Methods("GET")
	whiteIp.HandleFunc("/auth/getoptionsdatasprav", s.handleReverseProxy("/auth/getoptionsdatasprav")).Methods("GET")
	whiteIp.HandleFunc("/auth/getpacketsdata", s.handleReverseProxy("/auth/getpacketsdata")).Methods("GET")
	//colors
	whiteIp.HandleFunc("/auth/getcolorsdata", s.handleReverseProxy("/auth/getcolorsdata")).Methods("GET")

	//autoretail
	//open
	whiteIp.HandleFunc("/service/authentication", s.handleReverseProxyAutoretail("/service/authentication")).Methods("POST")
	//private
	whiteIp.HandleFunc("/auth/service/carservicedata", s.handleReverseProxyAutoretail("/auth/service/carservicedata")).Methods("POST")
}

//middleware white ip
func (s *server) middleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		whiteIpList := s.config.Spec.WhiteIp

		remoteIp := s.getAddr(r)

		for _, k := range whiteIpList {
			if k == remoteIp {
				next.ServeHTTP(w, r)
			}
		}
		s.error(w, r, http.StatusUnauthorized, errIp)

	})

}

//proxy method to general server
func (s *server) handleReverseProxy(urlExternal string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		urlInternal := fmt.Sprintf("%s%s", s.config.Spec.ProxyAddr.AddrServerRest, urlExternal)

		origin, err := url.Parse(urlInternal)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}

		proxy.ServeHTTP(w, r)

	}
}

//proxy method to autoretail server
func (s *server) handleReverseProxyAutoretail(urlExternal string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		urlInternal := fmt.Sprintf("%s%s", s.config.Spec.ProxyAddr.AddrServerRestAutoretail, urlExternal)

		origin, err := url.Parse(urlInternal)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}

		proxy.ServeHTTP(w, r)

	}
}

//get ip of client
func (s *server) getAddr(req *http.Request) string {

	remoteIP := ""
	if parts := strings.Split(req.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}

	if xff := strings.Trim(req.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}

	} else if xri := req.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP

}
