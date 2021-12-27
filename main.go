package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func dbInit() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot open the database")
	}
	db.AutoMigrate(&Board{}, &Comment{})
}

type Board struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string
	Comments []Comment `gorm:"foreignkey:BoardID"`
}

func getAllBoards() []Board {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot open the db in getAllBoards")
	}

	var boards []Board
	db.Order("created_at desc").Find(&boards)
	return boards
}

func createBoard(title string) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot create new board")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()

	db.Create(&Board{ID: id, Title: title, Comments: nil})
}

type Comment struct {
	gorm.Model
	ID      string `gorm:"primaryKey"`
	BoardID string
	Content string
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/**")

	dbInit()

	router.GET("/", func(ctx *gin.Context) {
		allBoards := getAllBoards()
		ctx.HTML(200, "index.html", gin.H{
			"boards": allBoards,
		})
	})

	router.GET("new/board/", func(ctx *gin.Context) {
		ctx.HTML(200, "create_board.html", gin.H{"title": "new board"})
	})

	router.POST("new/board/post/", func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		createBoard(title)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
