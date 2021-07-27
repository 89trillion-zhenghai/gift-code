package router

import (
	"gift-code/internal/ctrl"
	"gift-code/internal/globalError"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) *gin.Engine{
	r.POST("/createAndGetGiftCode",globalError.ErrorHandler(ctrl.CreateAndGetGiftCode))
	r.POST("/getGiftDetail",globalError.ErrorHandler(ctrl.GetGiftDetail))
	r.POST("/redeemGift",globalError.ErrorHandler(ctrl.RedeemGift))
	return r
}
