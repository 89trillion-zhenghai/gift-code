package service

import (
	"gift-code/internal/globalError"
	"gift-code/internal/model"
	"gift-code/internal/utils"
	"strconv"
	"time"
)

func CreateAndGetGiftCode(gift model.Gift) (code string,err error) {
	//创建前先判断礼品码是否重复,重复则重新生成随机礼品码
	gCode := ""
	for {
		gCode = utils.GetGiftCode()
		if model.GiftIsExit(gCode){
			break
		}
	}
	gift.GiftCode = gCode
	//将gift以hash储存在redis里
	err = model.SetGift(gift)
	err = model.SetAvailableDetail(gift.GiftCode)
	err = model.SetAvailableTime(gift.GiftCode)
	if err != nil{
		return "",globalError.DBError("redis发送错误",globalError.RedisNotConnect)
	}
	return gCode,nil
}

func GetGiftDetail(code string) (interface{},error){
	resMap := make(map[string]interface{})
	//从MAIN_{code}获取gift主体信息
	gift, _ := model.GetGift(code)
	if len(gift) == 0 {
		return nil,globalError.GiftCodeError("礼品码有误或不存在",globalError.GiftCodeNotExist)
	}
	//从AVAILABLE_{code}获取gift领取次数
	times, _ := model.GetAvailableTime(code)
	resMap["AvailableTime"] = times
	for k, v := range gift {
		resMap[k] = v
	}
	//从DETAIL_{code}获取gift领取详情
	detail, _ := model.GetAvailableDetail(code)
	if len(detail) == 0{
		resMap["AvailableDetail"] = ""
	}else{
		resMap["AvailableDetail"] = detail
	}
	return resMap,nil
}

func RedeemGift(code string,name string)(interface{},error){
	//判断礼品码是否存在
	gift, _ := model.GetGift(code)
	if len(gift) == 0 {
		return nil,globalError.GiftCodeError("礼品码有误或不存在",globalError.GiftCodeNotExist)
	}
	//判断礼品码是否失效
	now := time.Now().Unix()
	validity := gift["Validity"]
	val, _ := strconv.Atoi(validity)
	if int64(val) < now{
		return nil,globalError.GiftCodeError("礼品码失效",globalError.GiftCodeExpired)
	}
	//获取礼品码的种类
	giftType := gift["GiftType"]
	switch giftType {
	case "1":
	//第一类礼品码：指定用户一次性消耗
	//只需要判断领取列表是否有值，没有值就直接领取
		flag := model.GiftIsAvailed(code)
		if !flag {
			return nil,globalError.GiftCodeError("礼品码已失效",globalError.GiftCodeIsInvalid)
		}
		//领取次数+1，领取列表追加
		model.IncrAvailable(code)
		model.AppendUser(code,name)
		//给用户增加奖励
		detail := gift["GiftDetail"]
		model.SetUser(name)
		res := map[string]string{
			"1001":"10",
			"1002":"20",
			"1003":"30",
			"1004":"40",
			"1005":"50",
		}
		model.UserGetGift(name,res)
		return detail,nil
	case "2":
	//第二类礼品码：不指定用户限制兑换次数
	//先判断领取次数是否大于0，再判断是否领取过
		availableTime, _ := model.GetAvailableTime(code)
		time, _ := strconv.Atoi(gift["AvailableTimes"])
		times,_ := strconv.Atoi(availableTime)
		if time - times < 1 {
			return nil,globalError.GiftCodeError("礼品已被领取光了",globalError.GiftIsOver)
		}
		res, err := UserIsAvailed(code, name, gift)
		return res,err
	default:
	//第三类礼品码：不限用户不限次数兑换
	//判断用户是否领取过，没有则领取
		res, err := UserIsAvailed(code, name, gift)
		return res,err
	}

}

func UserIsAvailed(code string,name string,gift map[string]string) (interface{},error) {
	if !model.UserIsAvailed(code,name){
		return nil,globalError.GiftCodeError("你已经领取过本礼品了",globalError.GiftCodeReceived)
	}
	//领取次数+1，领取列表追加
	model.IncrAvailable(code)
	model.AppendUser(code,name)
	//给用户增加奖励
	detail := gift["GiftDetail"]
	model.SetUser(name)
	res := map[string]string{
		"1001":"10",
		"1002":"20",
		"1003":"30",
		"1004":"40",
		"1005":"50",
	}
	model.UserGetGift(name,res)
	return detail,nil
}




