package main

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/geckoboard/statsite-rewrite-sink/regexengine"
	"github.com/geckoboard/statsite-rewrite-sink/sinkformatter"
	"github.com/spf13/cobra"
)

func main() {
	cmdRewrite := &cobra.Command{
		Use:   "rewrite [input-file]",
		Short: "Rewrites the metrics in [input-file]",
		Long:  "Rewrites the metrics in [input-file] and writes them to STDOUT. Set [input-file] to `-` to read from STDIN.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sourceFile := args[0]
			input := os.Stdin

			var err error
			if sourceFile != "-" {
				input, err = os.Open(sourceFile)
				if err != nil {
					panic(err)
				}
			}

			rewriteToWriter(input, os.Stdout)
		},
	}
	cmdStream := &cobra.Command{
		Use:   "stream [file] [command] ([arg]...)",
		Short: "Rewrite the contents of [file] and stream results to [command]",
		Long:  "Rewrite the contents of [file] and stream results to [command]. Pass `-` as [file] to read from stdin",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			sourceFile := args[0]
			input := os.Stdin

			var err error
			if sourceFile != "-" {
				input, err = os.Open(sourceFile)
				if err != nil {
					panic(err)
				}
			}

			rewriteToCommand(input, args[1:])
		},
	}

	rootCmd := &cobra.Command{
		Use:   "statsite-rewrite-sink",
		Short: "Rewrite statsite metrics to include metric tags",
	}
	rootCmd.AddCommand(cmdRewrite, cmdStream)
	rootCmd.Execute()

}

func rewriteToWriter(in io.ReadCloser, out io.Writer) {
	regexengine.Stream(in, out, rules, sinkformatter.Librato)
}

func rewriteToCommand(in io.ReadCloser, command []string) {
	var wg sync.WaitGroup

	cmdIn, rewriteOut := io.Pipe()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rewriteToWriter(in, rewriteOut)
		// The semantics of `io.Pipe()` mean that invoking `Write()` on
		// the writer only returns once the reader has consumed the
		// bytes that were written. This means that when the
		// rewriter has finished we know that the command has read
		// everything we've written.
		//
		// Closing the writer indicates to `exec.Cmd` that it should stop
		// the process it's managing.
		rewriteOut.Close()
	}()

	args := []string{}
	if len(command) > 1 {
		args = command[1:]
	}
	cmd := exec.Command(command[0], args...)
	cmd.Stdin = cmdIn
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	wg.Wait()
}
