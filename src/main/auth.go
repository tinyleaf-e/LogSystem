package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
)

var RedisConn redis.Conn

func initRedis(redisHost string){
	var ee error
	RedisConn,ee =redis.Dial("tcp",redisHost)
	fmt.Println(ee)

}

func addAuthInfo(id string, role string) (token string,err error){
	initRedis(RedisHost)
	defer RedisConn.Close()
	currentToken,err:=redis.String(RedisConn.Do("GET",id))
	if(err==nil){
		exists,_:=RedisConn.Do("HEXISTS",currentToken,"id")
		if(exists.(int64)>0){
			token=currentToken
			return
		}
	}

	uid,_:=uuid.NewV4()
	token=uid.String()
	_,err=redis.String(RedisConn.Do("HMSET",token,"id",id,"role",role,"EX","7200"))
	if(err==nil){
		RedisConn.Do("EXPIRE",token,"7200")
		RedisConn.Do("SET",id,token)
		return
	}
	return
}


func checkAuthInfo(id string) int64{
	initRedis(RedisHost)
	defer RedisConn.Close()
	rel,err:=RedisConn.Do("HEXISTS",id,"id")
	if(err!=nil){
		fmt.Println("redis connection error")
		return -2
	}
	return rel.(int64)
}

func getRoleFromAuthInfo(token string) (role string, e error){
	initRedis(RedisHost)
	defer RedisConn.Close()
	role,e= redis.String(RedisConn.Do("HGET",token,"role"))
	return
}


func getUserIdFromAuthInfo(token string) (id string, e error){
	initRedis(RedisHost)
	defer RedisConn.Close()
	id,e= redis.String(RedisConn.Do("HGET",token,"id"))
	return
}


func Auth(w *http.ResponseWriter,r *http.Request)(access bool) {
	resource := mux.Vars(r)["resource"]
	id := mux.Vars(r)["id"]
	if(r.Header.Get("Authorization")==""){
		responseBuilder(w,http.StatusUnauthorized,"未获取到用户权限信息")
		return
	}
	token:=r.Header.Get("Authorization")
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
				if(err==nil&&userId==project.UserId){
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
					if(err==nil&&userId==project.UserId){
						access = true
					}else {
						access = false
					}

				}

			}

		}else {
			return false
		}
	case "format-item"://TODO 往下还得详细讨论
		if(role=="user"){
			formatId:=""
			if(id!=""){
				formatItem,_:=getFormatItem(id)
				formatId=formatItem.FormatId
			}else {
				formatId=r.FormValue("formatId")
			}
			format,err:=getFormat(formatId)
			//format,err=getFormat(r.PostFormValue("formatId"))
			if(err!=nil){
				access = false
			}
			project,err:=getProject(format.ProjectId)
			if(err==nil&&userId==project.UserId){
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
			if(err==nil&&userId==project.UserId){
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

func AuthFormatId(w *http.ResponseWriter,r *http.Request, formatId string) ( access bool) {
	access=false;
	if (r.Header.Get("Authorization") == "") {
		responseBuilder(w, http.StatusUnauthorized, "未获取到用户权限信息")
		return
	}
	token := r.Header.Get("Authorization")
	if (checkAuthInfo(token) <= 0) {
		responseBuilder(w, http.StatusUnauthorized, "用户权限信息已过期")
		return
	}
	userId, err2 := getUserIdFromAuthInfo(token)
	if (err2 != nil) {
		responseBuilder(w, http.StatusUnauthorized, "用户权限信息已过期")
		return
	}
	format, err := getFormat(formatId)
	if (err != nil) {
		responseBuilder(w, http.StatusBadGateway, err)
		return
	}
	project, err := getProject(format.ProjectId)
	if (err == nil && userId == project.UserId) {
		access = true
		return
	}
	responseBuilder(w, http.StatusUnauthorized, "用户权限信息已过期")
	return
}