package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func Stars() error {
	// Create a context that listens for Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop() // Restore default behavior when Build returns
	infileName := githubLinksFile
	cmdFileName := starCommandsFile
	tempfileName := starsTempResultsFile
	resultsfileName := starResultsFile
	transformer := func(ctx context.Context, scanner *bufio.Scanner, cmdWriter *bufio.Writer, encoder *json.Encoder) error {
		// read data
		inBytes := scanner.Bytes()
		if string(inBytes) == "language,url" {
			fmt.Println("skipping line")
			return nil
		}
		parts := strings.SplitN(string(inBytes), ",", 2)
		if len(parts) < 2 {
			return nil
		}
		language := parts[0]
		url := parts[1]
		// write command
		headers := `"Accept: application/vnd.github+json" -H "Authorization: Bearer $ghtoken"`
		commandText := fmt.Sprintf("curl -i -H %s %s", headers, url)
		if _, err := cmdWriter.WriteString(commandText); err != nil {
			return fmt.Errorf("failed to write command %w", err)
		}
		// run command
		data, err := getStars(url)
		if err != nil {
			return fmt.Errorf("Failed to get stars %v\n", err)
		}

		// write results

		imgInfo := StarResult{Language: language, Stars: data.StargazersCount}
		if err := encoder.Encode(imgInfo); err != nil {
			return fmt.Errorf("failed to marshal json %w", err)
		}
		return nil
	}
	if err := transform(ctx, infileName, transformer, cmdFileName, tempfileName); err != nil {
		return fmt.Errorf("failed to tansform %w", err)
	}
	if err := os.Rename(tempfileName, resultsfileName); err != nil {
		return fmt.Errorf("failed to finalize results: %w", err)
	}
	// log and return
	fmt.Printf("wrote to file: %s\n", resultsfileName)
	return nil
}

func getStars(url string) (GithubRepo, error) {
	githubRepo := GithubRepo{}
	// curl -i -H "Accept: application/vnd.github+json" \
	// 	-H "Authorization: Bearer $ghtoken" \
	// 	https://api.github.com/repos/golang/go
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	token := os.Getenv("ghtoken")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return githubRepo, fmt.Errorf("Error:", err)
	}
	defer resp.Body.Close() // Always close the body
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &githubRepo); err != nil {
		return githubRepo, fmt.Errorf("failed to unmarshal %w", err)
	}
	return githubRepo, nil
}
