package model

import "github.com/jinzhu/gorm"
import _"github.com/jinzhu/gorm/dialects/mssql"

type User struct {
	Id string
	Name string
	Passwd string
	Email string
	Role string
	Remark string
}

type Project struct {
	Id string
	UserId string
	Ticker string
	CreateTime string
	Desc string
}

type Format struct{
	Id string
	Name string
	ProjectId string
	Desc string
}

type FormatItem struct {
	Id string
	Name string
	FormatId string
	Order int
	Type string
	Desc string
	Example string
}

type Log struct{
	Id string
	PattenId string
	Type string
	CreateTime string
	IpAddress string
	Agent string
	ShortString1 string
	ShortString2 string
	ShortString3 string
	ShortString4 string
	ShortString5 string
	ShortString6 string
	ShortString7 string
	ShortString8 string
	ShortString9 string
	ShortString10 string
	LongString1 string
	LongString2 string
	LongString3 string
	LongString4 string
	LongString5 string
	Int1 int
	Int2 int
	Int3 int
	Int4 int
	Int5 int
	Bool1 bool
	Bool2 bool
	Bool3 bool
}

var DB *gorm.DB

func connectDB(connStr string)(db *gorm.DB, err error){
	db, err=gorm.Open("mssql",connStr)
	return
}

func getUser(id string)(user User, e error){
	e = DB.Table("user").Where("id=?",id).First(&user).Error
	return
}


func getAllUsers()(user[] User, e error){
	e = DB.Table("user").Scan(&user).Error
	return
}
func deleteUser(id string)(rel string, e error){
	e=DB.Table("user").Where("id=?",id).Delete(User{}).Error
	return
}

func addUser(user User)(rel string, e error){
	e=DB.Table("user").Create(&user).Error
	return
}
func updateUser(user User)(rel string, e error){
	e=DB.Table("user").Where("id=?",user.Id).Delete(User{}).Error
	if(e!=nil){
		return
	}
	e = DB.Table("user").Create(&user).Error
	return
}