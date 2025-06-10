package input

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func GetUserInput() (string, []string, error) {

	// Get prompt text
	fmt.Print("User prompt: ")
	reader := bufio.NewReader(os.Stdin)
	prompt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Enable image attachment
	fmt.Print("Path to images, seperated by commas, limit of 5 (optional): ")
	imageInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Parse comma-separated image paths
	var imageByteStrings []string
	if len(imageInput) > 1 { // Check if there's actual input beyond just newline
		imageInput = imageInput[:len(imageInput)-1] // Remove trailing newline
		if imageInput != "" {
			// Split by comma and trim whitespace
			paths := strings.Split(imageInput, ",")
			if len(paths) > 5 {
				return "", nil, fmt.Errorf("too many images provided")
			}
			for _, path := range paths {
				trimmedPath := strings.TrimSpace(path)
				if trimmedPath != "" {
					// encodedImage, err := encodeImageToBase64(path)
					encodedImage, err := encodeImageToStringTryAgain(trimmedPath)
					if err != nil {
						return "", nil, fmt.Errorf("failed to encode image %s: %v", path, err)
					}
					imageByteStrings = append(imageByteStrings, encodedImage)
				}
			}
		}
	}

	return prompt, imageByteStrings, err
}

func encodeImageToStringTryAgain(imagePath string) (string, error) {
	bytes, err := os.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}
	pdfString := base64.StdEncoding.EncodeToString(bytes)

	// Now you have the entire PDF as a string
	fmt.Printf("PDF content as string length: %d\n", len(pdfString))
	return pdfString, err
}

func encodeImageToString(imagePath string) (string, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get file info to allocate exact buffer size
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	size := stat.Size()
	buf := make([]byte, size)

	// Read the entire file at once (you can also use io.ReadAll)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}

	// Convert []byte to string
	// pdfString := string(buf)
	pdfString := base64.StdEncoding.EncodeToString(buf)

	// Now you have the entire PDF as a string
	fmt.Printf("PDF content as string length: %d\n", len(pdfString))
	return pdfString, err
}

func encodeImageToBase64(imagePath string) (string, error) {

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Read the entire file
	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	return encodedImage, err
}
