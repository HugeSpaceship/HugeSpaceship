package controllers

import (
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/internal/model/lbp_xml/npdata"
	"HugeSpaceship/internal/model/lbp_xml/recent_activity"
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
							ID:      "testdata-item",
							Subject: "This is a testdata news item",
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

											Content: "This is a pretty cool level\nYou Should play it\n&lt;slot type=&quot;user&quot; id=&quot;8&quot; icon=&quot;72d87f8f75b9b93f8cf6bdc8650c77cacb1cf22c&quot; big=&quot;true&quot;&gt;Shit Level&lt;/slot&gt;",
										},
									},
								},
							},
							Timestamp: time.Now().Unix(),
						},
					},
				},

				{
					ID:    "2",
					Title: "Other Category",
					Items: []lbp_xml.NewsItem{
						{
							ID:      "testdata-item-2",
							Subject: "This is a different testdata news item",
							Content: lbp_xml.NewsItemContent{
								Content: []lbp_xml.NewsFrame{
									{
										Width: "100",
										Title: "Other Test Frame",
										Item: lbp_xml.NewsFrameItem{
											Width: "50",
											NpHandle: npdata.NpHandle{
												Username: "Zaprit",
											},

											Content: "FUCK\n",
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
