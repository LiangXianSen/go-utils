package common

import (
	"bytes"
	"io"
	"os"
)

func CaptureStdout(f func()) []byte {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	os.Stdout = writer
	os.Stderr = writer

	out := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, reader); err != nil {
			out <- []byte{}
		}
		out <- buf.Bytes()
	}()

	f()
	writer.Close()

	return <-out
}
