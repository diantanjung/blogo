package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_userHttpDelivery "github.com/diantanjung/blogo/user-service/user/delivery/http"
	_userMiddleware "github.com/diantanjung/blogo/user-service/user/delivery/http/middleware"
	_userRepo "github.com/diantanjung/blogo/user-service/user/repository/psql"
	_userUcase "github.com/diantanjung/blogo/user-service/user/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _userMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	repo := _userRepo.NewPsqlUserRepository(db)
	us := _userUcase.NewUserUsecase(repo)

	_userHttpDelivery.NewUsersHandler(e, us)

	log.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
