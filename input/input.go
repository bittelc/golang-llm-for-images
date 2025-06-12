package input

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gen2brain/go-fitz"
	"github.com/nguyenthenguyen/docx"
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

		extension, err := detectFileType(trimmedPath)
		if err != nil {
			log.Fatal(err)
		}
		var encodedImage string
		switch extension {
		case ".pdf":
			slog.Info("This is a PDF", "trimmedPath", trimmedPath)
			pageImages, err := convertPdfToImages(trimmedPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to encode pdf into images: %s: error %v", trimmedPath, err)
			}
			imageByteStrings = append(imageByteStrings, pageImages...)
		case ".png", ".jpg", ".jpeg":
			slog.Info("png, jpg, or jpeg image detected", "fileType", extension)
			//TODO process jpeg and other images
			return "", nil, fmt.Errorf("program not yet able to process image file: %s", trimmedPath)
		case ".docx":
			slog.Info("doc file type detected", "fileType", extension)
			text, err := extractDocTypeText(trimmedPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to extract text from docx file: %s: error %v", trimmedPath, err)
			}
			slog.Info("text extracted from doc", "text", text)
			return "", nil, fmt.Errorf("program not yet able to process docx or doc files: %s", trimmedPath)
		default:
			slog.Error("Other type of file detected. Cannot parse file. Skipping.", "fileType", extension, "filePath", trimmedPath)
			continue
		}
		slog.Debug("Completed encoding of image", "path", trimmedPath, "encoded_length", len(encodedImage))
	}

	return prompt, imageByteStrings, err
}

func convertPdfToImages(filePath string) ([]string, error) {
	var pdfInB64 []string
	doc, err := fitz.New(filePath)
	if err != nil {
		return nil, err
	}
	defer doc.Close()
	if doc.NumPage() > 3 {
		return nil, fmt.Errorf("attached PDF is too long. Documents of more than 3 pages are not allowed.")
	}
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
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return "", err
	}
	return mtype.Extension(), nil
}

func extractDocTypeText(path string) (string, error) {
	doc, err := docx.ReadDocxFile(path)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	text := doc.Editable().GetContent()
	if err != nil {
		return "", err
	}

	return text, nil
}
