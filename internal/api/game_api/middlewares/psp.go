package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Header names for psp version information

const PSPExeHeader = "X-exe-v"
const PSPDataHeader = "X-data-v"

func PSPVersionMiddleware(ctx *gin.Context) {
	// If we're not on PSP, then bail
	if ctx.GetHeader(PSPExeHeader) == "" {
		return
	}

	// Pass through PSP Data and Exe headers
	// TODO: Make it so you can enforce a version
	ctx.Header(PSPExeHeader, ctx.GetHeader(PSPExeHeader))
	ctx.Header(PSPDataHeader, ctx.GetHeader(PSPDataHeader))
}
