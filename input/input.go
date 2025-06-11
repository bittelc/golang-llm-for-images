package input

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang-ai-server/logger"
	"image/png"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gen2brain/go-fitz"
)

func GetUserInput() (string, []string, error) {
	logger.LogProcessingStep("start_user_input", map[string]interface{}{
		"function": "GetUserInput",
	})

	// Get prompt text
	fmt.Print("User prompt: ")
	reader := bufio.NewReader(os.Stdin)
	slog.Debug("Created buffered reader for stdin")

	prompt, err := reader.ReadString('\n')
	if err != nil {
		logger.LogError("read_prompt", err, map[string]interface{}{
			"source": "stdin",
		})
		log.Fatal(err)
	}
	logger.LogUserInput("prompt", len(prompt))

	// Enable image attachment
	fmt.Print("Path to images, seperated by commas, limit of 5 (optional): ")
	imageInput, err := reader.ReadString('\n')
	if err != nil {
		logger.LogError("read_image_input", err, map[string]interface{}{
			"source": "stdin",
		})
		log.Fatal(err)
	}
	logger.LogUserInput("image_paths", len(imageInput))

	// Parse comma-separated image paths
	var imageByteStrings []string
	if len(imageInput) > 1 { // Check if there's actual input beyond just newline
		logger.LogProcessingStep("process_image_paths", map[string]interface{}{
			"raw_input_length": len(imageInput),
		})
		imageInput = imageInput[:len(imageInput)-1] // Remove trailing newline
		if imageInput != "" {
			// Split by comma and trim whitespace
			paths := strings.Split(imageInput, ",")
			slog.Debug("Split image paths", "path_count", len(paths), "paths", paths)

			if len(paths) > 5 {
				return "", nil, fmt.Errorf("too many images provided")
			}

			for i, path := range paths {
				trimmedPath := strings.TrimSpace(path)
				logger.LogProcessingStep("encode_image", map[string]interface{}{
					"index": i,
					"path":  trimmedPath,
				})

				if trimmedPath != "" {
					mime, err := detectFileType(trimmedPath)
					if err != nil {
						// handle error
						log.Fatal(err)
					}
					var encodedImage string
					switch mime {
					case "application/pdf":
						slog.Info("This is a PDF")
						pageImages, err := convertPdfToImages(trimmedPath)
						if err != nil {
							return "", nil, fmt.Errorf("failed to encode pdf into images: %s: error %v", path, err)
						}
						return prompt, pageImages, nil

					case "image/png", "image/jpeg":
						fmt.Println("This is a PNG/JPEG image")
					default:
						fmt.Println("Other type:", mime)
						return "", nil, fmt.Errorf("cannot parse files other than Pdf, jpeg, png: %s", path)
					}
					if err != nil {
						logger.LogError("encode_image", err, map[string]interface{}{
							"path":  trimmedPath,
							"index": i,
						})
						return "", nil, fmt.Errorf("failed to encode image %s: %v", path, err)
					}
					slog.Debug("Successfully encoded image", "path", trimmedPath, "encoded_length", len(encodedImage))
				}
			}
		}
	} else {
		slog.Debug("No image input provided")
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
		f, err := os.Create(fmt.Sprintf("page_%d.png", n+1))
		if err != nil {
			return nil, err
		}
		png.Encode(f, img)
		f.Close()

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

func encodeImageToStringTryAgain(imagePath string) (string, error) {
	logger.LogFileOperation("read_image", imagePath, 0)

	bytes, err := os.ReadFile(imagePath)
	if err != nil {
		logger.LogError("read_image_file", err, map[string]interface{}{
			"path": imagePath,
		})
		return "", err
	}

	logger.LogFileOperation("encode_image", imagePath, int64(len(bytes)))
	pdfString := base64.StdEncoding.EncodeToString(bytes)

	// Now you have the entire PDF as a string
	slog.Info("Successfully encoded file to base64", "path", imagePath, "encoded_length", len(pdfString))
	fmt.Printf("PDF content as string length: %d\n", len(pdfString))

	return pdfString, err
}

// This works if the incoming pdfString is from
// pdfString := base64.StdEncoding.EncodeToString(bytes)
func writeEncodedPdfToFile(pdfString string) {
	decoded, err := base64.StdEncoding.DecodeString(pdfString)
	if err != nil {
		log.Fatalf("Failed to decode pdfString: %v", err)
	}
	err = os.WriteFile("encoded.pdf", decoded, 0644)
	if err != nil {
		log.Fatalf("Failed to write encoded file: %v", err)
	}
	slog.Info("Successfully wrote file")
}

func encodeImageToString(imagePath string) (string, error) {
	logger.LogFileOperation("open_image", imagePath, 0)

	f, err := os.Open(imagePath)
	if err != nil {
		logger.LogError("open_image_file", err, map[string]interface{}{
			"path": imagePath,
		})
		return "", err
	}
	defer f.Close()
	slog.Debug("Successfully opened file", "path", imagePath)

	// Get file info to allocate exact buffer size
	stat, err := f.Stat()
	if err != nil {
		logger.LogError("get_file_stats", err, map[string]interface{}{
			"path": imagePath,
		})
		return "", err
	}

	size := stat.Size()
	logger.LogFileOperation("allocate_buffer", imagePath, size)
	buf := make([]byte, size)

	// Read the entire file at once (you can also use io.ReadAll)
	bytesRead, err := io.ReadFull(f, buf)
	if err != nil {
		logger.LogError("read_full_file", err, map[string]interface{}{
			"path":           imagePath,
			"expected_bytes": size,
			"actual_bytes":   bytesRead,
		})
		return "", err
	}
	slog.Debug("Successfully read full file", "path", imagePath, "bytes_read", bytesRead)

	// Convert []byte to string
	// pdfString := string(buf)
	pdfString := base64.StdEncoding.EncodeToString(buf)

	// Now you have the entire PDF as a string
	slog.Debug("Successfully encoded file to base64", "path", imagePath, "encoded_length", len(pdfString))
	fmt.Printf("PDF content as string length: %d\n", len(pdfString))
	return pdfString, err
}

func encodeImageToBase64(imagePath string) (string, error) {
	logger.LogFileOperation("open_image_b64", imagePath, 0)

	file, err := os.Open(imagePath)
	if err != nil {
		logger.LogError("open_image_file_b64", err, map[string]interface{}{
			"path": imagePath,
		})
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()
	slog.Debug("Successfully opened image file", "path", imagePath)

	// Read the entire file
	imageData, err := io.ReadAll(file)
	if err != nil {
		logger.LogError("read_image_file_b64", err, map[string]interface{}{
			"path": imagePath,
		})
		return "", fmt.Errorf("failed to read image file: %v", err)
	}
	logger.LogFileOperation("encode_image_b64", imagePath, int64(len(imageData)))

	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	slog.Debug("Successfully encoded image to base64", "path", imagePath, "encoded_length", len(encodedImage))
	return encodedImage, err
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
