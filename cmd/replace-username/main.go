package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type UserInfo struct {
	NameOld string `json:"name_old"`
	NameNew string `json:"name_new"`
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

	var userInfos []UserInfo
	if err := json.Unmarshal(fileContent, &userInfos); err != nil {
		return nil, err
	}

	oldToNewName := make(map[string]string)
	for _, info := range userInfos {
		oldToNewName[info.NameOld] = info.NameNew
	}

	return oldToNewName, nil
}

func processInput2(inputPath, outputPath string, oldToNewName map[string]string) error {
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
		if strings.HasPrefix(line, `{"type":"user"`) {
			// Process `user` type
			// Replace username
			username := regexp.MustCompile(`"username":"[a-z_\-\.]+"`)
			line = username.ReplaceAllStringFunc(line, func(username string) string {
				username = username[12 : len(username)-1] // Remove '"username":"' and '"'
				if newName, ok := oldToNewName[username]; ok {
					return fmt.Sprintf(`"username":"%s"`, newName)
				}
				return fmt.Sprintf(`"username":"%s"`, username)
			})
		} else if strings.HasPrefix(line, `{"type":"post"`) {
			// Process `post` type
			// Replace post author
			author := regexp.MustCompile(`"user":"[a-z_\-\.]+"`)
			line = author.ReplaceAllStringFunc(line, func(author string) string {
				author = author[8 : len(author)-1] // Remove '"user":"' and '"'
				if newName, ok := oldToNewName[author]; ok {
					return fmt.Sprintf(`"user":"%s"`, newName)
				}
				return fmt.Sprintf(`"user":"%s"`, author)
			})
			// Replace mentions
			mentions := regexp.MustCompile(`@[a-z_\-\.]+`)
			line = mentions.ReplaceAllStringFunc(line, func(mention string) string {
				mention = mention[1:] // Remove '@' prefix
				if newName, ok := oldToNewName[mention]; ok {
					return "@" + newName
				}
				return "@" + mention
			})
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
