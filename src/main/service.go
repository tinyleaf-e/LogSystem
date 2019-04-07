package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Pong"))
}

func PreprocessXHR(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type")
}


var Conf=NewConf("conf.ini")

func main(){
	Conf.parse()
	connStr,_ := Conf.get("connectionString")
	DB,_ = connectDB(connStr)
	r := mux.NewRouter()
	r.HandleFunc("/ping",Ping)
	//r.HandleFunc("/test",testbody).Methods("POST")
	r.HandleFunc("/create",createTables).Methods("GET")
	r.HandleFunc("/user/{id}",getUserById).Methods("GET")
	r.HandleFunc("/user",listAllUser).Methods("GET")
	r.HandleFunc("/user",addUserdbdb).Methods("POST")
	r.HandleFunc("/user/{id}",updateUserById).Methods("POST")
	r.HandleFunc("/user/{id}",deleteUserById).Methods("DELETE")

	r.HandleFunc("/project/{id}",getProjectById).Methods("GET")
	r.HandleFunc("/project",getProjectsByUserIddbdb).Methods("GET")
	r.HandleFunc("/project",addProjectdbdb).Methods("POST")
	r.HandleFunc("/project/{id}",updateProjectById).Methods("POST")
	r.HandleFunc("/project/{id}",deleteProjectById).Methods("DELETE")

	r.HandleFunc("/format/{id}",getFormatById).Methods("GET")
	r.HandleFunc("/format",getFormatsByProjectIddbdb).Methods("GET")
	r.HandleFunc("/format",addFormatdbdb).Methods("POST")
	r.HandleFunc("/format/{id}",updateFormatById).Methods("POST")
	r.HandleFunc("/format/{id}",deleteFormatById).Methods("DELETE")


	r.HandleFunc("/format-item",getFormatItemByFormatIddbdb).Methods("GET")
	r.HandleFunc("/format-item",addFormatItemdbdb).Methods("POST")
	r.HandleFunc("/format-item/{id}",updateFormatItemById).Methods("POST")
	r.HandleFunc("/format-item/{id}",deleteFormatItemById).Methods("DELETE")


	r.HandleFunc("/log",getLogByQuery).Methods("GET")
	r.HandleFunc("/log",addLogdbdb).Methods("POST")
	r.HandleFunc("/log/{id}",updateLogById).Methods("POST")

	http.Handle("/",r)
	http.ListenAndServe("localhost:8011",nil)
}