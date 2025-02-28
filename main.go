package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Load dotenv file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create an HTTP server and Listen on port 7777
	http.HandleFunc("/obsidian/", HelloHandler)
	http.HandleFunc("/obsidian/ask", AskHandler)
	http.HandleFunc("/obsidian/upload", UploadHandler)

	log.Fatal(http.ListenAndServe(":7777", nil))
	log.Println("Server started on port 7777")
}
