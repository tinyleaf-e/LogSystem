package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var RedisConn redis.Conn

func initRedis(){
	RedisConn,_ =redis.Dial("tcp","localhost:6379")

}

func addAuthInfo(id string, passwd string, role string) string{
	initRedis()
	defer RedisConn.Close()
	rel,err:=RedisConn.Do("HMSET",id,"passwd",passwd,"role",role)
	if(err!=nil){
		fmt.Println("redis connection error")
		return err.Error()
	}
	fmt.Println(rel)
	return "success"
}


func checkAuthInfo(id string) string{
	initRedis()
	defer RedisConn.Close()
	rel,err:=RedisConn.Do("HEXISTS",id,"passwd")
	if(err!=nil){
		fmt.Println("redis connection error")
		return err.Error()
	}
	fmt.Println(rel)
	return ""
}