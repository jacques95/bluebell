package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func SignUp(params *models.ParamSignUp) {

	mysql.InsertUser()
}
