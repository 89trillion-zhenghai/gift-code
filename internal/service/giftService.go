package service

import (
	"encoding/json"
	"errors"
	"gift-code/internal/globalError"
	"gift-code/internal/model"
	"gift-code/internal/utils"
	"github.com/go-redis/redis"
	"time"
)

func CreateAndGetGiftCode(gift model.Gift) (code string,err error) {
	tm,err := time.ParseDuration(gift.Validity)
	//过期时间 = 当前时间 + 有效期
	expireDate := time.Now().Add(tm).Format("2006-01-02 15:04:05")
	gift.Validity = expireDate
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})
	defer client.Close()
	//创建前先判断礼品码是否重复,重复则重新生成随机礼品码
	gCode := ""
	for {
		gCode = utils.GetGiftCode()
		isHas := client.SAdd("GIFT_CODE",gCode)
		if isHas.Val() == 1 {
			break
		}
	}
	gift.GiftCode = gCode
	//如果是第三类的礼品码，则将可领取次数设置为-1
	if gift.GiftType == "3"{
		gift.AvailableTimes = "-1"
	}
	gift.AvailedTimes = "0"
	//将gift转成map
	m := gift.BeanToMap()
	//将gift以hash储存在redis里
	hmSet := client.HMSet("MAIN_"+gCode, m)
	if hmSet.Val() == "OK"{
		return gCode,nil
	}
	return gCode,errors.New("unknown error! please try again")
}

func GetGiftDetail(gc string) (resMap map[string]string,err error){
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})
	defer client.Close()
	//直接从hash里面去那取
	result := client.HGetAll("MAIN_" + gc).Val()
	//去REC_{gc}里面获取领取信息
	rec, err := client.HGetAll("REC_" + gc).Result()
	if err == nil {
		bts, err := json.Marshal(rec)
		if err != nil {

		}else{
			result["RecDetail"] = string(bts)
		}
	}else {
		result["RecDetail"] = ""
	}
	return result,nil
}

func RedeemGift(gc string,uuid string)(resMap map[string]string,testErr globalError.TestError){
	resMap = make(map[string]string)
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})
	defer client.Close()
	//判断礼品码合不合法(存在&未失效)
	flag, _ := client.SIsMember("GIFT_CODE", gc).Result()
	if !flag {
		return nil,globalError.TestError{
			Err:    errors.New("礼品码不存在！"),
			Status: "10001",
		}
	}
	vail, _ := client.HGet("MAIN_"+gc, "Validity").Result()
	now := time.Now().Format("2006-01-02 15:04:05")
	if now > vail {
		return  nil,globalError.TestError{
			Err:    errors.New("礼品码已失效！"),
			Status: "10002",
		}
	}
	//拿到该礼品码的类型
	gType := client.HGet("MAIN_"+gc, "GiftType").Val()
	//类型1
	if gType == "1" {
		//判断该用户是否领取过
		flag, _ = client.HExists("REC_"+gc, uuid).Result()
		if flag {
			return  nil,globalError.TestError{
				Err:    errors.New("您已经领取过本礼包！请检查"),
				Status: "10003",
			}
		}
		//判断可领取次数是否>1
		err := utils.IsAvailable(gc, client)
		if err != nil {
			return  nil,globalError.TestError{
				Err:    errors.New("礼包已经发送完毕"),
				Status: "10004",
			}
		}
		//将可领取次数-1,	已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailableTimes", -1)
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := utils.UpdateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
		//类型2
	}else if gType == "2"{
		//判断可领取次数是否>1
		err := utils.IsAvailable(gc, client)
		if err != nil {
			return  nil,globalError.TestError{
				Err:    errors.New("礼包已经发送完毕"),
				Status: "10004",
			}
		}
		//将可领取次数-1,	已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailableTimes", -1)
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := utils.UpdateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
		//类型3
	}else if gType == "3"{
		//已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := utils.UpdateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
	}

	return resMap,globalError.TestError{
		Status: "200",
	}
}




