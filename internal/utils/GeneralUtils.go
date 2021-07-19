package utils

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//StrToMap 将json型字符串转成map
func StrToMap(str string) map[string]interface{} {
	res := make(map[string]interface{})
	bts := []byte(str)
	err := json.Unmarshal(bts, &res)
	if err != nil {

	}
	return res
}

//GetGiftCode 随机获得8位礼品码 由大写字母和数字组成
func GetGiftCode() string {
	str := make([]string,6)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 8; i++ {
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

//IsAvailable 判断该礼品码是否可以被领取
func IsAvailable(gc string,client *redis.Client) (err error){
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

//UpdateREC 在领取列表追加领取记录
func UpdateREC(gc string,uuid string,time string,client *redis.Client)  string{
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
		client.HMSet("USER_"+uuid,StrToMap(giftDetail))
	}else{
		//如果用户记录存在，则将用户仓库取出来
		detail := StrToMap(giftDetail)
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