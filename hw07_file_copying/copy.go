package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type ProgressBar struct {
	writer   io.Writer
	total    int64
	progress int64
}

func NewProgressBar(writer io.Writer, total int64) *ProgressBar {
	return &ProgressBar{
		writer: writer,
		total:  total,
	}
}

func (b *ProgressBar) Increment(inc int64) {
	b.progress += inc
	progressPer := 100 * b.progress / b.total
	progressDone := 30 * progressPer / 100
	progressUndone := 30 - progressDone
	progressString := fmt.Sprintf(
		"copy: %s>%s %d%%\n",
		strings.Repeat("=", int(progressDone)),
		strings.Repeat(" ", int(progressUndone)),
		progressPer,
	)
	b.writer.Write([]byte(progressString))
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromStat, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if !fromStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fromStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if limit == 0 || limit > fromStat.Size() {
		limit = fromStat.Size() - offset
	}

	return doCopy(fromFile, toFile, limit, offset)
}

func doCopy(from io.ReadSeeker, to io.WriteSeeker, limit, offset int64) error {
	var (
		totalCopied int64
		batchSize   int64 = 256
	)

	if batchSize > limit {
		batchSize = limit
	}

	_, err := from.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	progressBar := NewProgressBar(os.Stdout, limit)

	for {
		if totalCopied+batchSize > limit {
			batchSize = limit - totalCopied
		}

		copied, err := io.CopyN(to, from, batchSize)
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			}
			return err
		}

		progressBar.Increment(copied)

		totalCopied += copied
		if totalCopied >= limit {
			break
		}
	}

	return nil
}
