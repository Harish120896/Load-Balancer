package Load_Balancer

type server struct {
	serverName string
	hostname string
	port int
	scheme string
	connections int
}

func (s server) URL() string{
	return s.scheme + "://" + s.hostname + ":" + string(s.port)
}
