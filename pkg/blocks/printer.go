package blocks

import (
	"fmt"
	"io"
)

type Printer struct {
	writer io.Writer

	Message <-chan string
}

func NewPrinter(writer io.Writer) *Printer {
	return &Printer{writer: writer}
}

func (c *Printer) Process() {
	for msg := range c.Message {
		fmt.Fprintln(c.writer, msg)
	}
}
