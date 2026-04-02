package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	solutionsFile, err := os.OpenFile("solutions.tsv", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer solutionsFile.Close()
	r := csv.NewReader(solutionsFile)
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Print(records)

	outFile, err := os.Create("images.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	w := csv.NewWriter(outFile)
	w.Comma = '\t'
	defer w.Flush()
	w.Write([]string{"lang", "solution", "image"})
	results := make([][]string, 0, len(records))
	// skip header
	for _, v := range records[1:] {
		fmt.Println(v)
		lang, dockerfile, solution := v[0], v[1], v[2]
		dir := "./solutions/" + lang
		fname := dir + "/" + dockerfile
		image := "bday/" + solution
		command := fmt.Sprintf("docker build -f %s %s -t %s", fname, dir, image)
		fields := strings.Fields(command)
		fmt.Println(strings.Join(fields, " "))
		c := exec.Command(fields[0], fields[1:]...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		// _, err := c.Output()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(string(out))
		row := []string{lang, solution, image}
		fmt.Println(row)
		w.Write(row)
		w.Flush() // needed to get it to write now
		results = append(results, row)
		fmt.Println()
	}
	fmt.Println(results)
	// w.WriteAll(results)
}
