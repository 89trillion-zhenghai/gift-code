package ctrl

import (
	"gift-code/internal/handler"
	"gift-code/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//CreateAndGetGiftCode 创建一个礼品对象，返回一个礼品码
func CreateAndGetGiftCode(c *gin.Context)  {
	createName := c.PostForm("userName")
	createTime := time.Now().Format("2006-01-02 15:04:05")
	description := c.PostForm("description")
	giftType := c.PostForm("giftType")
	validity := c.PostForm("validity")
	availableTimes := c.PostForm("availableTimes")
	giftDetail := c.PostForm("giftDetail")
	gift := model.NewGift(createName, createTime, "", description, giftType, validity, availableTimes, "", giftDetail)
	giftCode, err := handler.CreateAndGetGiftCode(gift)
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"message":err.Error(),
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"giftCode":giftCode,
		})
	}
}

//GetGiftDetail 查询礼品信息
func GetGiftDetail(c *gin.Context){
	giftCode := c.PostForm("giftCode")
	resMap, err := handler.GetGiftDetail(giftCode)
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"message":err.Error(),
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"message":resMap,
		})
	}
}


//RedeemGift 兑换礼品，返回礼品内容
func RedeemGift(c *gin.Context){
	giftCode := c.PostForm("giftCode")
	userName := c.PostForm("uuid")
	resMap,err := handler.RedeemGift(giftCode,userName)
	if err.Err != nil{
		c.JSON(http.StatusOK,gin.H{
			"message":nil,
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"message":resMap,
		})
	}
}
