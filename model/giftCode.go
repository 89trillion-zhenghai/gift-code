package model

import "reflect"

type Gift struct {
	CreateUser 		string				//创建人员
	CreateTime 		string				//创建时间
	GiftCode 		string				//礼品码
	Description 	string				//礼品描述
	GiftType 		string				//礼品码种类	1、指定用户一次性消耗 2、不指定用户限制兑换次数 3、不限用户不限次数兑换
	Validity 		string				//有效期		单位：天
	AvailableTimes 	string				//可领取次数
	AvailedTimes   	string				//已领取次数
	GiftDetail	   	string				//礼品内容列表
}

type RecDetail struct {
	GiftCode	string		//礼品码
	Detail		[]string	//领取信息列表
}

func NewGift(crtU string,crtT string,gC string,desc string,gifT string,valid string,avlTimes string,availedTimes string,gifD string) Gift{
	gift:= Gift{
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
	return gift
}

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


