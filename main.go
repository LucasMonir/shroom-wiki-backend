package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	auth "shroom-wiki-backend/Auth"
	middleware "shroom-wiki-backend/Middleware"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB

type shroom struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Genus       string `json:"genus"`
	Species     string `json:"species"`
	Img         string `json:"img"`
	Edible      string `json:"edible"`
	Toxic       string `json:"toxic"`
}

func hasError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getUser1(writter http.ResponseWriter, request *http.Request) {
	rows, err := db.Query(`select * from users where id = 1`)

	hasError(err)

	defer rows.Close()

	users := make([]auth.User, 0)

	for rows.Next() {
		user := auth.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		hasError(err)
		users = append(users, user)
	}

	hasError(err)
	fmt.Println(users)
	middleware.JSON(writter, http.StatusOK, users)
}

func getShrooms(writter http.ResponseWriter, request *http.Request) {
	rows, err := db.Query(`select mushroom.id, name, description, img, genus.genus, species, edible, toxic  from mushroom join genus on genus.id = mushroom.genus`)
	hasError(err)

	defer rows.Close()

	shrooms := make([]shroom, 0)

	for rows.Next() {

		mushroom := shroom{}

		err = rows.Scan(&mushroom.ID, &mushroom.Name, &mushroom.Description, &mushroom.Img, &mushroom.Genus, &mushroom.Species, &mushroom.Edible, &mushroom.Toxic)
		hasError(err)

		shrooms = append(shrooms, mushroom)
	}

	hasError(err)
	fmt.Println(shrooms)
	middleware.JSON(writter, http.StatusOK, shrooms)
}

func getRandomShroom(writter http.ResponseWriter, request *http.Request) {
	rows, err := db.Query(`select mushroom.id, name, description, img, genus.genus, species, edible, toxic from mushroom join genus on genus.id = mushroom.genus order by random() limit 1;`)

	hasError(err)

	defer rows.Close()

	randomMushroom := shroom{}

	for rows.Next() {

		err = rows.Scan(&randomMushroom.ID, &randomMushroom.Name, &randomMushroom.Description, &randomMushroom.Img, &randomMushroom.Genus, &randomMushroom.Species, &randomMushroom.Edible, &randomMushroom.Toxic)
		hasError(err)
	}

	hasError(err)
	fmt.Println(randomMushroom)
	middleware.JSON(writter, http.StatusOK, randomMushroom)
}

func getShroomById(writter http.ResponseWriter, request *http.Request) {
	id, parseError := strconv.Atoi(request.URL.Query().Get("id"))

	hasError(parseError)

	var err error
	var rows *sql.Rows

	if id != 0 {
		rows, err = db.Query(`select mushroom.id, name, description, img, genus.genus, species, edible, toxic from mushroom join genus on genus.id = mushroom.genus where mushroom.id = $1`, id)
		defer rows.Close()

		shrooms := make([]shroom, 0)

		for rows.Next() {

			mushroom := shroom{}

			err = rows.Scan(&mushroom.ID, &mushroom.Name, &mushroom.Description, &mushroom.Img, &mushroom.Genus, &mushroom.Species, &mushroom.Edible, &mushroom.Toxic)
			hasError(err)

			shrooms = append(shrooms, mushroom)
		}

		hasError(err)
		fmt.Println(shrooms)
		middleware.JSON(writter, http.StatusNotFound, shrooms)
	}

	middleware.JSON(writter, http.StatusNotFound, "Nothing found")
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
	initDb()
	router := mux.NewRouter()

	router.HandleFunc("/shrooms", getShrooms).Methods(http.MethodGet)
	router.HandleFunc("/shroom", getShroomById).Methods(http.MethodGet)
	router.HandleFunc("/randomShroom", getRandomShroom).Methods(http.MethodGet)
	router.HandleFunc("/users", getUser1).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(":4200", handler)
}
