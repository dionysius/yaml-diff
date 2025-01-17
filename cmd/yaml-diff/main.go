package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sters/yaml-diff/yamldiff"
)

func main() {
	ignoreEmptyFields := flag.Bool("ignore-empty-fields", false, "Ignore empty field")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: yaml-diff file1 file2")
		os.Exit(1)
	}
	file1 := args[0]
	file2 := args[1]

	yamls1, err := yamldiff.Load(load(file1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	yamls2, err := yamldiff.Load(load(file2))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	opts := []yamldiff.DoOptionFunc{}
	if *ignoreEmptyFields {
		opts = append(opts, yamldiff.EmptyAsNull())
	}

	fmt.Printf("--- %s\n+++ %s\n\n", file1, file2)
	for _, diff := range yamldiff.Do(yamls1, yamls2, opts...) {
		fmt.Println(diff.Dump())
	}

	fmt.Print()
}

func load(f string) string {
	file, err := os.Open(f)
	defer func() { _ = file.Close() }()
	if err != nil {
		log.Printf("%+v, %s", err, f)

		return ""
	}

	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("%+v, %s", err, f)

		return ""
	}

	return string(b)
}
