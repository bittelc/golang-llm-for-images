package input

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang-ai-server/logger"
	"image/png"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gen2brain/go-fitz"
)

func GetUserInput() (string, []string, error) {
	fmt.Print("User prompt: ")
	reader := bufio.NewReader(os.Stdin)

	prompt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Enable image attachment
	fmt.Print("Path to image (optional): ")
	imageInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Image path(s) received", "pathString", imageInput)

	// Parse comma-separated image paths
	var imageByteStrings []string
	imageInput = strings.TrimSuffix(imageInput, "\n") // Remove trailing newline
	if len(imageInput) < 1 || imageInput == "" {
		slog.Debug("No image input provided")
		return prompt, nil, nil
	}

	// Split by comma and trim whitespace
	paths := strings.Split(imageInput, ",")
	slog.Debug("Split image paths", "path_count", len(paths), "paths", paths)

	if len(paths) > 5 {
		return "", nil, fmt.Errorf("too many images provided")
	}

	for _, path := range paths {
		trimmedPath := strings.TrimSpace(path)

		if trimmedPath == "" {
			continue // empty path provided
		}

		mime, err := detectFileType(trimmedPath)
		if err != nil {
			log.Fatal(err)
		}
		var encodedImage string
		switch mime {
		case "application/pdf":
			slog.Info("This is a PDF", "trimmedPath", trimmedPath)
			pageImages, err := convertPdfToImages(trimmedPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to encode pdf into images: %s: error %v", path, err)
			}
			return prompt, pageImages, nil
		case "image/png", "image/jpeg", "image/jpg":
			slog.Info("png, jpeg, or jpg image detected", "fileType", mime)
			//TODO process jpeg and other images
		default:
			slog.Info("Other type of file detected:", "fileType", mime)
			return "", nil, fmt.Errorf("cannot parse files other than Pdf, jpeg, png: %s", path)
		}
		if err != nil {
			return "", nil, fmt.Errorf("failed to encode file %s: %v", path, err)
		}
		slog.Debug("Successfully encoded image", "path", trimmedPath, "encoded_length", len(encodedImage))
	}

	logger.LogProcessingStep("complete_user_input", map[string]interface{}{
		"prompt_length": len(prompt),
		"image_count":   len(imageByteStrings),
	})
	return prompt, imageByteStrings, err
}

func convertPdfToImages(filePath string) ([]string, error) {
	var pdfInB64 []string
	doc, err := fitz.New(filePath)
	if err != nil {
		return nil, err
	}
	defer doc.Close()
	for n := range doc.NumPage() {
		img, err := doc.Image(n)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			log.Printf("Failed to encode page %d to PNG: %v", n, err)
			continue
		}

		pageInB64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		pdfInB64 = append(pdfInB64, pageInB64)
	}
	return pdfInB64, nil
}

func detectFileType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512) // only need up to the first 512 bytes
	n, err := f.Read(buf)
	if err != nil {
		return "", err
	}

	kind := http.DetectContentType(buf[:n])
	return kind, nil
}
