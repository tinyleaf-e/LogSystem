package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"fmt"
	"model"
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

func Add(w http.ResponseWriter,r *http.Request){
	PreprocessXHR(&w)
	r.ParseMultipartForm(32 << 20)
	userid := r.MultipartForm.Value["userId"][0]
	passwd:=r.MultipartForm.Value["passwd"][0]
	role:=r.MultipartForm.Value["role"][0]
	rel := addAuthInfo(userid,passwd,role)
	dataMap:=make(map[string]interface{})
	dataMap["status"] = "ok"
	dataMap["rel"] = rel
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}

func getUserById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	user,err:=getUser(id)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = user
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}

func listAllUser(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	user,err:=getAllUsers()
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = user
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func deleteUserById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	user,err:=deleteUser(id)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = user
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func updateUserById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	user,err:=getUser(id)
	if err==nil{
		name:=r.FormValue("name")
		if name!=""{
			user.Name = name
		}
		passwd:=r.FormValue("passwd")
		if passwd!=""{
			user.Passwd = passwd
		}
		email:=r.FormValue("email")
		if email!=""{
			user.Email = email
		}
		role:=r.FormValue("role")
		if role!=""{
			user.Role = role
		}
		remark:=r.FormValue("remark")
		if remark!=""{
			user.Remark = remark
		}
		rel,err:=updateUser(user)
		if(err!=nil){
			dataMap["status"] = "error"
			dataMap["rel"] = fmt.Sprintf("%s", err)
		}else{
			dataMap["status"] = "ok"
			dataMap["rel"] = fmt.Sprintf("%s", rel)
		}

	}else{
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)

	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}
func addUser1(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	userid := r.FormValue("userId")
	name:=r.FormValue("name")
	passwd:=r.FormValue("passwd")
	email:=r.FormValue("email")
	role:=r.FormValue("role")
	remark:=r.FormValue("remark")
	user:=User{userid,name,passwd,email,role,remark}
	rel,err:=addUser(user)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = err
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = rel
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


var Conf=NewConf("conf.ini")

func main(){
	Conf.parse()
	connStr,_ := Conf.get("connectionString")
	DB,_ = connectDB(connStr)
	r := mux.NewRouter()
	r.HandleFunc("/ping",Ping)
	r.HandleFunc("/add",Add)
	r.HandleFunc("/user/{id}",getUserById).Methods("GET")
	r.HandleFunc("/user",listAllUser).Methods("GET")
	r.HandleFunc("/user",addUser1).Methods("POST")
	r.HandleFunc("/user/{id}",updateUserById).Methods("POST")
	r.HandleFunc("/user/{id}",deleteUserById).Methods("DELETE")
	http.Handle("/",r)
	http.ListenAndServe("localhost:8011",nil)
}