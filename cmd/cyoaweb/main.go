package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA application on")
	filename := flag.String("file", "gopher.json", "The JSON file with the Choose Your Own Adventure story.")
	flag.Parse()
	fmt.Printf("Using the story in %s\n", *filename)

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, error := cyoa.JSONStory(file)
	if error != nil {
		panic(error)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
