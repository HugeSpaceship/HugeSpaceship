package controllers

import (
	"HugeSpaceship/internal/model/lbp_xml"
	"bytes"
	"encoding/xml"
	"github.com/gin-gonic/gin"
)

// NotificationController is a stub due to unknown schema
func NotificationController() gin.HandlerFunc {
	return func(c *gin.Context) {
		notifications := []lbp_xml.Notification{{
			XMLName:  xml.Name{Local: "notification"},
			Type:     "moderationNotification",
			Text:     "Really cool Test Notification that should get shown multiple times",
			Extended: "",
		},
			{
				XMLName:  xml.Name{Local: "notification"},
				Type:     "moderationNotification",
				Text:     "I shitted all over the place",
				Extended: "",
			},
		}
		notificationBuf := new(bytes.Buffer)
		for _, notification := range notifications {
			b, err := xml.Marshal(&notification)
			if err != nil {
				c.Error(err)
				break
			}
			notificationBuf.Write(b)
		}

		c.Data(200, "text/xml", notificationBuf.Bytes())
	}
}
