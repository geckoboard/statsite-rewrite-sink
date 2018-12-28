package main

import (
	"io"
	"os"
	"path/filepath"
)

func loadExample(name string) (io.Reader, error) {
	f := filepath.Join("./examples", name)

	return os.Open(f)
}

func main() {
	r, _ := loadExample("consul.in.dump")

	RegexScanner(r, os.Stdout, rules)
}
