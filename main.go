package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_userHttpDeliver "github.com/budhilaw/blog/user/delivery/http"
	_userRepo "github.com/budhilaw/blog/user/repository"
	_userUcase "github.com/budhilaw/blog/user/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	// Variables
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", conn, val.Encode())
	dbConn, err := sql.Open("mysql", dsn)
	Check(err)

	err = dbConn.Ping()
	Check(err)

	defer dbConn.Close()

	e := echo.New()
	if viper.GetBool("debug") {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, url=${url}, status=${status}\n",
		}))
	}

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	userRepo := _userRepo.New(dbConn)
	userUcase := _userUcase.New(userRepo, timeoutContext)

	_userHttpDeliver.New(e, userUcase)

	e.Start(viper.GetString("server.address"))
}

// Check : Check if something wrong
func Check(err error) {
	if err != nil && viper.GetBool("debug") {
		log.Fatal(err)
		return
	}
}
