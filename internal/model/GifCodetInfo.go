package model

import "reflect"

type Gift struct {
	CreateUser 		string			`json:"createUser"`		//创建人员
	CreateTime 		string			`json:"createTime"`		//创建时间
	GiftCode 		string			`json:"giftCode"`		//礼品码
	Description 	string			`json:"description"`	//礼品描述
	GiftType 		string			`json:"giftType"`		//礼品码种类	1、指定用户一次性消耗 2、不指定用户限制兑换次数 3、不限用户不限次数兑换
	Validity 		string			`json:"validity"`		//有效期		单位：天
	AvailableTimes 	string			`json:"availableTimes"`	//可领取次数
	AvailedTimes   	string			`json:"availedTimes"`	//已领取次数
	GiftDetail	   	string			`json:"giftDetail"`		//礼品内容列表
}

//BeanToMap Gift转map
func (g Gift)BeanToMap() map[string]interface{} {
	m := make(map[string]interface{})
	typeOf := reflect.TypeOf(g)
	valueOf := reflect.ValueOf(g)
	for i := 0; i < typeOf.NumField(); i++ {
		key := typeOf.Field(i).Name
		value := valueOf.Field(i).Interface()
		m[key] = value
	}
	return m
}

func NewGift(crtU string,crtT string,gC string,desc string,gifT string,valid string,avlTimes string,availedTimes string,gifD string) Gift{
	return Gift{
		crtU,
		crtT,
		gC,
		desc,
		gifT,
		valid,
		avlTimes,
		availedTimes,
		gifD,
	}
}