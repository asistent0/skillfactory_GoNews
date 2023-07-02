package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db1 := memdb.New()

	// Реляционная БД PostgresSQL.
	constr, err := constrPostgres()
	if err != nil {
		log.Fatal(err)
	}
	db2, err := postgres.New(constr)
	if err != nil {
		log.Fatal(err)
	}
	// Документная БД MongoDB.
	constr, err = constrMongo()
	if err != nil {
		log.Fatal(err)
	}
	db3, err := mongo.New(constr)
	if err != nil {
		log.Fatal(err)
	}
	_, _, _ = db1, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db1

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	err = http.ListenAndServe(":8080", srv.api.Router())
	if err != nil {
		return
	}
}

func constrPostgres() (string, error) {
	user, exists := os.LookupEnv("POSTGRES_USER")
	if !exists {
		return "", errors.New("not exist postgres user")
	}
	pwd, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if !exists {
		return "", errors.New("not exist postgres password")
	}
	pdb, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		return "", errors.New("not exist postgres db")
	}
	port, exists := os.LookupEnv("POSTGRES_PORT")
	if !exists {
		return "", errors.New("not exist postgres port")
	}
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		return "", errors.New("not exist postgres host")
	}
	constr := "postgres://" + user + ":" + pwd + "@" + host + ":" + port + "/" + pdb
	return constr, nil
}

func constrMongo() (string, error) {
	user, exists := os.LookupEnv("MONGO_USER")
	if !exists {
		return "", errors.New("not exist mongo user")
	}
	pwd, exists := os.LookupEnv("MONGO_PASSWORD")
	if !exists {
		return "", errors.New("not exist mongo password")
	}
	port, exists := os.LookupEnv("MONGO_PORT")
	if !exists {
		return "", errors.New("not exist mongo port")
	}
	host, exists := os.LookupEnv("MONGO_HOST")
	if !exists {
		return "", errors.New("not exist mongo host")
	}
	constr := "mongodb://" + user + ":" + pwd + "@" + host + ":" + port + "/"
	return constr, nil
}
