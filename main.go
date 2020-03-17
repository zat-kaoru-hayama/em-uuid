package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
)

var rxUUID = regexp.MustCompile(`@UUID@`)

func mains(args []string) error {
	for _, src := range args {
		fd, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("%s: %w", src, err)
		}
		dst := src + ".tmp"
		w, err := os.Create(dst)
		if err != nil {
			fd.Close()
			return fmt.Errorf("%s: %w", dst, err)
		}
		sc := bufio.NewScanner(fd)
		for sc.Scan() {
			line := sc.Text()
			line = rxUUID.ReplaceAllStringFunc(line, func(_ string) string {
				return uuid.New().String()
			})
			fmt.Fprintln(w, line)
		}
		fd.Close()
		if err = w.Close(); err != nil {
			return fmt.Errorf("%s: close: %w", dst, err)
		}
		backup := src + ".bak"
		os.Remove(backup)
		if err = os.Rename(src, backup); err != nil {
			return fmt.Errorf("%s -> %s: %w", src, backup, err)
		}
		if err = os.Rename(dst, src); err != nil {
			return fmt.Errorf("%s -> %s: %w", dst, src, err)
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
