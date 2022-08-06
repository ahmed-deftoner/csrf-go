package main

import (
	"log"

	"github.com/ahmed-deftoner/csrf-go/db"
	"github.com/ahmed-deftoner/csrf-go/server"
	myjwt "github.com/ahmed-deftoner/csrf-go/server/middleware/myJwt"
)

var host = "localhost"
var port = "9000"

func main() {
	db.DBInit()
	jwtErr := myjwt.JWTInit()

	if jwtErr != nil {
		log.Println("error initializing jwt")
		log.Fatal(jwtErr)
	}

	serverErr := server.StartServer(host, port)
	if serverErr != nil {
		log.Println("error initializing server")
		log.Fatal(serverErr)
	}
}
