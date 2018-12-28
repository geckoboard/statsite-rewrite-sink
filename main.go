package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/geckoboard/statsite-rewrite-sink/regexengine"
	"github.com/geckoboard/statsite-rewrite-sink/sinkformatter"
)

func loadExample(name string) (io.Reader, error) {
	f := filepath.Join("./examples", name)

	return os.Open(f)
}

func main() {
	r, _ := loadExample("consul.in.dump")

	regexengine.Stream(r, os.Stdout, rules, sinkformatter.Librato)
}
