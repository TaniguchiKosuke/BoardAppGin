package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	Comments []Comment `gorm:"foreignKey:BoardID"`
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

func createBoardComment(boardId string, content string) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot create comment")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()
	db.Create(&Comment{ID: id, BoardID: boardId, Content: content})
}

func getAllBoardComments(boardId string) []Comment {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot get comments")
	}

	var comments []Comment
	db.Table("comments").Select("COALESCE(BoardID,?)", boardId).Find(&comments)
	return comments
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

	router.GET("/", func(c *gin.Context) {
		allBoards := getAllBoards()
		c.HTML(200, "index.html", gin.H{
			"boards": allBoards,
		})
	})

	router.GET("/new/board", func(c *gin.Context) {
		c.HTML(200, "create_board.html", gin.H{"title": "new board"})
	})

	router.POST("/new/board/post", func(c *gin.Context) {
		title := c.PostForm("title")
		createBoard(title)
		c.Redirect(302, "/")
	})

	router.GET("/board/:id", func(c *gin.Context) {
		boardId := c.Param("id")
		allBoardComments := getAllBoardComments(boardId)
		c.HTML(200, "board.html", gin.H{
			"ID":       boardId,
			"comments": allBoardComments,
		})
	})

	router.POST("/board/:id/comment", func(c *gin.Context) {
		content := c.PostForm("content")
		boardId := c.Param("id")
		createBoardComment(boardId, content)
		c.Redirect(302, "/board/"+boardId)
	})

	router.Run()
}
