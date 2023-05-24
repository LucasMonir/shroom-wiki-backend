package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type shroom struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var shrooms = []shroom{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getShrooms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, shrooms)
}

func postShrooms(c *gin.Context) {
	var newShroom shroom

	if err := c.BindJSON(&newShroom); err != nil {
		return
	}

	shrooms = append(shrooms, newShroom)
	c.IndentedJSON(http.StatusCreated, newShroom)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getShrooms)
	router.POST("/albums", postShrooms)

	router.Run("localhost:4200")
}
