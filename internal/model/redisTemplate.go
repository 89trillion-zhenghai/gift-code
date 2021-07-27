package model

import (
	"gift-code/internal/globalError"
	"strconv"
	"time"
)

//GiftIsExit 判断礼品码是否已存在，存在返回false，不存在返回ture
func GiftIsExit(code string) bool {
	result, _ := Rc.Keys("MAIN_" + code).Result()
	return len(result) == 0
}

//SetGift 保存礼品码信息,将礼品码主体信息保存到hash
func SetGift(gift Gift) error {
	//礼品码主体信息保存到hash里
	gMap := gift.BeanToMap()
	_, err := Rc.HMSet("MAIN_"+gift.GiftCode, gMap).Result()
	return err
}

//GetGift 查询礼品码主体信息
func GetGift(code string) (map[string]string,error){
	resMap, err := Rc.HGetAll("MAIN_" + code).Result()
	if err != nil {
		return nil,err
	}
	return resMap,nil
}

//SetAvailableTime 领取次数保存到string
func SetAvailableTime(code string) error {
	_, err := Rc.Set("AVAILABLE_"+code,"0",0).Result()
	return err
}

//IncrAvailable 领取次数加1
func IncrAvailable(code string) error{
	_, err := Rc.Incr("AVAILABLE_" + code).Result()
	return err
}

//GetAvailableTime 查询领取次数
func GetAvailableTime(code string) (string, error){
	result, err := Rc.Get("AVAILABLE_" + code).Result()
	return result,err
}

//SetAvailableDetail 领取列表信息保存到hash
func SetAvailableDetail(code string) error {
	_, err := Rc.HMSet("DETAIL_"+code, nil).Result()
	return err
}

func GetAvailableDetail(code string) (map[string]string,error) {
	result, err := Rc.HGetAll("DETAIL_" + code).Result()
	return result,err
}

//GiftIsAvailed 判断礼品是否被领取过 领取过返回false
func GiftIsAvailed(code string) bool {
	result, _ := Rc.Keys("DETAIL_" + code).Result()
	return len(result) == 0
}

//UserIsAvailed 判断用户是否领取过 领取过返回false
func UserIsAvailed(code string,name string) bool {
	result, _ := Rc.HGet("DETAIL_"+code, name).Result()
	return len(result) == 0
}


//AppendUser 保存用户领取信息
func AppendUser(code string,name string) bool {
	now := time.Now().Format("2006-01-02 15:04:05")
	res, _ := Rc.HSet("DETAIL_"+code, name, now).Result()
	return res
}

//SetUser 用户信息用一个hash维护
func SetUser(name string) {
	isExit, _ := Rc.Keys("USER_" + name).Result()
	if len(isExit) == 0 {
		depot := map[string]interface{}{
			"1001":"0",
			"1002":"0",
			"1003":"0",
			"1004":"0",
			"1005":"0",
		}
		Rc.HMSet("USER_"+name, depot).Result()
	}
}

//GetUser 返回用户信息
func GetUser(name string) (map[string]string, error){
	result, err := Rc.HGetAll("USER_" + name).Result()
	return result,err
}

//UserGetGift 用户成功兑换奖励
func UserGetGift(name string,detail map[string]string) error{
	var err error
	for k, v := range detail {
		value, err := strconv.Atoi(v)
		if err != nil {
			return globalError.ServerError("服务器错误！")
		}
		_, err = Rc.HIncrBy("USER_"+name, k, int64(value)).Result()
	}
	return err
}
