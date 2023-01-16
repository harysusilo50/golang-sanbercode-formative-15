package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"main.go/controllers"
	dbmigrate "main.go/database"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Failed load file environment")
	} else {
		fmt.Println("Success read file environment")
	}
	port, _ := strconv.Atoi(os.Getenv("PGPORT"))
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password =%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		port,
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"))
	DB, err = sql.Open("postgres", psqlInfo)
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB connection failed")
		panic(err)
	} else {
		fmt.Println("DB connection success")
	}

	dbmigrate.DbMigrate(DB)

	defer DB.Close()
	// Router
	router := gin.Default()
	router.GET("/persons", controllers.GetAllPerson)
	router.POST("/persons", controllers.InsertPerson)
	router.PUT("/persons/:id", controllers.UpdatePerson)
	router.DELETE("/persons/:id", controllers.DeletePerson)

	router.Run(":" + os.Getenv("PORT"))
}
