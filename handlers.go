package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func encodeImageToBase64(imagePath string) (string, error) {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imgBytes), nil
}

func processImage(imagePath string) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	base64Image, err := encodeImageToBase64(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		log.Fatal("OPENAI_BASE_URL is not set")
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	ctx := context.Background()

	ext := filepath.Ext(imagePath)
	mimeType := ""
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	default:
		mimeType = "image/jpeg"
	}

	image := openai.ImagePart(fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image))

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant that is going to help me manage my obsidian vault."),
			openai.UserMessage("I am giving you an image of some handwritten note. Convert all text and math equations that you find and turn it into a markdown document that I can store in my obsidian vault. Make sure that all the contents of the image is wrapped in a ```markdown``` block so its easy for me to extract it."),
			openai.UserMessageParts(image),
		}),
		Seed:        openai.Int(1),
		Model:       openai.F("gpt-4o-2024-08-06-moderated"),
		Temperature: openai.Float(0.0),
	})
	if err != nil {
		panic(err)
	}

	content := completion.Choices[0].Message.Content
	// extract the markdown block from the content
	markdownBlock := regexp.MustCompile("(?s)```markdown\\s*(.*?)\\s*```").FindStringSubmatch(content)
	if len(markdownBlock) > 1 {
		// extract the markdown content
		markdownContent := strings.TrimSpace(markdownBlock[1])

		fileName := filepath.Base(imagePath)

		markdownLink := fmt.Sprintf("\n\n![[%s]]", fileName)
		markdownContent += markdownLink
		os.MkdirAll("/app/data", 0755)
		// write the markdown content to a file
		file, err := os.Create(fmt.Sprintf("/app/data/main/%s.md", fileName))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(markdownContent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Markdown file created successfully!")
		// Move the image to the main folder
		imagePath := fmt.Sprintf("/app/data/uploads/%s", fileName)
		err = os.Rename(imagePath, fmt.Sprintf("/app/data/main/attachments/%s", fileName))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No markdown block found in the response.")
	}

}

func ask(question string) {
	// Get the OpenAI API key from the environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		log.Fatal("OPENAI_BASE_URL is not set")
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	ctx := context.Background()

	print("> ")
	println(question)
	println()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:        openai.Int(1),
		Model:       openai.F(openai.ChatModelGPT4o2024_05_13),
		Temperature: openai.Float(0.7),
	})
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}

// helloHandler handles the root route.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// askHandler handles the /ask route.
func AskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		q := r.URL.Query().Get("q")
		ask(q)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// uploadHandler handles the /upload route.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Create the uploads directory if it doesn't exist.
		uploadDir := "/app/data/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, 0755)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating upload directory: %s", err.Error()), http.StatusInternalServerError)
				return
			}
		}

		// Parse the multipart form with a maximum file size of 10MB.
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		// Get the files from the form.
		files := r.MultipartForm.File["images"]
		if files == nil || len(files) == 0 {
			http.Error(w, "No images provided", http.StatusBadRequest)
			return
		}

		// Get the first file from the form.
		fileHeader := files[0]
		if fileHeader == nil {
			http.Error(w, "No image provided", http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Generate a timestamped filename.
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		timestamp = strings.ReplaceAll(timestamp, ":", "-")
		fileExt := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%s_%s%s", strings.TrimSuffix(fileHeader.Filename, fileExt), timestamp, fileExt)
		// Create the file path.
		filePath := filepath.Join(uploadDir, fileName)

		// Create the file.
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the file to the destination.
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying file: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		log.Printf("File '%s' uploaded successfully to '%s'", fileHeader.Filename, filePath)

		fmt.Fprintln(w, "Image uploaded successfully!")

		go processImage(filePath)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
