package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/knyazev-ro/vulcan/config"
)

var GlobalLimit = make(chan struct{}, 64)
var DB *sql.DB

func Init() {
	config := config.GetConfig()
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.Driver, config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("pgx", dsn) // pgx через database/sql
	if err != nil {
		panic(err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		panic(err)
	}
	DB = db
}
