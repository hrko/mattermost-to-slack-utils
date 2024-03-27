package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("Usage: go run main.go <zipfile> <attachment-dir> <input.jsonl> <output.jsonl>")
	}

	zipFileName := os.Args[1]
	idToFilename := getIdToFilename(zipFileName)
	for id, filename := range idToFilename {
		fmt.Printf("%s: %s\n", id, filename)
	}

	dir := os.Args[2]
	moveFiles(dir, idToFilename)

	jsonlInputPath := os.Args[3]
	jsonlOutputPath := os.Args[4]
	replacePathInJsonl(jsonlInputPath, jsonlOutputPath, idToFilename)
}

func getIdToFilename(zipFileName string) map[string]string {
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	idToFilename := make(map[string]string)
	for _, file := range zipFile.File {
		re := regexp.MustCompile(`^__uploads/([A-Z0-9]+)/(.+)$`)
		matches := re.FindStringSubmatch(file.Name)
		if len(matches) != 3 {
			continue
		}
		idToFilename[matches[1]] = matches[2]
	}

	return idToFilename
}

// dir: Path to "bulk-export-attachments" directory
func moveFiles(dir string, idToFilename map[string]string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		re := regexp.MustCompile(`^([A-Z0-9]+)_.+$`)
		matches := re.FindStringSubmatch(info.Name())
		if len(matches) != 2 {
			return nil
		}
		id := matches[1]
		filename, ok := idToFilename[id]
		if !ok {
			return nil
		}
		newPath := filepath.Join(dir, id, filename)
		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return err
		}
		if err := os.Rename(path, newPath); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func replacePathInJsonl(jsonlInputPath, jsonlOutputPath string, idToFilename map[string]string) {
	inputFile, err := os.Open(jsonlInputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(jsonlOutputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, `{"type":"post"`) {
			re := regexp.MustCompile(`"path":"bulk-export-attachments/([A-Z0-9]+)_.+?"`)
			line = re.ReplaceAllStringFunc(line, func(path string) string {
				id := re.FindStringSubmatch(path)[1]
				filename, ok := idToFilename[id]
				if !ok {
					return path
				}
				return fmt.Sprintf(`"path":"bulk-export-attachments/%s/%s"`, id, filename)
			})
		}
		if _, err := outputFile.WriteString(line + "\n"); err != nil {
			log.Fatal(err)
		}
	}
}
