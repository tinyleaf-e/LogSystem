package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func Ping(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Pong"))
}

func PreprocessXHR(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "false")
	(*w).Header().Add("Access-Control-Allow-Headers", "Authorization")
}



var Conf=NewConf("conf.ini")

func main(){
	Conf.parse()
	connStr,_ := Conf.get("connectionString")
	redisHost,_ := Conf.get("redisHost")
	port,_ := Conf.get("port")
	DB,_ = connectDB(connStr)
	initRedis(redisHost)
	r := mux.NewRouter()
	r.HandleFunc("/ping",Ping)
	r.HandleFunc("/token",getToken).Methods("GET")
	//r.HandleFunc("/test",testbody).Methods("POST")
	r.HandleFunc("/create",createTables).Methods("GET")
	r.HandleFunc("/{resource}/{id}",getUserById).Methods("GET")
	r.HandleFunc("/{resource:user}",Options).Methods("OPTIONS")
	r.HandleFunc("/{resource:user}",listAllUser).Methods("GET")
	r.HandleFunc("/{resource:user}",addUserdbdb).Methods("POST")
	r.HandleFunc("/{resource:user}/{id}",updateUserById).Methods("POST")
	r.HandleFunc("/{resource:user}/{id}",deleteUserById).Methods("DELETE")

	r.HandleFunc("/{resource:project}/{id}",getProjectById).Methods("GET")
	r.HandleFunc("/{resource:project}",Options).Methods("OPTIONS")
	r.HandleFunc("/{resource:project}",getProjectsByUserIddbdb).Methods("GET")
	r.HandleFunc("/{resource:project}",addProjectdbdb).Methods("POST")
	r.HandleFunc("/{resource:project}/{id}",updateProjectById).Methods("POST")
	r.HandleFunc("/{resource:project}/{id}",deleteProjectById).Methods("DELETE")

	r.HandleFunc("/{resource:format}/{id}",getFormatById).Methods("GET")
	r.HandleFunc("/{resource:format}",getFormatsByProjectIddbdb).Methods("GET")
	r.HandleFunc("/{resource:format}",addFormatdbdb).Methods("POST")
	r.HandleFunc("/{resource:format}/{id}",updateFormatById).Methods("POST")
	r.HandleFunc("/{resource:format}/{id}",deleteFormatById).Methods("DELETE")


	r.HandleFunc("/{resource:format-item}",getFormatItemByFormatIddbdb).Methods("GET")
	r.HandleFunc("/{resource:format-item}",addFormatItemdbdb).Methods("POST")
	r.HandleFunc("/{resource:format-item}/{id}",updateFormatItemById).Methods("POST")
	r.HandleFunc("/{resource:format-item}/{id}",deleteFormatItemById).Methods("DELETE")


	r.HandleFunc("/{resource:log}",getLogByQuery).Methods("GET")
	r.HandleFunc("/{resource:log}",addLogdbdb).Methods("POST")
	r.HandleFunc("/{resource:log}/{id}",updateLogById).Methods("POST")

	http.Handle("/",r)
	fmt.Println("system started, listening on :"+port)
	http.ListenAndServe(":"+port,nil)
}