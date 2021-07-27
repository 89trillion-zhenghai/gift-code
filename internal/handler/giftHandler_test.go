package handler

import (
	"fmt"
	"gift-code/internal/model"
	"gift-code/internal/service"
	"testing"
	"time"
)

func TestCreateAndGetGiftCode(t *testing.T) {
	model.InitRc()
	gift := model.Gift{
		CreateUser:  "admin",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Description: "十周年活动奖励",
		GiftType:    "1",
		GiftDetail:  "{\"1001\":\"10\",\"1002\":\"5\"}",
	}
	validity := "10m"
	tm,err := time.ParseDuration(validity)
	//过期时间 = 当前时间 + 有效期
	expireDate := time.Now().Add(tm).Unix()
	gift.Validity = expireDate
	code, err := service.CreateAndGetGiftCode(gift)
	if err != nil {
		t.Error(err)
	}
	println(code)
}

func TestGetGiftDetail(t *testing.T) {
	model.InitRc()
	giftCode := "2X23871F"
	detail, err := service.GetGiftDetail(giftCode)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(detail)
}

func TestRedeemGift(t *testing.T) {
	model.InitRc()
	giftCode := "2X23871F"
	name := "smallBai"
	gift, err := service.RedeemGift(giftCode, name)
	if err != nil{
		t.Error(err)
	}
	fmt.Println(gift)
}
