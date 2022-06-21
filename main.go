package main

import (
	"database/sql"
	"log"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "dicky123"
	dbname   = "learn_go"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Music struct {
	MusicID    string `json:"music_id"`
	MusicTitle string `json:"music_title"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Music `json:"data"`
	Message string  `json:"message"`
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func asciiJson(c *gin.Context) {
	data := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK, data)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// integrasi dgn postgres
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)

	CheckError(err)

	return db
}

// Create a music

// response and request handlers
func createMusic(c *gin.Context) {
	var music Music
	if err := c.ShouldBindJSON(&music); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	musicID := music.MusicID
	musicTitle := music.MusicTitle
	log.Println(musicID)
	log.Println(musicTitle)

	if musicID == "" || musicTitle == "" {
		// response = JsonResponse{Type: "error", Message: "You are missing musicID or musicTitle parameter."}
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are missing musicID or musicTitle parameter"})
	} else {
		printMessage("Inserting music into DB")

		fmt.Println("Inserting new music with ID: " + musicID + " and name: " + musicTitle)

		db := setupDB()

		var lastInsertID int
		err := db.QueryRow("INSERT INTO music.music(musicid, musictitle) VALUES($1, $2) returning id;", musicID, musicTitle).Scan(&lastInsertID)

		// check errors
		CheckError(err)
		c.JSON(http.StatusBadRequest, gin.H{"type": "success", "message": "The music has been inserted successfully!"})

		// response = JsonResponse{Type: "success", Message: "The music has been inserted successfully!"}
	}
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func main() {
	// init router gin
	router := gin.Default()

	// // connection string
	// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// // open database
	// db, err := sql.Open("postgres", psqlconn)
	// CheckError(err)

	// // close database
	// defer db.Close()

	// // check db
	// err = db.Ping()
	// CheckError(err)

	fmt.Println("Connected!")

	router.GET("/get_albums", getAlbums)
	router.POST("/post_albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.GET("/ascii", asciiJson)

	// api connect with postgres
	router.POST("/add_music", createMusic)
	// router.HandleFunc("/movies/", createMusic).Methods("POST")

	router.Run("localhost:8080")
}
