package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

type exampleTestcase struct {
	input  string
	output string
}

func listExamples(t *testing.T) []exampleTestcase {
	files, err := filepath.Glob("examples/*.in.dump")
	if err != nil {
		t.Fatal(err)
	}

	testcases := make([]exampleTestcase, 0, len(files))

	for _, f := range files {
		c := exampleTestcase{
			input:  f,
			output: strings.Replace(f, "in.dump", "out.dump", -1),
		}
		testcases = append(testcases, c)
	}

	return testcases
}

func TestExamples(t *testing.T) {
	for _, input := range listExamples(t) {
		inputReader, err := os.Open(input.input)
		if err != nil {
			t.Fatal(t)
		}
		outputWriter := strings.Builder{}

		RegexScanner(inputReader, &outputWriter, rules)

		output := outputWriter.String()

		if os.Getenv("UPDATE_GOLDEN_OUTPUT") != "" {
			ioutil.WriteFile(input.output, []byte(output), os.ModePerm)
			continue
		}

		expected, err := ioutil.ReadFile(input.output)
		if err != nil {
			t.Fatal(err)
		}

		if string(expected) != string(output) {
			diff := diff.LineDiff(string(expected), string(output))
			t.Errorf("unexpected diff: %s", diff)
		}
	}
}
