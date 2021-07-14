package service

import (
	"encoding/json"
	"errors"
	"gift-code/model"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"strings"
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
		gCode = getGiftCode()
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

func RedeemGift(gc string,uuid string)(resMap map[string]string,err error){
	resMap = make(map[string]string)
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})
	defer client.Close()
	//判断礼品码合不合法(存在&未失效)
	flag, err := client.SIsMember("GIFT_CODE", gc).Result()
	if !flag {
		return nil,errors.New("礼品码不存在！")
	}
	vail, err := client.HGet("MAIN_"+gc, "Validity").Result()
	now := time.Now().Format("2006-01-02 15:04:05")
	if now > vail {
		return  nil,errors.New("礼品码已失效！")
	}
	//拿到该礼品码的类型
	gType := client.HGet("MAIN_"+gc, "GiftType").Val()
	//类型1
	if gType == "1" {
		//判断该用户是否领取过
		flag, err = client.HExists("REC_"+gc, uuid).Result()
		if flag {
			return nil,errors.New("您已经领取过本礼包！请检查")
		}
		//判断可领取次数是否>1
		err := isAvailable(gc, client)
		if err != nil {
			return nil,err
		}
		//将可领取次数-1,	已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailableTimes", -1)
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := updateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
		//类型2
	}else if gType == "2"{
		//判断可领取次数是否>1
		err := isAvailable(gc, client)
		if err != nil {
			return nil,err
		}
		//将可领取次数-1,	已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailableTimes", -1)
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := updateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
	//类型3
	}else if gType == "3"{
		//已领取次数+1
		client.HIncrBy("MAIN_"+gc, "AvailedTimes", 1)
		//在领取列表追加领取记录
		detail := updateREC(gc, uuid, now, client)
		resMap["GiftDetail"] = detail
	}

	return resMap,nil
}

//判断该礼品码是否可以被领取
func isAvailable(gc string,client *redis.Client) (err error){
	//判断可领取次数是否>1
	num := client.HGet("MAIN_"+gc, "AvailableTimes").Val()
	res, err := strconv.Atoi(num)
	if err != nil {
		return errors.New("数据异常，请联系管理员检查！")
	}
	if res < 1{
		return errors.New("来晚一步！礼品已经发完了")
	}
	return nil
}
//在领取列表追加领取记录
func updateREC(gc string,uuid string,time string,client *redis.Client)  string{
	flag := client.HExists("REC_"+gc, uuid).Val()
	if flag {
		times := client.HGet("REC_"+gc, uuid).Val()
		// times :     2006-01-02 15:04:05,2006-01-03 15:04:10
		ts := times+","+time
		client.HSet("REC_"+gc,uuid,ts)
	}else{
		//在领取列表追加领取记录
		client.HSet("REC_"+gc,uuid,time)
	}
	giftDetail := client.HGet("MAIN_"+gc, "GiftDetail").Val()
	bol := client.Exists("USER_" + uuid).Val()
	if bol == 0 {
		client.HMSet("USER_"+uuid,strToMap(giftDetail))
	}else{
		//如果用户记录存在，则将用户仓库取出来
		detail := strToMap(giftDetail)
		old := client.HGetAll("USER_" + uuid).Val()
		for k, v := range old {
			if _, ok := detail[k]; ok {
				// 存在,将detail[k]的值加上old[k]的值
				oldN, _ := strconv.Atoi(v)
				newN, _ := strconv.Atoi(detail[k].(string))
				detail[k] = oldN+newN
			}
			if _, ok := old[k]; !ok {
				// 不存在
				detail[k] = v
			}
		}
		client.HMSet("USER_"+uuid,detail)
	}
	//获取礼品内容
	return giftDetail
}

//将json型字符串转成map
func strToMap(str string) map[string]interface{} {
	res := make(map[string]interface{})
	bts := []byte(str)
	err := json.Unmarshal(bts, &res)
	if err != nil {

	}
	return res
}

//随机获得8位礼品码 由大写字母和数字组成
//ascii 大写字母 A～Z [65,90] 0~9 [48,57]
//随机0，1，1就在[65,90]随机 0就在[30,39]随机
func getGiftCode() string {
	str := make([]string,6)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		bt := rand.Intn(10)
		if bt > 5 {
			bt = rand.Intn(26)+65
		}else{
			bt = rand.Intn(9)+48
		}
		str = append(str, string(rune(bt)))
	}
	return strings.Join(str,"")
}