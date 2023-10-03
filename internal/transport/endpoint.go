package transport

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) GetById(c *gin.Context) {

	uid := c.Param("id")

	data, err := app.orderSvc.Get(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "not found",
			"data":   "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"data":   data,
	})

}
