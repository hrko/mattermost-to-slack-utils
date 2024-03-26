package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ChannelInfo stores information about a channel from input1
type ChannelInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: program <input1.json> <input2.jsonl> <output.jsonl>")
		return
	}

	input1Path := os.Args[1]
	input2Path := os.Args[2]
	outputPath := os.Args[3]

	// Read and parse input1.json
	idToName, err := readInput1(input1Path)
	if err != nil {
		panic(err)
	}

	// Process input2.jsonl and write to output.jsonl
	if err := processInput2(input2Path, outputPath, idToName); err != nil {
		panic(err)
	}
}

func readInput1(filePath string) (map[string]string, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var channelInfos []ChannelInfo
	if err := json.Unmarshal(fileContent, &channelInfos); err != nil {
		return nil, err
	}

	idToName := make(map[string]string)
	for _, info := range channelInfos {
		idToName[info.ID] = info.Name
	}

	return idToName, nil
}

func processInput2(inputPath, outputPath string, idToName map[string]string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, `{"type":"channel"`) {
			for id, name := range idToName {
				// Replace display_name if it matches the lowercase id
				if strings.Contains(line, fmt.Sprintf(`"display_name":"%s"`, strings.ToLower(id))) {
					line = strings.ReplaceAll(line, fmt.Sprintf(`"display_name":"%s"`, strings.ToLower(id)), fmt.Sprintf(`"display_name":"%s"`, name))
					break
				}
			}
		}
		if _, err := outputFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
