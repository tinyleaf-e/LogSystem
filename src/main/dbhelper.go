package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"time"
)

type User struct {
	Id string `gorm:"type:varchar(20)";primary_key`
	Name string `gorm:"type:varchar(20)"`
	Passwd string `gorm:"type:varchar(20)"`
	Email string `gorm:"type:varchar(30)"`
	Role string `gorm:"type:varchar(10)"`
	Remark string `gorm:"type:varchar(50)"`
}

type Project struct {
	Id string `gorm:"type:varchar(40)";primary_key`
	Name string `gorm:"type:varchar(30)"`
	UserId string `gorm:"type:varchar(20)"`
	Ticket string `gorm:"type:varchar(40)"`
	CreateTime time.Time `gorm:"type:datetime"`
	Desc string `gorm:"type:varchar(300)"`
}

type Format struct{
	Id string `gorm:"type:varchar(40)";primary_key`
	Name string `gorm:"type:varchar(30)"`
	ProjectId string `gorm:"type:varchar(40)"`
	CreateTime time.Time `gorm:"type:datetime"`
	Desc string `gorm:"type:varchar(300)"`
}

type FormatItem struct {
	Id string `gorm:"type:varchar(40)";primary_key`
	Name string `gorm:"type:varchar(30)"`
	FormatId string `gorm:"type:varchar(40)"`
	Order int `gorm:"type:int"`
	Type string `gorm:"type:varchar(15)"`
	Desc string `gorm:"type:varchar(300)"`
	Example string `gorm:"type:varchar(100)"`
	Editable bool `gorm:"type:tinyint"`
}

type Log struct{
	Id string `gorm:"type:varchar(40)";primary_key`
	FormatId string `gorm:"type:varchar(40)"`
	Type string `gorm:"type:varchar(10)"`
	CreateTime time.Time `gorm:"type:datetime"`
	IpAddress string `gorm:"type:varchar(30)"`
	Agent string `gorm:"type:varchar(200)"`
	ShortString1 string `gorm:"type:varchar(50)"`
	ShortString2 string `gorm:"type:varchar(50)"`
	ShortString3 string `gorm:"type:varchar(50)"`
	ShortString4 string `gorm:"type:varchar(50)"`
	ShortString5 string `gorm:"type:varchar(50)"`
	ShortString6 string `gorm:"type:varchar(50)"`
	ShortString7 string `gorm:"type:varchar(50)"`
	ShortString8 string `gorm:"type:varchar(50)"`
	ShortString9 string `gorm:"type:varchar(50)"`
	ShortString10 string `gorm:"type:varchar(50)"`
	LongString1 string `gorm:"type:varchar(1000)"`
	LongString2 string `gorm:"type:varchar(1000)"`
	LongString3 string `gorm:"type:varchar(1000)"`
	LongString4 string `gorm:"type:varchar(1000)"`
	LongString5 string `gorm:"type:varchar(1000)"`
	Int1 int `gorm:"type:int"`
	Int2 int `gorm:"type:int"`
	Int3 int `gorm:"type:int"`
	Int4 int `gorm:"type:int"`
	Int5 int `gorm:"type:int"`
	Bool1 bool `gorm:"type:tinyint"`
	Bool2 bool `gorm:"type:tinyint"`
	Bool3 bool `gorm:"type:tinyint"`
}

var DB *gorm.DB

func connectDB(connStr string)(db *gorm.DB, err error){
	db, err=gorm.Open("mssql",connStr)
	return
}

func createTable(){
	DB.Table("user").CreateTable(&User{})
	DB.Table("project").CreateTable(&Project{})
	DB.Table("format").CreateTable(&Format{})
	DB.Table("format_item").CreateTable(&FormatItem{})
	DB.Table("log").CreateTable(&Log{})
}


//User Dao
func getUser(id string)(user User, e error){
	e = DB.Table("user").Where("id=?",id).First(&user).Error
	return
}


func getAllUsers()(user[] User, e error){
	e = DB.Table("user").Scan(&user).Error
	return
}
func deleteUser(id[] string)(rel string, e error){
	e=DB.Table("user").Where("id in (?)",id).Delete(User{}).Error
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

//Project Dao
func getProject(id string)(project Project, e error){
	e = DB.Table("project").Where("id=?",id).First(&project).Error
	return
}


func getProjectsByUserId(userId string)(project[] Project, e error){
	e = DB.Table("project").Where("user_id=?",userId).Scan(&project).Error
	return
}
func deleteProject(id[] string)(rel string, e error){
	e=DB.Table("project").Where("id in (?)",id).Delete(Project{}).Error
	return
}

func deleteProjectByUserId(userId string)(rel string, e error){
	e=DB.Table("project").Where("user_id =?",userId).Delete(Project{}).Error
	return
}

func addProject(project Project)(rel string, e error){
	e=DB.Table("project").Create(&project).Error
	return
}
func updateProject(project Project)(rel string, e error){
	e=DB.Table("project").Where("id=?",project.Id).Delete(Project{}).Error
	if(e!=nil){
		return
	}
	e = DB.Table("project").Create(&project).Error
	return
}

//Format Dao
func getFormat(id string)(format Format, e error){
	e = DB.Table("format").Where("id=?",id).First(&format).Error
	return
}


func getFormatsByProjectId(projectId string)(format[] Format, e error){
	e = DB.Table("format").Where("project_id=?",projectId).Scan(&format).Error
	return
}
func deleteFormat(id[] string)(rel string, e error){
	e=DB.Table("format").Where("id in (?)",id).Delete(Format{}).Error
	return
}

func addFormat(format Format)(rel string, e error){
	e=DB.Table("Format").Create(&format).Error
	return
}
func updateFormat(format Format)(rel string, e error){
	e=DB.Table("format").Where("id=?",format.Id).Delete(Format{}).Error
	if(e!=nil){
		return
	}
	e = DB.Table("format").Create(&format).Error
	return
}

//FormatItem Dao
func getFormatItem(id string)(formatItem FormatItem, e error){
	e = DB.Table("format_item").Where("id=?",id).First(&formatItem).Error
	return
}


func getFormatItemsByFormatId(formatId string)(formatItem[] FormatItem, e error){
	e = DB.Table("format_item").Where("format_id=?",formatId).Order("order").Scan(&formatItem).Error
	return
}
func countFormatItemsByFormatId(formatId string)(count int, e error){
	e = DB.Table("format_item").Where("format_id=?",formatId).Count(&count).Error
	return
}
func deleteFormatItem(id[] string)(rel string, e error){
	e=DB.Table("format_item").Where("id in (?)",id).Delete(FormatItem{}).Error
	return
}

func addFormatItem(formatItem FormatItem)(rel string, e error){
	e=DB.Table("format_item").Create(&formatItem).Error
	return
}
func updateFormatItem(formatItem FormatItem)(rel string, e error){
	e=DB.Table("format_item").Where("id=?",formatItem.Id).Delete(FormatItem{}).Error
	if(e!=nil){
		return
	}
	e = DB.Table("format_item").Create(&formatItem).Error
	return
}

//Log Dao
func getLog(id string)(log Log, e error){
	e = DB.Table("log").Where("id=?",id).First(&log).Error
	return
}


func getLogsByQuery(query string, para...interface{})(log[] Log, e error){
	e = DB.Table("log").Where(query,para...).Scan(&log).Error
	return
}
func deleteLog(id[] string)(rel string, e error){
	e=DB.Table("log").Where("id in (?)",id).Delete(Log{}).Error
	return
}

func addLog(log Log)(rel string, e error){
	e=DB.Table("Log").Create(&log).Error
	return
}
func updateLog(log Log)(rel string, e error){
	e=DB.Table("log").Where("id=?",log.Id).Delete(Log{}).Error
	if(e!=nil){
		return
	}
	e = DB.Table("log").Create(&log).Error
	return
}