package main

import (
	"fmt"
	"gift-code/model"
	"gift-code/service"
	"testing"
	"time"
)

func TestCreateAndGetGiftCode(t *testing.T) {
	createName := "admin"
	createTime := time.Now().Format("2006-01-02 15:04:05")
	description := "十周年活动"
	giftType := "2"
	validity := "10m"
	availableTimes := "10"
	giftDetail := "{\"道具\":\"20\",\"小兵\":\"20\",\"金币\":\"1000\",\"钻石\":\"10\"}"
	gift := model.NewGift(createName, createTime, "", description, giftType, validity, availableTimes, "", giftDetail)
	giftCode,err := service.CreateAndGetGiftCode(gift)
	if err != nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println(giftCode)
	}
}

func TestGetGiftDetail(t *testing.T) {
	giftCode := "8N440Q"
	resMap, err := service.GetGiftDetail(giftCode)
	if err != nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println(resMap)
	}
}

func TestRedeemGift(t *testing.T) {
	giftCode := "8N440Q"
	userName := "user01"
	resMap,err := service.RedeemGift(giftCode,userName)
	if err != nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println(resMap)
	}
}
