package controllers

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
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
											NpHandle: lbp_xml.NpHandle{
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
