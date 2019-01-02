package Load_Balancer

import "net/http"

type proxy struct {
	hostname string
	port int
	scheme string
	servers []server
}

func (p proxy) handler(w http.ResponseWriter, r * http.Request) {
	//handle
}

