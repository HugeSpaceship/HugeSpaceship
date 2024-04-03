package slots

import (
	"fmt"
	"io"
	"net/http"
)

func UploadScoreHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
}
