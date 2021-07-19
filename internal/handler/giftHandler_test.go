package handler

import (
	"gift-code/internal/model"
	"gift-code/internal/service"
	"testing"
	"time"
)

func TestCreateAndGetGiftCode(t *testing.T) {
	gift := model.NewGift("admin", time.Now().Format("2006-01-02 15:04:05"), "", "十周年活动奖励", "1", "20m", "80", "", "{\"1001\":\"10\",\"1002\":\"5\"}")
	code, err := service.CreateAndGetGiftCode(gift)
	if err != nil {
		t.Error(err)
	}
	println(code)
}

func TestGetGiftDetail(t *testing.T) {
	giftCode := "88I6NC4N"
	detail, err := service.GetGiftDetail(giftCode)
	if err != nil {
		t.Error(err)
	}
	for k,v:=range detail{
		println(k,":",v)
	}
}

func TestRedeemGift(t *testing.T) {
	giftCode := "81414012"
	uuid := "smallBai"
	gift, err := service.RedeemGift(giftCode, uuid)
	if err.Err != nil{
		t.Error(err.Err)
	}
	for k,v:=range gift{
		println(k,":",v)
	}
}
