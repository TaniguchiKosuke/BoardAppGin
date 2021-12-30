package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func dbInit() {
	db := dbConnect()
	db.AutoMigrate(&User{}, &Board{}, &Comment{})
}

func dbConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot open the database")
	}

	return db
}

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `form:"username" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}

func signUp(c *gin.Context) {
	var form User
	// バリデーション処理
	if err := c.Bind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
		c.Abort()
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")
		// 登録ユーザーが重複していた場合にはじく処理
		if err := createUser(username, password); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
		}
		c.Redirect(302, "/")
	}
}

func createUser(username string, password string) error {
	db := dbConnect()

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println("cannot generate uuid")
		return err
	}
	id := uuid.String()

	passwordEncrypt, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	passwordErr := db.Create(&User{ID: id, Username: username, Password: string(passwordEncrypt)}).Error
	if passwordErr != nil {
		log.Println("An error occured when creating new User, most likely the username is already used")
		return passwordErr
	}
	return nil
}

func login(c *gin.Context) {
	dbPassword := getUser(c.PostForm("username")).Password
	log.Println(dbPassword)
	formPassword := c.PostForm("password")

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(formPassword)); err != nil {
		log.Println("login failed")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
		c.Abort()
	} else {
		log.Println("login successed")
		c.Redirect(302, "/")
	}
}

func getUser(username string) User {
	db := dbConnect()

	var user User
	db.First(&user, "username = ?", username)
	return user
}

type Board struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string
	Comments []Comment `gorm:"foreignKey:BoardID"`
}

func getAllBoards(c *gin.Context) {
	db := dbConnect()

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
	db := dbConnect()

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic("cannot generate uuid")
	}

	id := uuid.String()
	title := c.PostForm("title")
	db.Create(&Board{ID: id, Title: title, Comments: nil})

	c.Redirect(302, "/")
}

type Comment struct {
	gorm.Model
	ID      string `gorm:"primaryKey"`
	BoardID string
	Content string
}

func createBoardComment(c *gin.Context) {
	db := dbConnect()

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
	db := dbConnect()

	// var board Board
	boardId := c.Param("id")
	// boardObj := db.First(&board, "id = ?", boardId)

	var comments []Comment
	searchParam, queryExist := c.GetQuery("comment")
	if queryExist {
		db.Where("content LIKE ?", "%"+searchParam+"%").Where(&Comment{BoardID: boardId}).Order("created_at desc").Find(&comments)
	} else {
		db.Where(&Comment{BoardID: boardId}).Order("created_at desc").Find(&comments)
	}

	c.HTML(200, "board.html", gin.H{
		"ID":       boardId,
		"comments": comments,
	})
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/**")

	dbInit()

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{})
	})
	router.POST("/signup", signUp)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	router.POST("/login", login)
	router.GET("/", getAllBoards)
	router.GET("/new/board", func(c *gin.Context) {
		c.HTML(200, "create_board.html", gin.H{"title": "new board"})
	})
	router.POST("/new/board/post", createBoard)
	router.GET("/board/:id", getAllBoardComments)
	router.POST("/board/:id/comment", createBoardComment)

	router.Run()
}
