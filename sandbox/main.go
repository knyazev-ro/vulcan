package main

import (
	"github.com/knyazev-ro/vulcan/orm/db"
)

func main() {
	db.Init()
	RealExampleORM()
}
