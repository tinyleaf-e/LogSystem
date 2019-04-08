package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var RedisConn redis.Conn

func initRedis(){
	RedisConn,_ =redis.Dial("tcp","localhost:6379")

}

func addAuthInfo(id string, passwd string, role string) string{
	rel,err:=RedisConn.Do("HMSET",id,"passwd",passwd,"role",role)
	if(err!=nil){
		fmt.Println("redis connection error")
		return err.Error()
	}
	fmt.Println(rel)
	return "success"
}


func checkAuthInfo(id string) int64{
	rel,err:=RedisConn.Do("HEXISTS",id,"role")
	if(err!=nil){
		fmt.Println("redis connection error")
		return -2
	}
	return rel.(int64)
}

func getRoleFromAuthInfo(token string) (role string, e error){
	role,e= redis.String(RedisConn.Do("HGET",token,"role"))
	return
}


func getUserIdFromAuthInfo(token string) (id string, e error){
	id,e= redis.String(RedisConn.Do("HGET",token,"userId"))
	return
}


func Auth(w *http.ResponseWriter,r *http.Request)(access bool) {
	resource := mux.Vars(r)["resource"]
	id := mux.Vars(r)["id"]
	token:=r.Header.Get("Authorization")[7:]
	if(checkAuthInfo(token)<=0){
		responseBuilder(w,http.StatusUnauthorized,"用户权限信息已过期")
		return false
	}
	role,err1:=getRoleFromAuthInfo(token)
	userId,err2:=getUserIdFromAuthInfo(token)
	if(err1!=nil||err2!=nil){
		responseBuilder(w,http.StatusUnauthorized,"用户权限信息已过期")
		return false
	}
	access=false
	switch resource {
	case "user":
		if(role=="admin"){
			access = true
		}else{
			if(id!=""&&id==userId){
				access = true
			}else {
				access = false
			}
		}
	case "project":
		if(role=="user"){
			if(r.FormValue("userId")!=""){
				if(userId==r.FormValue("userId")){
					access = true
				}else {
					access = false
				}
			}else{
				project,err:=getProject(id)
				if(err!=nil){
					if(err.Error()=="record not found"){
						responseBuilder(w,http.StatusNotFound,"该资源不存在")
					}else{
						responseBuilder(w,http.StatusBadGateway,err)
					}
				}else if(userId==project.UserId){
					access = true
				}else {
					access = false
				}
			}

		}else {
			access = false
		}
	case "format":
		if(role=="user"){
			if(r.FormValue("projectId")!=""){
				project,err:=getProject(r.FormValue("projectId"))
				if(err!=nil&&userId==project.UserId){
					access = true
				}else {
					access = false
				}
			}else{
				format,err:=getFormat(id)
				if(err!=nil){
					if(err.Error()=="record not found"){
						responseBuilder(w,http.StatusNotFound,"该资源不存在")
					}else{
						responseBuilder(w,http.StatusBadGateway,err)
					}
				}else
				{
					project,err:=getProject(format.ProjectId)
					if(err!=nil&&userId==project.UserId){
						access = true
					}else {
						access = false
					}

				}

			}

		}else {
			return false
		}
	case "format-item":
		if(role=="user"){
			format,err:=getFormat(r.FormValue("formatId"))
			if(err!=nil){
				access = false
			}
			project,err:=getProject(format.ProjectId)
			if(err!=nil&&userId==project.UserId){
				access = true
			}else {
				access = false
			}
		}else {
			access = false
		}
	case "log":
		if(role=="user"){
			format,err:=getFormat(r.FormValue("formatId"))
			if(err!=nil){
				access = false
			}
			project,err:=getProject(format.ProjectId)
			if(err!=nil&&userId==project.UserId){
				access = true
			}else {
				access = false
			}
		}else {
			access = false
		}
	}
	if(!access){
		responseBuilder(w,http.StatusUnauthorized,"权限不足，或未获取到权限信息")
	}
	return
}