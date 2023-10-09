package controllers

import (
	"HugeSpaceship/internal/model/lbp_xml/recent_activity"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"time"
)

func StreamHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timestamp := time.Now().UTC().UnixMilli()
		ctx.XML(200, recent_activity.Stream{
			StartTimestamp: timestamp - 50000,
			EndTimestamp:   timestamp,
			Groups: recent_activity.Groups{Groups: []recent_activity.Group{
				{Type: "news", NewsID: "1", Timestamp: timestamp, Events: recent_activity.GroupEvents{Events: []recent_activity.Event{
					{Type: "news_post", Timestamp: timestamp - 25000, NewsID: "1"},
				}}},
			}},
			News: recent_activity.News{Items: []recent_activity.NewsItem{
				{XMLName: xml.Name{
					Local: "item",
				}, ID: 1, Title: "Test Post", Text: "Test Text", Summary: "Test Summary", Date: timestamp - 25000, Image: recent_activity.Image{Alignment: "left"}, Category: "no_category", Background: ""},
			}},
		})
	}
}
