package tests

import (
	"io"
	"testing"
)

func Closer(t *testing.T) func(io.Closer) {
	return func(c io.Closer) {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
