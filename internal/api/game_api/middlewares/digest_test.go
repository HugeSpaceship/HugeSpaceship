package middlewares

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestDigest = "1234"
const TestAltDigest = "5678"

const TestAuthCookie = "0123456789"

var cfg = config.Config{
	API: struct {
		EnforceDigest      bool `default:"false"`
		DigestKey          string
		AlternateDigestKey string
	}{EnforceDigest: true, DigestKey: TestDigest, AlternateDigestKey: TestAltDigest},
}

// The digests that we'll expect from the testdata
var expectedDigest = utils.CalculateDigest("/testdata", TestAuthCookie, TestDigest, []byte("Ok!"), false)
var expectedAltDigest = utils.CalculateDigest("/testdata", TestAuthCookie, TestAltDigest, []byte("Ok!"), false)

func setupDigestTestRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/testdata", DigestMiddleware(&cfg), func(ctx *gin.Context) {
		ctx.String(200, "Ok!")
	})
	return r
}
func TestDigestMiddleware(t *testing.T) {
	t.Run("testPrimaryDigest", testPrimaryDigest)
	t.Run("testAlternateDigest", testAlternateDigest)
}

func testPrimaryDigest(t *testing.T) {
	r := setupDigestTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", bytes.NewReader([]byte{}))
	req.AddCookie(&http.Cookie{
		Name:  "MM_AUTH",
		Value: TestAuthCookie,
	})
	digest := utils.CalculateDigest("/testdata", TestAuthCookie, TestDigest, nil, false)
	req.Header.Add(DigestHeaderA, digest)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Response code was not Ok")
	assert.Equal(t, "Ok!", w.Body.String())
	assert.Equal(t, w.Header().Get(DigestHeaderA), expectedDigest)
	assert.Equal(t, w.Header().Get(DigestHeaderB), digest)
}

func testAlternateDigest(t *testing.T) {
	r := setupDigestTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", bytes.NewReader([]byte{}))
	req.AddCookie(&http.Cookie{
		Name:  "MM_AUTH",
		Value: TestAuthCookie,
	})
	digest := utils.CalculateDigest("/testdata", TestAuthCookie, TestAltDigest, nil, false)
	req.Header.Add(DigestHeaderA, digest)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Response code was not Ok")
	assert.Equal(t, "Ok!", w.Body.String())
	assert.Equal(t, w.Header().Get(DigestHeaderA), expectedAltDigest)
	assert.Equal(t, w.Header().Get(DigestHeaderB), digest)
}
