package resources

import (
	"HugeSpaceship/pkg/common/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetResourceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		resource, tx, size, err := db.GetResource(dbCtx, ctx.Param("hash"))
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}
		defer func(tx pgx.Tx, dbCtx context.Context, ctx *gin.Context) {
			err := tx.Commit(dbCtx)
			if err != nil {
				ctx.Error(err)
			}
		}(tx, dbCtx, ctx)
		ctx.DataFromReader(200, size, "application/octet-stream", resource, nil)
		err = resource.Close()
		if err != nil {
			ctx.Error(err)
		}
	}
}
