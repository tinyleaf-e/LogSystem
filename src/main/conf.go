package main

import (
	"strings"
	"os"
	"bufio"
	"io"
	"errors"
)

type ConfSet struct{
	fname string
	parsed bool
	item map[string]string
}

func NewConf(fname string) *ConfSet {
	return &ConfSet{fname, false, make(map[string]string)}
}

func (c ConfSet) parse() error{
	c.parsed = true

	fp, err := os.Open(c.fname)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(fp)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		if len(line) == 0 {
			continue
		}

		l := strings.TrimSpace(string(line))

		if strings.HasPrefix(l,"#"){
			continue
		}


		parts := strings.SplitN(l, "=", 2)
		name, value := parts[0], parts[1]
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		c.item[name]=value
	}
	return nil
}

func (c ConfSet) get(key string) (string,error){
	value,ok:=c.item[key]
	if(ok){
		return value,nil
	}else{
		return "",errors.New("no such key")
	}
}