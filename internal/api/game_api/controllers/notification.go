package controllers

import (
	"net/http"
)

// NotificationController is a stub due to unknown schema
func NotificationController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: Write notification code, see #51
		//notificationBuf := new(bytes.Buffer)
		//for _, notification := range notifications {
		//	b, err := xml.Marshal(&notification)
		//	if err != nil {
		//		c.Error(err)
		//		break
		//	}
		//	notificationBuf.Write(b)
		//}
		//
		//c.Data(200, "text/xml", notificationBuf.Bytes())

	}
}
