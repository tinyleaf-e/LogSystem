package main

import (
	"strings"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"time"
	"github.com/satori/go.uuid"
)

func main1(){
	sli:=[]string{"1","2"}
	str:=strings.Join(sli,"\",\"")
	fmt.Println(str)
}

func createTables(w http.ResponseWriter,r *http.Request){
	createTable()
}

func testbody(w http.ResponseWriter,r *http.Request) {
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



		PreprocessXHR(&w)
		id,_:=uuid.NewV4()
		projectId := r.FormValue("projectId")
		format:=Format{id.String(),name,projectId,time.Now(),desc}
		rel,err:=addFormat(format)
		if(err!=nil){
			dataMap["status"] = "error"
			dataMap["rel"] = err
		}else{
			hasError :=false
			for index, i := range item {
				//就在这里对di进行类型判断
				itemData, _ := i.(map[string]string)
				itemId,_:=uuid.NewV4()
				formatItem:=FormatItem{itemId.String(),itemData["name"],id.String(),index,itemData["type"],itemData["desc"],itemData["example"],itemData["editable"]=="1"}
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
