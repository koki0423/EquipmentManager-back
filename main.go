package main

import (
	"equipmentManager/internal/database/db"
	"equipmentManager/internal/handler"
	"equipmentManager/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	_ "equipmentManager/docs"
)

// @title           備品管理システム
// @version         1.0
// @description     研究室内の備品を管理するシステムです。
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"*", //開発検証用
		},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ハンドラーインスタンスを生成
	sh := handler.NewSystemHandler()
	auh := handler.NewAuthHandler(db)
	ah:=handler.NewAssetHandler(db)
	lh:= handler.NewLendHandler(db)
	dh:= handler.NewDisposalHandler(db)

	router.InitRouter(r, sh, auh,ah,lh,dh)

	r.Run("0.0.0.0:8080")
}
