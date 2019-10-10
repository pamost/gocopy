package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

// Copy
func goCopy(from string, to string, offset int, limit int) error {

	// Open file
	src, err := os.Open(from)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("don't specify a source file %v", src)
			return err
		}
		fmt.Printf("can't open the file %v", src)
		return err
	}

	// Info  file
	file, err := src.Stat()
	if err != nil {
		fmt.Printf("not got info about the file %v", file)
		return err
	}
	defer src.Close()

	// Check offset and limit
	switch {
	case limit == 0: // 0, 0 - size ok
		limit = int(file.Size())
	case limit <= offset: // >, < && size ok
		fmt.Printf("-offset value %d bytes must not be greater than or equal to -limit %d bytes", offset, limit)
		return err
	case limit > int(file.Size()): // 0, > && size <
		fmt.Printf("-limit value %d bytes greater than file size %d bytes", limit, int(file.Size()))
		return err
	}
	if offset >= int(file.Size()) { // >, 0 && size <
		fmt.Printf("-offset value %d bytes greater than file size %d bytes", offset, int(file.Size()))
		return err
	}

	buf := make([]byte, limit) // buffer of the desired size

	// Read file
	if read, err := io.ReadFull(src, buf); err != nil {
		fmt.Printf("can't read the file %v %v", src, read)
		return err
	}
	srcReader := bytes.NewReader(buf[offset:])

	// Create file
	dst, err := os.Create(to)
	if err != nil {
		fmt.Printf("not specified file destination %v", dst)
		return err
	}
	defer dst.Close()

	// Copy progress bar
	err = copyProgressBar(dst, srcReader, limit)
	if err != nil {
		fmt.Printf("failed to copy %v", err)
		return err
	}
	return nil
}

// Copy progress bar
func copyProgressBar(dst io.Writer, srcReader io.Reader, limit int) error {

	// Create bar
	p := mpb.New(mpb.WithWidth(64))

	bar := p.AddBar(int64(limit),
		mpb.PrependDecorators(decor.Counters(decor.UnitKiB, "% .1f / % .1f")),
		mpb.AppendDecorators(decor.Percentage()),
	)

	// Create proxy reader
	proxyReader := bar.ProxyReader(srcReader)
	defer proxyReader.Close()

	// Copy file
	if written, err := io.Copy(dst, proxyReader); err != nil {
		fmt.Printf("failed to written %v %v", dst, written)
		return err
	}
	p.Wait()
	return nil
}

func main() {
	var from, to string
	var offset, limit int

	// Flags
	flag.StringVar(&from, "from", "", "/path/to/source")
	flag.StringVar(&to, "to", "", "/path/to/dest")
	flag.IntVar(&offset, "offset", 0, "offset")
	flag.IntVar(&limit, "limit", 0, "limit")
	flag.Parse()

	// gocopy -from=src.txt -to=dst.txt -offset=8 -limit=128
	if err := goCopy(from, to, offset, limit); err != nil {
		fmt.Println(err)
		return
	}
}
