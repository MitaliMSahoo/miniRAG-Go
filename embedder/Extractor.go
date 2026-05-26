package embedder

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/dslipak/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/tmc/langchaingo/textsplitter"
)

func ExtractText(file multipart.File, ext string) (string, error) {
	log.Printf("Extracting Text.....")
	switch ext {
	case ".txt", ".md":
		b, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("reading txt: %w", err)
		}
		return string(b), nil

	case ".pdf":
		// read into buffer first — pdf needs ReadSeeker
		buf, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("reading pdf: %w", err)
		}

		r, err := pdf.NewReader(bytes.NewReader(buf), int64(len(buf)))
		if err != nil {
			return "", fmt.Errorf("parsing pdf: %w", err)
		}

		var sb strings.Builder
		for i := 1; i <= r.NumPage(); i++ {
			page := r.Page(i)
			if page.V.IsNull() {
				continue
			}
			text, err := page.GetPlainText(nil)
			if err != nil {
				continue
			}
			sb.WriteString(text)
		}
		return sb.String(), nil

	case ".docx":

		buf, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("reading docx file: %w", err)
		}
		tmpFile, err := os.CreateTemp("", "*.docx")
		if err != nil {
			return "", err
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Write(buf)
		tmpFile.Close()

		r, err := docx.ReadDocxFile(tmpFile.Name())
		if err != nil {
			return "", fmt.Errorf("parsing docx: %w", err)
		}
		defer r.Close()
		return r.Editable().GetContent(), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", ext)
	}

}

func ChunkText(text string) ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(500),   // 500 chars per chunk
		textsplitter.WithChunkOverlap(50), // 50 char overlap between chunks
	)
	return splitter.SplitText(text)
}
