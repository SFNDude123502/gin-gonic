package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type jsonAccount struct {
	Doc        string                 `json:"doc"`
	Login_data map[string]interface{} `json:"login_data"`
	missing    bool
}

var uri string
var last string
var uname string
var pword string
var auth int

func main() {

	err := godotenv.Load("go.env")
	eh(err)
	uri = os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("UNSET 'MONGODB_URI' ENV VAR")
	}

	srv := gin.Default()
	srv.LoadHTMLGlob("templates/*")
	srv.Static("/tmpl", "./templates")
	srv.StaticFile("/favicon.ico", "./favicon.ico")

	store := cookie.NewStore([]byte("passwd"))
	srv.Use(sessions.Sessions("session", store))

	srv.GET("/ping", getPing)
	srv.GET("/register", getRegister)
	srv.POST("/make", make)
	srv.GET("/login", getLogin)
	srv.POST("/fetch", fetch)
	srv.GET("/", getHome)

	srv.Run(":8080")
}

func getRegister(c *gin.Context) {
	last = "register"
	c.HTML(200, "register.go.html", gin.H{"err": c.Query("err")})
}
func getLogin(c *gin.Context) {
	last = "login"
	c.HTML(200, "login.go.html", gin.H{"err": c.Query("err")})
}

func getHome(c *gin.Context) {
	//session := sessions.Default(c)
	if uname == "" || pword == "" || auth == 0 {
		c.Redirect(301, "/login")
	}
	c.HTML(200, "index.go.html", gin.H{
		"username":  uname,
		"password":  pword,
		"authority": auth,
	})
}

func fetch(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")
	//session := sessions.Default(c)
	if last == "" || username == "" || password == "" {
		c.Redirect(301, "/login?err=Missing Username or Password")
		return
	}

	acc := getMongoAcc(username)
	fmt.Println(acc)
	if acc.missing {
		c.Redirect(301, "/login?err=Account Does Not Exist")
		return
	}
	pass, ok := acc.Login_data["password"].(string)
	if !ok {
		panic("invalid password")
	}
	aut, ok := acc.Login_data["auth"].(float64)

	passwd, err := base64.StdEncoding.DecodeString(pass)
	eh(err)
	if ok {
		if string(passwd) == password {
			uname = username
			pword = password
			auth = int(aut)
			c.Redirect(301, "/")

			return
		}
	}

	c.Redirect(301, "/login?err=Wrong Password")
}

func make(c *gin.Context) {
	username, password, auth := c.PostForm("username"), c.PostForm("password"), c.PostForm("auth")
	if username == "" || password == "" {
		c.Redirect(301, "/register?err=Missing Username or Password")
		return
	}
	if auth == "" {
		auth = "1"
	}
	acc := getMongoAcc(username)

	if !acc.missing {
		c.Redirect(301, "/register?err=Username is Taken")
		return
	}

	authority, err := strconv.Atoi(auth)
	eh(err)

	addMongoAcc(username, password, authority)

	c.Redirect(301, "/login")
}

func getPing(c *gin.Context) {
	c.JSON(200, gin.H{"message": "gin is gonic-ing"})
}

func eh(err error) {
	if err != nil {
		panic(err)
	}
}
