package lbp_xml

import (
	"github.com/HugeSpaceship/HugeSpaceship/pkg/utils"
	"net/http"
	"strconv"
)

type PaginationData struct {
	Size   uint `form:"pageSize"`
	Start  uint `form:"pageStart"`
	Page   uint `form:"page"`
	Domain uint `form:"-"`
}

func GetPaginationData(r *http.Request) (PaginationData, error) {
	page := r.URL.Query().Get("page")
	var pageNumber int
	var err error

	var pageSize, pageStart int
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			return PaginationData{}, err
		}
		goto pastPageStart // If anyone want to rewrite this function to be less hellish that's just fine by me
		// For now though we're using the goto
	}

	pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		return PaginationData{}, err
	}
	pageStart, err = strconv.Atoi(r.URL.Query().Get("pageStart"))
	if err != nil {
		return PaginationData{}, err
	}

pastPageStart:

	pageData := PaginationData{
		Domain: utils.GetContextValue[uint](r.Context(), "domain"),
		Size:   uint(pageSize),
		Start:  uint(pageStart),
		Page:   uint(pageNumber),
	}

	if pageData.Page != 0 {
		pageData.Size = 50
		pageData.Start = 1 * pageData.Page
	}
	return pageData, err
}

func (p PaginationData) GetData() (uint, uint) {
	return p.Size, p.Start
}
