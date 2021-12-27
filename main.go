package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("cannot open the database")
    }
    db.AutoMigrate(&Board{}, &Comment{})
    defer db.Close()
}

type Board struct {
	gorm.Model
	ID        string
	Title     string
	Comments  []Comment `gorm:"foreignkey:BoardID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getAllBoards() []Board {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("cannot open the db in getAllBoards")
	}

	var boards []Board
	db.Order("create_at desc").Find(&boards)
	db.Close()
	return boards
}

func createBoard(title string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("cannot create new board")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()

	db.Create(&Board{ID: id, Title: title, Comments: nil})
	defer db.Close()
}

type Comment struct {
	gorm.Model
	ID        string
	BoardID   string
	content   string
	CreatedAt time.Time
	UpdateAt  time.Time
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