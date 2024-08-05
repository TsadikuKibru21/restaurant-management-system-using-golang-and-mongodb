package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/invoices", controller.GetInvoices())
	incommingRoutes.GET("/invoices/:invoice_id", controller.GetInvoice())
	incommingRoutes.POST("/invoices", controller.CreateInvoice())
	incommingRoutes.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())

}
