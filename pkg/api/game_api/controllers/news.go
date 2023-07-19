package controllers

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"HugeSpaceship/pkg/common/model/lbp_xml/npdata"
	"HugeSpaceship/pkg/common/model/lbp_xml/recent_activity"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"time"
)

func NewsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.XML(200, &lbp_xml.News{
			Categories: []lbp_xml.SubCategory{
				{
					ID:    "1",
					Title: "Test Category",
					Items: []lbp_xml.NewsItem{
						{
							ID:      "test-item",
							Subject: "This is a test news item",
							Content: lbp_xml.NewsItemContent{
								Content: []lbp_xml.NewsFrame{
									{
										Width: "100",
										Title: "Test Frame",
										Item: lbp_xml.NewsFrameItem{
											Width: "50",
											NpHandle: npdata.NpHandle{
												Username: "Zaprit282",
											},
											Content: "Test Frame Item",
										},
									},
								},
							},
							Timestamp: time.Now().Unix(),
						},
					},
				},
			},
		})
	}
}

func LBP2NewsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.XML(200, recent_activity.NewsItem{
			XMLName:    xml.Name{Local: "news"},
			ID:         1,
			Title:      "Test Title",
			Text:       "Test Text",
			Summary:    "Summary",
			Date:       time.Now().UTC().UnixMilli() - 25000,
			Image:      recent_activity.Image{Alignment: "left"},
			Category:   "no_category",
			Background: "",
		})
	}
}
