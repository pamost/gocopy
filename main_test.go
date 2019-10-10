package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

const src, dst = "src.txt", "dst.txt"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generationSrcFile() {
	buf := make([]byte, 1024)
	for i := range buf {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	file, err := os.Create(src)
	if err != nil {
		fmt.Printf("failed to create src file %v", file)
		return
	}
	if _, err := file.Write(buf); err != nil {
		fmt.Printf("failed to written in src file %v", file)
		return
	}
	file.Close()

}

func TestGoCopy(t *testing.T) {
	testTable := []struct {
		from   string
		to     string
		offset int
		limit  int
	}{
		{src, dst, 0, 0},
		{src, dst, 0, 32},
		{src, dst, 64, 256},
		{src, dst, 512, 0},
	}

	generationSrcFile()

	for i, item := range testTable {
		err := goCopy(item.from, item.to, item.offset, item.limit)
		if err != nil {
			t.Errorf("error on item %d %v", i, err)
		}
	}

	os.Remove(dst)
	os.Remove(src)
}
