package Load_Balancer

import "net/http"

const configFileName  = "config.yml"

func main(){

	p,err := ReadConfig(configFileName)
	if err != nil{
		LogError("Error during configuration")
		return
	}
	http.HandleFunc("/",p.handler)
	http.ListenAndServe(":"+string(p.port),nil)
}
