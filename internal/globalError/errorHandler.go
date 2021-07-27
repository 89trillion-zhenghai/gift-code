package globalError

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ServerExpression = 1001		//服务器异常
	RedisNotConnect	 = 1002		//redis连接异常
	GiftCodeExpired  = 1003		//礼品码已过期
	GiftCodeReceived = 1004		//该用户已经领取过礼品码了
	GiftCodeNotExist = 1005		//礼品码不存在/错误
	GiftCodeIsInvalid= 1006		//礼品码已失效
	GiftIsOver		 = 1007		//礼品被领取完毕
	ParamIsEmpty	 = 1008		//参数为空
	ParamIsIllegal	 = 1009		//参数不合法
)

type GlobalHandler func(c *gin.Context) (interface{}, error)

func ErrorHandler(handler GlobalHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			globalError := err.(GlobalError)
			c.JSON(globalError.Status, globalError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": result,
			"message":"success",
		})
	}
}

//ServerError 服务器异常
func ServerError(message string) GlobalError {
	return GlobalError{
		Status: http.StatusInternalServerError,
		Code: ServerExpression,
		Message: message,
	}
}

//DBError 数据库异常
func DBError(message string,code int) GlobalError {
	return GlobalError{
		Status: http.StatusInternalServerError,
		Code: code,
		Message: message,
	}
}

//GiftCodeError 礼品码异常
func GiftCodeError(message string, code int) GlobalError {
	return GlobalError{
		Status: http.StatusBadRequest,
		Code: code,
		Message: message,
	}
}

//ParamError 请求参数异常
func ParamError(message string,code int) GlobalError {
	return GlobalError{
		Status: http.StatusBadRequest,
		Code: code,
		Message: message,
	}
}