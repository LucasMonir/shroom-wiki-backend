package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type shroom struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Genus       string `json:"genus"`
	Species     string `json:"species"`
	Edible      string `json:"edible"`
	Toxic       string `json:"toxic"`
	Img         string `json:"img"`
}

var shrooms = []shroom{
	{ID: "1", Name: "Amanita phalloides", Description: "Amanita phalloides, commonly known as the death cap, is a deadly poisonous basidiomycete fungus, one of many in the genus Amanita.", Img: "https://encrypted-tbn2.gstatic.com/licensed-image?q=tbn:ANd9GcSjfmn3dnqwqbz-trVZPIF4nX9SIVKGJxa0wfVjQ1CJAx9-zUmHuEXmvMKdhfI4veqUrzFxce1MhFjEkLA"},
	{ID: "4", Name: "Amanita phanterina", Description: "Amanita pantherina, also known as the panther cap, false blusher, and the panther amanita due to its similarity to the true blusher, is a species of fungus found in Europe and Western Asia", Img: "https://encrypted-tbn1.gstatic.com/licensed-image?q=tbn:ANd9GcS2jPfF9E740FTgearAS0JLauc6sn_nPleiKl4Yg56krHxH5-K94dsnxx4xj8FEa8YdQxGhcKqUWsAFYz4"},
	{ID: "3", Name: "Amanita muscaria", Description: "Amanita muscaria, commonly known as the fly agaric or fly amanita, is a basidiomycete of the genus Amanita. It is also a muscimol mushroom.", Img: "https://www.naturezadivina.com.br/media/amasty/blog/amanita-muscaria.jpg"},
	{ID: "2", Name: "Amanita roseolamellata", Description: "This species is a readily recognisable Amanita with few collections. The habitat of this species is highly fragmented and in decline. Many of these areas are small patches of forest surrounded by urban or agricultural land, with recent (within the last few decades) severe decline in forest extent and increasing pressure for urban development.", Img: "http://www.amanitaceae.org/image/uploaded/r/roseolam139307_web.jpg"},
}

func hasError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getShrooms(c *gin.Context) {
	rows, err := db.Query(`select mushroom.id, name, description, img, genus.genus, species from mushroom join genus on genus.id = mushroom.genus`)
	fmt.Println(rows)
	hasError(err)

	defer rows.Close()

	shrooms := make([]shroom, 0)

	for rows.Next() {

		mushroom := shroom{}

		err = rows.Scan(&mushroom.ID, &mushroom.Name, &mushroom.Description, &mushroom.Img, &mushroom.Genus, &mushroom.Species)
		hasError(err)

		shrooms = append(shrooms, mushroom)
	}

	hasError(err)
	fmt.Println(shrooms)
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

func initDb() {
	var err error

	connStr := "postgresql://postgres:lu123cas@localhost:5432/shroomwiki?sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func main() {

	router := gin.Default()
	initDb()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	router.GET("/shrooms", getShrooms)
	router.POST("/shrooms", postShrooms)

	router.Run("localhost:4200")
}
