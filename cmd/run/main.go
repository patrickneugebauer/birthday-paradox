package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const infile = "images.tsv"
const outfile = "runs.tsv"
const comma = '\t'

func main() {
	solutionsFile, err := os.OpenFile(infile, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("read: " + infile)
	defer solutionsFile.Close()
	r := csv.NewReader(solutionsFile)
	r.Comma = comma
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	outFile, err := os.Create("runs.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	w := csv.NewWriter(outFile)
	w.Comma = comma
	defer w.Flush()
	w.Write([]string{"lang", "solution", "image", "iterations", "sample-size", "percent", "seconds", "ips"})
	results := make([][]string, 0, len(records))
	// skip header
	for _, v := range records[1:] {
		fmt.Println(v)
		lang, solution, image := v[0], v[1], v[2]
		iterations := 1000
		command := fmt.Sprintf("docker run --rm %s %d", image, iterations)
		fields := strings.Fields(command)
		fmt.Println(strings.Join(fields, " "))
		c := exec.Command(fields[0], fields[1:]...)
		data, err := c.Output()
		if err != nil {
			log.Fatal(err)
		}
		seconds := getResults(data)
		secondsNum, err := strconv.ParseFloat(seconds, 64)
		if err != nil {
			log.Fatal(err)
		}
		ips := float64(iterations) / secondsNum
		row := []string{lang, solution, image, strconv.Itoa(iterations), seconds, ips}
		fmt.Println(row)
		w.Write(row)
		w.Flush() // needed to get it to write now
		results = append(results, row)
		fmt.Println()
	}
}

func getResults(rawOutput []byte) string {
	// "iterations: xxxx\nsample-size: xx\npercent: xx.xx\nseconds: x.xxxxxx\n"
	outputString := strings.TrimSpace(string(rawOutput))
	fmt.Printf("%q\n", outputString)
	lines := strings.Split(outputString, "\n")
	if len(lines) < 4 {
		log.Fatal("fourth row not found!")
	}
	secondsLine := lines[3]
	_, after, found := strings.Cut(secondsLine, ": ")
	if !found {
		fmt.Println(outputString)
		fmt.Println(lines)
		log.Fatal("output row string not split")
	}
	seconds := strings.TrimSpace(after)
	return seconds
}
