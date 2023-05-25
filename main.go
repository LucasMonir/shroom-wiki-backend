package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type shroom struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genus       string `json:"genus"`
	Species     string `json:"species"`
	Img         string `json:"img"`
}

var shrooms = []shroom{
	{ID: "1", Title: "Amanita phalloides", Description: "Amanita phalloides, commonly known as the death cap, is a deadly poisonous basidiomycete fungus, one of many in the genus Amanita.", Img: "https://encrypted-tbn2.gstatic.com/licensed-image?q=tbn:ANd9GcSjfmn3dnqwqbz-trVZPIF4nX9SIVKGJxa0wfVjQ1CJAx9-zUmHuEXmvMKdhfI4veqUrzFxce1MhFjEkLA"},
	{ID: "4", Title: "Amanita phanterina", Description: "Amanita pantherina, also known as the panther cap, false blusher, and the panther amanita due to its similarity to the true blusher, is a species of fungus found in Europe and Western Asia", Img: "https://encrypted-tbn1.gstatic.com/licensed-image?q=tbn:ANd9GcS2jPfF9E740FTgearAS0JLauc6sn_nPleiKl4Yg56krHxH5-K94dsnxx4xj8FEa8YdQxGhcKqUWsAFYz4"},
	{ID: "3", Title: "Amanita muscaria", Description: "Amanita muscaria, commonly known as the fly agaric or fly amanita, is a basidiomycete of the genus Amanita. It is also a muscimol mushroom.", Img: "https://www.naturezadivina.com.br/media/amasty/blog/amanita-muscaria.jpg"},
	{ID: "2", Title: "Amanita roseolamellata", Description: "This species is a readily recognisable Amanita with few collections. The habitat of this species is highly fragmented and in decline. Many of these areas are small patches of forest surrounded by urban or agricultural land, with recent (within the last few decades) severe decline in forest extent and increasing pressure for urban development.", Img: "http://www.amanitaceae.org/image/uploaded/r/roseolam139307_web.jpg"},
}

func getShrooms(c *gin.Context) {
	c.JSON(http.StatusOK, shrooms)
}

func postShrooms(c *gin.Context) {
	var newShroom shroom

	if err := c.BindJSON(&newShroom); err != nil {
		return
	}

	shrooms = append(shrooms, newShroom)
	c.JSON(http.StatusCreated, newShroom)
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	router.GET("/shrooms", getShrooms)
	router.POST("/shrooms", postShrooms)

	router.Run("localhost:4200")
}
