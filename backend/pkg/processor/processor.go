package processor

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

var logs []LogEntry

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	if strings.HasSuffix(file.Filename, ".tar.gz") {
		if err := extractTarGz(f); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .tar.gz supported"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "uploaded"})
}

func extractTarGz(r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		path := filepath.Join(os.TempDir(), hdr.Name)
		os.MkdirAll(filepath.Dir(path), 0755)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return err
		}
		f.Close()
		parseFile(path)
	}
	return nil
}

func parseFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var timestampStr, level, msg string
		if strings.HasPrefix(line, "[") {
			if idx := strings.Index(line, "]"); idx != -1 {
				timestampStr = line[1:idx]
				rest := strings.TrimSpace(line[idx+1:])
				parts := strings.SplitN(rest, " ", 2)
				if len(parts) == 2 {
					level = parts[0]
					msg = parts[1]
				}
			}
		}
		if level == "" {
			level = "INFO"
			msg = line
		}
		ts, _ := time.Parse(time.RFC3339, timestampStr)
		logs = append(logs, LogEntry{Timestamp: ts, Level: level, Message: msg})
	}
}

func AnalyzeHandler(c *gin.Context) {
	level := c.Query("level")
	var result []LogEntry
	for _, l := range logs {
		if level == "" || strings.EqualFold(l.Level, level) {
			result = append(result, l)
		}
	}
	c.JSON(http.StatusOK, result)
}
