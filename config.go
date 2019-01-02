package Load_Balancer

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func setDefaults(p * proxy){
	if p.port == 0{
		p.port = 80
	}
	if p.scheme == ""{
		p.scheme = "http"
	}
	if p.hostname == ""{
		p.hostname = "localhost"
	}
}

func Validate(condition bool, errorMsg string) string {
	if condition{
		return errorMsg
	} else{
		return ""
	}
}

func removeEmpty(errList []string)[]string{
	filtered := []string{}
	for i,e := range errList{
		if e != ""{
			filtered = append(filtered,e)
		}
	}
	return filtered
}

func generateValidationErrors(p proxy) []string{
	return removeEmpty([]string{
		Validate(len(p.servers) == 0,"Should not be empty"),
	})
}

func validateFields(p proxy) error{
	errorList := generateValidationErrors(p)
	if len(errorList) != 0{
		return errors.New("Invalid configuration")
	}
	return nil
}

func ReadConfig(configFileName string) (proxy,error){
	proxyObj := proxy{}
	file,err := ioutil.ReadFile(configFileName)
	if err != nil{
		LogError("Error on reading file " + configFileName)
		return proxyObj,err
	}
	err = yaml.Unmarshal(file,&proxyObj)
	if err != nil{
		LogError("Error while unmarshalling!!")
		return proxyObj,err
	}
	setDefaults(&proxyObj)
	err = validateFields(proxyObj)
	if err != nil{
		return proxyObj,err
	}
	return proxyObj,nil
}
