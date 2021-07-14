package main

import (
	"gift-code/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/createAndGetGiftCode",controller.CreateAndGetGiftCode)
	r.POST("/getGiftDetail",controller.GetGiftDetail)
	r.POST("/redeemGift",controller.RedeemGift)
	r.Run(":8080")
}




