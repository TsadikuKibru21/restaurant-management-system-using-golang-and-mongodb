package routes

import (
	controller "restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(incommingRoutes *gin.Engine){

	incommingRoutes.GET("/menus",controller.GetMenus())
	incommingRoutes.GET("/menus/:menu_id",controller.GetMenu())
	incommingRoutes.GET("/menus",controller.CreateMenu())
	incommingRoutes.GET("/menus/:menu_id",controller.UpdateMenu())
}

