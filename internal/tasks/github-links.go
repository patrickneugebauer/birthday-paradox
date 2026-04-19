package tasks

import (
	"fmt"
	"os"
)

func CheckGithubLinks() error {
	_, err := os.Stat(githubLinksFile)
	// var output *os.File
	if os.IsNotExist(err) {
		makeNewFile()
	} else {
		fmt.Println("file was found")
		// output, err = os.OpenFile(githubLinksFile, os.O_RDWR, 0755)
		// if err != nil {
		// 	return fmt.Errorf("failed to create output file: %w", err)
		// }
		// defer output.Close()
	}
	return nil
}

func makeNewFile() error {
	output, err := os.Create(githubLinksFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	fmt.Printf("created file: %s\n", githubLinksFile)
	defer output.Close()
	languages, err := ReadDockerfileList()
	if err != nil {
		return err
	}
	output.WriteString("language,url\n")
	for _, v := range languages {
		text := fmt.Sprintf("%s,\n", v.Language)
		output.WriteString(text)
	}
	return nil
}
