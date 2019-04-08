package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"strings"
	"github.com/satori/go.uuid"
	"time"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"strconv"
)

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
	if(!Auth(&w,r)){ return }

	id := mux.Vars(r)["id"]
	user,err:=getUser(id)
	if(err!=nil){
		responseBuilder(&w,http.StatusInternalServerError,err)
	}else{
		responseBuilder(&w,http.StatusOK,user)
	}
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
	ids:=strings.Split(id,",")//TODO 还要删该用户下的一系列东西
	user,err:=deleteUser(ids)
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
func addUserdbdb(w http.ResponseWriter,r *http.Request) {
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

//Project

func getProjectById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	project,err:=getProject(id)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = project
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}

func getProjectsByUserIddbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	userId:=r.FormValue("userId")
	projects,err:=getProjectsByUserId(userId)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = projects
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func deleteProjectById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	ids:=strings.Split(id,",")
	rel,err:=deleteProject(ids)//TODO delete related
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = rel
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func updateProjectById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	project,err:=getProject(id)
	if err==nil{
		name:=r.FormValue("name")
		if name!=""{
			project.Name = name
		}
		desc:=r.FormValue("desc")
		if desc!=""{
			project.Desc = desc
		}
		rel,err:=updateProject(project)
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

func refreshProjectTicketById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	ticket,_:=uuid.NewV4()
	project,err:=getProject(id)
	if err==nil{
		project.Ticket=ticket.String()
		rel,err:=updateProject(project)
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
}
func addProjectdbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id,_:=uuid.NewV4()
	userid := r.FormValue("userId")
	name:=r.FormValue("name")
	desc:=r.FormValue("desc")
	ticktet,_:=uuid.NewV4()
	project:=Project{id.String(),name,userid,ticktet.String(),time.Now(),desc}
	rel,err:=addProject(project)
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


//Format

func getFormatById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	format,err:=getFormat(id)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = format
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}

func getFormatsByProjectIddbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	projectId:=r.FormValue("projectId")
	formats,err:=getFormatsByProjectId(projectId)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = formats
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func deleteFormatById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	ids:=strings.Split(id,",")
	rel,err:=deleteFormat(ids)//TODO delete related
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = rel
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func updateFormatById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	format,err:=getFormat(id)
	if err==nil{
		name:=r.FormValue("name")
		if name!=""{
			format.Name = name
		}
		desc:=r.FormValue("desc")
		if desc!=""{
			format.Desc = desc
		}
		rel,err:=updateFormat(format)
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

func addFormatdbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dataMap["status"] = "error"
		dataMap["rel"] = "read body err"
		return
	}else{
		bodyJson, _ := simplejson.NewJson([]byte(string(body)))
		item, _ := bodyJson.Get("item").Array()
		name,_:=bodyJson.Get("name").String()
		desc,_:=bodyJson.Get("desc").String()
		projectId,_:=bodyJson.Get("projectId").String()

		id,_:=uuid.NewV4()
		format:=Format{id.String(),name,projectId,time.Now(),desc}
		rel,err:=addFormat(format)
		if(err!=nil){
			dataMap["status"] = "error"
			dataMap["rel"] = err
		}else{
			hasError :=false
			for index, i := range item {
				//就在这里对di进行类型判断
				itemData, _ := i.(map[string]interface{})
				itemId,_:=uuid.NewV4()
				formatItem:=FormatItem{itemId.String(),itemData["name"].(string),id.String(),index,itemData["type"].(string),itemData["desc"].(string),itemData["example"].(string),itemData["editable"].(bool)}
				_,err:=addFormatItem(formatItem)
				if(err!=nil){
					hasError=true
					dataMap["status"] = "error"
					dataMap["rel"] = err
					break
				}
			}
			if(hasError){

				dataMap["status"] = "error"
				dataMap["rel"] = "error when create item"
			}else {
				dataMap["status"] = "ok"
				dataMap["rel"] = rel
			}
		}
	}


	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


//FormatItem


func getFormatItemByFormatIddbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	formatId:=r.FormValue("formatId")
	formatItems,err:=getFormatItemsByFormatId(formatId)
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = formatItems
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func deleteFormatItemById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id := mux.Vars(r)["id"]
	ids:=strings.Split(id,",")
	rel,err:=deleteFormatItem(ids)//TODO delete related
	dataMap:=make(map[string]interface{})
	if(err!=nil){
		dataMap["status"] = "error"
		dataMap["rel"] = fmt.Sprintf("%s", err)
	}else{
		dataMap["status"] = "ok"
		dataMap["rel"] = rel
	}
	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}


func updateFormatItemById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	formatItem,err:=getFormatItem(id)
	if err==nil{
		name:=r.FormValue("name")
		if name!=""{
			formatItem.Name=name
		}
		itemType:=r.FormValue("type")
		if itemType!=""{
			formatItem.Type=itemType
		}
		desc:=r.FormValue("desc")
		if desc!=""{
			formatItem.Desc=desc
		}
		example:=r.FormValue("example")
		if example!=""{
			formatItem.Example=example
		}
		editable:=r.FormValue("editable")
		if editable!=""{
			formatItem.Editable=(editable=="1")
		}
		rel,err:=updateFormatItem(formatItem)
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

func addFormatItemdbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	id,_:=uuid.NewV4()
	name:=r.FormValue("name")
	formatId := r.FormValue("formatId")
	order,_:=countFormatItemsByFormatId(formatId)
	itemType:=r.FormValue("type")
	desc:=r.FormValue("desc")
	example:=r.FormValue("example")
	editable:=r.FormValue("editable")
	formatItem:=FormatItem{id.String(),name,formatId,order,itemType,desc,example,editable=="1"}
	rel,err:=addFormatItem(formatItem)
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



//Log


func getLogByQuery(w http.ResponseWriter,r *http.Request) {

	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dataMap["status"] = "error"
		dataMap["rel"] = "read body err"
		return
	}else{
		bodyJson, _ := simplejson.NewJson([]byte(string(body)))
		queryItem, err := bodyJson.Get("query").Array()
		formatId, err := bodyJson.Get("formatId").String()
		//name,_:=bodyJson.Get("name").String()

		queryFormat :=[]string{}
		var paraSlice []interface{}

		var sysnaxErrMsg string

		formatItems,err:=getFormatItemsByFormatId(formatId)

		if(err!=nil){
			dataMap["status"] = "error"
			dataMap["rel"] = err
		}else{
			hasError :=false
			for _, i := range queryItem {
				//就在这里对di进行类型判断

				itemData, _ := i.(map[string]interface{})

				for  _,fi:=range formatItems{
					if(itemData["key"].(string)==fi.Name){
						switch fi.Type {
						case "shortString1", "shortString2", "shortString3", "shortString4", "shortString5", "shortString6", "shortString7", "shortString8", "shortString9", "shortString10":
							switch itemData["condition"].(string) {
							case "=","<>":
								queryFormat = append(queryFormat,camel2underline(fi.Type)+" "+itemData["condition"].(string)+" ?" )
								paraSlice = append(paraSlice,itemData["value"].(string))
							case "like":
								queryFormat = append(queryFormat,camel2underline(fi.Type)+" "+itemData["condition"].(string)+" %?%" )
								paraSlice = append(paraSlice,itemData["value"].(string))
							case "in":
								queryFormat = append(queryFormat,camel2underline(fi.Type)+" "+itemData["condition"].(string)+" (?)" ) //TODO maybe wrong
								paraSlice = append(paraSlice,itemData["value"].([]interface{}))
							default:
								sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持条件 \""+itemData["condition"].(string)+"\"运算符"
							}
						case "int1", "int2", "int3", "int4", "int5":
							switch itemData["condition"].(string) {
							case "<",">","=","<=",">=","<>":
								queryFormat = append(queryFormat,fi.Type+" "+itemData["condition"].(string)+" ?" )
								paraSlice = append(paraSlice,itemData["value"].(int))
							case "in":
								queryFormat = append(queryFormat,fi.Type+" "+itemData["condition"].(string)+" (?)" ) //TODO maybe wrong
								paraSlice = append(paraSlice,itemData["value"].([]interface{}))
							default:
								sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持条件 \""+itemData["condition"].(string)+"\"运算符"
							}
						case "bool1", "bool2", "bool3":
							switch itemData["condition"].(string) {
							case "=","<>":
								queryFormat = append(queryFormat,fi.Type+" "+itemData["condition"].(string)+" ?" )
								paraSlice = append(paraSlice,itemData["value"].(bool))
							default:
								sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持条件 \""+itemData["condition"].(string)+"\"运算符"
							}
						default:
							sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持检索"
						}
					}
				}

				switch itemData["key"].(string) {
				case "type","ip_address","agent":
					switch itemData["condition"].(string) {
					case "=","<>":
						queryFormat = append(queryFormat,itemData["key"].(string)+" "+itemData["condition"].(string)+" ?" )
						paraSlice = append(paraSlice,itemData["value"].(string))
					case "like":
						queryFormat = append(queryFormat,itemData["key"].(string)+" "+itemData["condition"].(string)+" %?%" )
						paraSlice = append(paraSlice,itemData["value"].(string))
					case "in":
						queryFormat = append(queryFormat,itemData["key"].(string)+" "+itemData["condition"].(string)+" (?)" ) //TODO maybe wrong
						paraSlice = append(paraSlice,itemData["value"].([]interface{}))
					default:
						sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持条件 \""+itemData["condition"].(string)+"\"运算符"
					}
				case "create_time":
					switch itemData["condition"].(string) {
					case "<",">","<=",">=":
						queryFormat = append(queryFormat,itemData["key"].(string)+" "+itemData["condition"].(string)+" ?" )
						dt,_:=time.Parse("2006-01-02 15:04:05", itemData["value"].(string))
						paraSlice = append(paraSlice,dt)
					default:
						sysnaxErrMsg="类型\""+itemData["key"].(string)+"\" 暂不支持条件 \""+itemData["condition"].(string)+"\"运算符"
					}
				}

				if(err!=nil){
					hasError=true
					dataMap["status"] = "error"
					dataMap["rel"] = sysnaxErrMsg
					break
				}
			}
			logs,err:=getLogsByQuery(strings.Join(queryFormat," AND "), paraSlice...)


			//var mapResult []map[string]interface{}
			mapResult:= make([]map[string]interface{},0)
			for _,log:=range logs{
				logToShow:=make( map[string]interface{})
				logToShow["id"]=log.Id
				logToShow["format_id"]=log.FormatId
				logToShow["type"]=log.Type
				logToShow["create_time"]=log.CreateTime
				logToShow["ip_address"]=log.IpAddress
				logToShow["agent"]=log.Agent

				for _,item :=range formatItems {
					switch item.Type {
					case "shortString1":
						logToShow[item.Name]=log.ShortString1
					case "shortString2":
						logToShow[item.Name]=log.ShortString2
					case "shortString3":
						logToShow[item.Name]=log.ShortString3
					case "shortString4":
						logToShow[item.Name]=log.ShortString4
					case "shortString5":
						logToShow[item.Name]=log.ShortString5
					case "shortString6":
						logToShow[item.Name]=log.ShortString6
					case "shortString7":
						logToShow[item.Name]=log.ShortString7
					case "shortString8":
						logToShow[item.Name]=log.ShortString8
					case "shortString9":
						logToShow[item.Name]=log.ShortString9
					case "shortString10":
						logToShow[item.Name]=log.ShortString10
					case "longString1":
						logToShow[item.Name]=log.LongString1
					case "longString2":
						logToShow[item.Name]=log.LongString2
					case "longString3":
						logToShow[item.Name]=log.LongString3
					case "longString4":
						logToShow[item.Name]=log.LongString4
					case "longString5":
						logToShow[item.Name]=log.LongString5
					case "int1":
						logToShow[item.Name]=log.Int1
					case "int2":
						logToShow[item.Name]=log.Int2
					case "int3":
						logToShow[item.Name]=log.Int3
					case "int4":
						logToShow[item.Name]=log.Int4
					case "int5":
						logToShow[item.Name]=log.Int5
					case "bool1":
						logToShow[item.Name]=log.Bool1
					case "bool2":
						logToShow[item.Name]=log.Bool2
					case "bool3":
						logToShow[item.Name]=log.Bool3
					}
				}

				mapResult = append(mapResult,logToShow)
			}

			if(hasError){

				dataMap["status"] = "error"
				dataMap["rel"] = "error when create item"
			}else if(err!=nil){

				dataMap["status"] = "error"
				dataMap["rel"] = err
			}else {
				dataMap["status"] = "ok"
				dataMap["rel"] = mapResult

			}
		}
	}

	j,_:=json.Marshal(dataMap)
	fmt.Fprint(w,string(j))
}



func updateLogById(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	dataMap:=make(map[string]interface{})
	id := mux.Vars(r)["id"]
	log,err:=getLog(id)
	formatItems,_:=getFormatItemsByFormatId(log.FormatId)
	//TODO validate
	if err==nil{
		for _,item :=range formatItems {
			if(item.Editable!=true){
				continue
			}
			value := r.FormValue(item.Name)
			if (value == "") {
				continue
			}
			valueInt,_:=strconv.Atoi(value)
			switch item.Type {
			case "shortString1":
				log.ShortString1=value
			case "shortString2":
				log.ShortString2=value
			case "shortString3":
				log.ShortString3=value
			case "shortString4":
				log.ShortString4=value
			case "shortString5":
				log.ShortString5=value
			case "shortString6":
				log.ShortString6=value
			case "shortString7":
				log.ShortString7=value
			case "shortString8":
				log.ShortString8=value
			case "shortString9":
				log.ShortString9=value
			case "shortString10":
				log.ShortString10=value
			case "longString1":
				log.LongString1=value
			case "longString2":
				log.LongString2=value
			case "longString3":
				log.LongString3=value
			case "longString4":
				log.LongString4=value
			case "longString5":
				log.LongString5=value
			case "int1":
				log.Int1=valueInt
			case "int2":
				log.Int2=valueInt
			case "int3":
				log.Int3=valueInt
			case "int4":
				log.Int4=valueInt
			case "int5":
				log.Int5=valueInt
			case "bool1":
				log.Bool1=(value=="1")
			case "bool2":
				log.Bool2=(value=="1")
			case "bool3":
				log.Bool3=(value=="1")
			}
		}
		rel,err:=updateLog(log)
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

func addLogdbdb(w http.ResponseWriter,r *http.Request) {
	PreprocessXHR(&w)
	formatId := r.FormValue("formatId")

	log :=Log{}
	formatItems,_:=getFormatItemsByFormatId(formatId)
	for _,item :=range formatItems{
		value:=r.FormValue(item.Name)
		if(value==""){
			continue
		}
		valueInt,_:=strconv.Atoi(value)
		switch item.Type {
		case "shortString1":
			log.ShortString1=value
		case "shortString2":
			log.ShortString2=value
		case "shortString3":
			log.ShortString3=value
		case "shortString4":
			log.ShortString4=value
		case "shortString5":
			log.ShortString5=value
		case "shortString6":
			log.ShortString6=value
		case "shortString7":
			log.ShortString7=value
		case "shortString8":
			log.ShortString8=value
		case "shortString9":
			log.ShortString9=value
		case "shortString10":
			log.ShortString10=value
		case "longString1":
			log.LongString1=value
		case "longString2":
			log.LongString2=value
		case "longString3":
			log.LongString3=value
		case "longString4":
			log.LongString4=value
		case "longString5":
			log.LongString5=value
		case "int1":
			log.Int1=valueInt
		case "int2":
			log.Int2=valueInt
		case "int3":
			log.Int3=valueInt
		case "int4":
			log.Int4=valueInt
		case "int5":
			log.Int5=valueInt
		case "bool1":
			log.Bool1=(value=="1")
		case "bool2":
			log.Bool2=(value=="1")
		case "bool3":
			log.Bool3=(value=="1")
		}
	}

	logId,_:=uuid.NewV4()
	log.Id=logId.String()
	log.FormatId=formatId
	log.Type=r.FormValue("TYPE")
	log.CreateTime=time.Now()
	log.IpAddress=ClientIP(r)
	log.Agent=r.Header.Get("User-Agent")

	rel,err:=addLog(log)
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
