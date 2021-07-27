package main

import (
	"gift-code/app/http"
	"gift-code/internal/model"
	"log"
)

func main() {
	//项目入口
	err := model.InitRc()
	if err != nil {
		log.Println(err.Error())
		return
	}
	http.InitRun()
}




