package users

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
)

func UpdateUserHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		user := lbp_xml.User{}
		err := context.BindXML(&user)
		if err != nil {
			context.Error(err)
		}
	}
}
