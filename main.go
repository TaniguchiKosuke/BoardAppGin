package main

import (
	"fmt"

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

func getAllBoards(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot open the db in getAllBoards")
	}
	
	var boards []Board
	searchParam, queryExist := c.GetQuery("title")
	if queryExist {
		db.Where("title LIKE ?", "%"+searchParam+"%").Order("created_at desc").Find(&boards)		
	} else {
		db.Order("created_at desc").Find(&boards)
	}

	c.HTML(200, "index.html", gin.H{
		"boards": boards,
	})
}

func createBoard(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot create new board")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()
	title := c.PostForm("title")
	db.Create(&Board{ID: id, Title: title, Comments: nil})

	c.Redirect(302, "/")
}

func createBoardComment(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot create comment")
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()
	content := c.PostForm("content")
	boardId := c.Param("id")
	db.Create(&Comment{ID: id, BoardID: boardId, Content: content})

	c.Redirect(302, "/board/"+boardId)
}

func getAllBoardComments(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot get comments")
	}
	
	var board Board
	boardId := c.Param("id")
	boardObj := db.Where("id ?", boardId).First(&board)
	fmt.Println(board)

	var comments []Comment
	searchParam, queryExist := c.GetQuery("comment")
	if queryExist {
		db.Where("content LIKE ?", "%"+searchParam+"%").Where(&Comment{BoardID: boardId}).Order("created_at desc").Find(&comments)
	} else {
		db.Where(&Comment{BoardID: boardId}).Order("created_at desc").Find(&comments)
	}

	c.HTML(200, "board.html", gin.H{
		"board" : boardObj,
		"ID":       boardId,
		"comments": comments,
	})
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

	router.GET("/", getAllBoards)
	router.GET("/new/board", func(c *gin.Context) {
		c.HTML(200, "create_board.html", gin.H{"title": "new board"})
	})
	router.POST("/new/board/post", createBoard)
	router.GET("/board/:id", getAllBoardComments)
	router.POST("/board/:id/comment", createBoardComment)

	router.Run()
}
