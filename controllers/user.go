package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	// 参数提取并校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断 err 是否为 validator 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSONP(http.StatusOK, gin.H{"msg": err.Error()})
			return
		}
		c.JSONP(http.StatusOK, gin.H{"msg": removeTopStruct(errs.Translate(trans))})
		return
	}

	// 业务处理
	logic.SignUp(p)

	// 返回响应
	c.JSONP(http.StatusOK, gin.H{"msg": "success"})
}
