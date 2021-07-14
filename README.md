## 项目简介

​	1、本项目主要实现了创建和验证礼品码，最终效果为管理员通过接口创建礼品码，用户领取礼品码，并且管理员可以查询礼品码使用情况。

## 快速上手

​	1、本项目由Go语言开发，数据库采用redis。需要配置Go开发环境，以及安装redis数据库

​	2、进入项目根目录，通过命令行命令启动项目

​			1、go build main.go		---编译

​			2、./main 						  ---启动项目

## 功能介绍

​	1、本项目核心数据为礼品信息，为其构建的数据结构为

```go
type gift struct {
   createUser        string          		//创建人员
   createTime        string          		//创建时间
   giftCode       	string          		//礼品码
   description    	string          		//礼品描述
   giftType      		string               //礼品码种类
   validity      		string          		//有效期
   availableTimes    string               //可领取次数
   availedTimes      string               //已领取次数
   giftDetail        string    	         //礼品内容列表
}
```

​	2、实现三个接口

```go
//CreateAndGetGiftCode 创建一个礼品对象，返回一个礼品码
func CreateAndGetGiftCode(c *gin.Context)  {...}

//GetGiftDetail 查询礼品信息
func GetGiftDetail(c *gin.Context){...}

//RedeemGift 兑换礼品
func RedeemGift(c *gin.Context){...}
```

## 数据库设计

​	1、本项目采用redis保存相关数据，redis有五种基本数据类型和几种特殊数据类型，根据项目需求分析，存储的数据需要完成以下功能。

​		1、设置有效时间，如果到期则将其设置为已过期并且将其持久化

​		2、用户输入符合要求的礼品码返回奖品内容

​		3、当用户成功兑换礼品时需要将已领取次数+1

​		4、管理员查看礼品码对应礼品领取情况

2、用到以下数据类型保存数据

​		1、set维护礼品码，保证礼品码不重复 	eg: GIFT_CODE 123    gift_code 234    

​		2、hash 维护领取列表，记录领取信息(领取人：领取时间) 	eg: REC_{giftCode} uuid1:times uuid2:times

​		3、hash保存礼品主题信息

## 流程图
未命名文件 (4).jpg![未命名文件 (4)](https://user-images.githubusercontent.com/86946999/125587537-409f3b8d-a4e4-4b7c-80f5-6ef39edb27d3.jpg)



