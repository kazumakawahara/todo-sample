package main

import (
	"log"

	"github.com/kazumakawahara/todo-sample/infrastructure/router"
)

func main() {
	if err := router.Run(); err != nil {
		log.Fatalln(err)
	}
}
