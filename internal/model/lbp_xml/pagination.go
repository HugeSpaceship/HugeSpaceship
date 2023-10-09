package lbp_xml

import "github.com/gin-gonic/gin"

type PaginationData struct {
	Size  uint `form:"pageSize"`
	Start uint `form:"pageStart"`
	Page  uint `form:"page"`
}

func GetPageinationData(ctx *gin.Context) (PaginationData, error) {
	pageData := PaginationData{}
	err := ctx.ShouldBindQuery(&pageData)
	if pageData.Page != 0 {
		pageData.Size = 50
		pageData.Start = 1 * pageData.Page
	}
	return pageData, err
}

func (p PaginationData) GetData() (uint, uint) {
	return p.Size, p.Start
}
