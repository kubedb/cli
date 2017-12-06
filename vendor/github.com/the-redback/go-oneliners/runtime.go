package oneliners

import (
	"fmt"
	"log"
	"runtime"
	"encoding/json"
	"strings"
)

func FILE(a ...interface{}) {
	_, file, ln, ok := runtime.Caller(1)
	if ok {
		fmt.Println("__FILE__", file, "__LINE__", ln)
		if len(a) > 0 {
			fmt.Println(a...)
		}
	} else {
		log.Fatal("Failed to detect runtime caller info.")
	}
}


func PrettyJson(a interface{},msg ...string) {
	_, file, ln, ok := runtime.Caller(1)
	if ok {
		fmt.Println("__FILE__", file, "__LINE__", ln)
		if a!=nil {
			data,_:=json.MarshalIndent(a,"","   ")
			if len(msg) >0 {
				str:= strings.Trim(fmt.Sprintf("%v",msg), "[]")
				fmt.Println("=====================[",str,"]=====================")
			} else {
				fmt.Println("===============================================================")
			}
			fmt.Println(string(data))
		}
	} else {
		log.Fatal("Failed to detect runtime caller info.")
	}
}