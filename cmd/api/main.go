package main

import (
	"fmt"
	"recruit-sys/internal/server"
)

func main() {

	server := server.NewServer()

	fmt.Println("server running at : ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
