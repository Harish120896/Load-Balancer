package Load_Balancer

import (
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type proxy struct {
	hostname string
	port int
	scheme string
	servers []server
}

func (p proxy) origin() string {
	return p.scheme + "://" + p.hostname + ":" + strconv.Itoa(p.port)
}

func (p proxy) chooseServer(ignoreList map[string]bool) *server {
	min := 100
	index := 0
	for i,ser := range p.servers{
		if _,ok := ignoreList[ser.serverName]; ok{
			continue
		}
		if min > ser.connections {
			min = ser.connections
			index = i
		}
	}
	return &p.servers[index]
}

func (p proxy) reverseProxy(w http.ResponseWriter, r * http.Request,server * server) error {
	u,err := url.Parse(server.URL()+r.RequestURI)
	if err != nil{
		//handle
	}
	r.URL = u
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("Origin", p.origin())
	r.Host = server.URL()
	r.RequestURI = ""

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// TODO: If the server doesn't respond, try a new web server
	// We could return a status code from this function and let the handler try passing the request to a new server.
	resp, err := client.Do(r)
	if err != nil {
		// For now, this is a fatal error
		// When we can fail to another webserver, this should only be a warning.
		LogError("connection refused")
		return err
	}
	LogInfo("Recieved response: " + strconv.Itoa(resp.StatusCode))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError("Proxy: Failed to read response body")
		http.NotFound(w, r)
		return err
	}

	buffer := bytes.NewBuffer(bodyBytes)
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ";"))
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, buffer)
	return nil
}

func (p proxy) attemptServers(w http.ResponseWriter, r * http.Request, ignoreList map[string]bool){
	if float64(len(ignoreList)) >= math.Min(float64(3),float64(len(p.servers))){
		LogError("No servers found to get served!!")
		http.NotFound(w,r)
		return
	}
	server := p.chooseServer(ignoreList)
	server.connections += 1
	err := p.reverseProxy(w,r,server)
	server.connections -= 1
	if err != nil && strings.Contains(err.Error(),"connection refused"){
		ignoreList[server.hostname] = true
		p.attemptServers(w,r,ignoreList)
		return
	}

}

func (p proxy) handler(w http.ResponseWriter, r * http.Request) {
	p.attemptServers(w,r,make(map[string]bool))
}

