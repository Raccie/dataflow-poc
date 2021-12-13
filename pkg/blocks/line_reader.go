package blocks

import "io"

type LineReader struct {
	Reader <-chan io.Reader
}
