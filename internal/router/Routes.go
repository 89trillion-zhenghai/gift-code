package router

import (
	"gift-code/internal/ctrl"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) *gin.Engine{
	r.POST("/createAndGetGiftCode",ctrl.CreateAndGetGiftCode)
	r.POST("/getGiftDetail",ctrl.GetGiftDetail)
	r.POST("/redeemGift",ctrl.RedeemGift)
	return r
}
