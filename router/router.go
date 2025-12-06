package router

import (
	"go-challenge/handlers"
	"go-challenge/repositories"
	"go-challenge/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"status": "OK"})
	})

	return r
}

// SetupRouter はDBを受け取りルーターをセットアップする
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// ヘルスチェック
	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"status": "OK"})
	})

	// 依存性の注入
	locationRepo := repositories.NewLocationRepository(db)
	locationService := services.NewLocationService(locationRepo)
	locationHandler := handlers.NewLocationHandler(locationService)

	// エンドポイント登録
	api := r.Group("/api")
	{
		api.GET("/locations", locationHandler.GetLocations)
	}

	return r
}
