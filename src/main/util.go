package main

import (
	"net/http"
	"strings"
	"net"
	"regexp"
	"fmt"
	"encoding/json"
)

func ClientIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != ""  {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != ""  {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

//下划线转小驼峰
func camel2underline(str string)string{
	f := func(s string) string {
		return "_"+strings.ToLower(s)
	}
	re, _ := regexp.Compile("[A-Z]")
	return re.ReplaceAllStringFunc(str, f)

}

func responseBuilder(w *http.ResponseWriter,state int,rel interface{}){

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "false")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type")
	(*w).WriteHeader(state)

	dataMap:=make(map[string]interface{})
	if(state>=200&&state<300){
		dataMap["data"]=rel
	}else {
		dataMap["msg"]=fmt.Sprintf("%s", rel)
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(*w,string(j))
}
