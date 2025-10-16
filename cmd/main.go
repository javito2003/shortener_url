package main

import "github.com/javito2003/shortener_url/internal"

func main() {
	server := internal.NewServer()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
