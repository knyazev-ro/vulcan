package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/knyazev-ro/vulcan/config"
)

var GlobalLimit chan struct{}
var DB *sql.DB

func Init() {
	config := config.GetConfig()
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.Driver, config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("pgx", dsn) // pgx через database/sql
	if err != nil {
		panic(err)
	}

	GlobalLimit = make(chan struct{}, config.SemaphoreLimit)

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		panic(err)
	}
	DB = db
}
