package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApiRoute(e *echo.Echo, db *gorm.DB) {
	//aGroup := e.Group("api/v1")
	//telcoScoreController := config.InjectTelcoScoreController(db)
	//scoreproController := config.InjectScoreproController(db)
	//kmbController := config.InjectKmbScoreproController(db)
	//wgController := config.InjectWgScoreproController(db)
	//telcoGroup := aGroup.Group("/score")
	//{
	//	telcoGroup.POST("/credit/:phoneNumber", telcoScoreController.CreditScore)
	//	telcoGroup.POST("/credit/:phoneNumber/limit", telcoScoreController.CreditScoreLimit)
	//	telcoGroup.GET("/credit/detail/:id", telcoScoreController.Detail)
	//	telcoGroup.POST("/experian", telcoScoreController.Experian)
	//	telcoGroup.POST("/token", telcoScoreController.GetToken)
	//	telcoGroup.POST("/pickle", telcoScoreController.InternalScoring)
	//}
	//
	//{
	//	scoreproGroup := aGroup.Group("/scorepro")
	//	scoreproGroup.POST("/idx", scoreproController.Scoring)
	//	scoreproGroup.GET("/detail/:id", scoreproController.Detail)
	//}
	//
	//{
	//	kmbGroup := aGroup.Group("/scorepro/kmb")
	//	kmbGroup.POST("/idx", kmbController.Scoring)
	//}
	//
	//{
	//	wgGroup := aGroup.Group("/scorepro/wg")
	//	wgGroup.POST("/idx", wgController.Scoring)
	//}

}
