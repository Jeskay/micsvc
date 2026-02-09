package main

import (
	"log"

	"github.com/Jeskay/micsvc/internal/db"
	transport "github.com/Jeskay/micsvc/internal/transport/http"
	"github.com/Jeskay/micsvc/internal/user"
)

func main() {
	svc := user.NewUserService(db.NewUserStorage())
	transport.RunHTTP("localhost:9090", svc, func(err error) {
		log.Println(err)
	})
}
