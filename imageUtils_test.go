package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConcatImage(t *testing.T) {
	wd, _ := os.Getwd()
	ConcatImage(filepath.Join(wd, "testdata", "pngs"))
}
