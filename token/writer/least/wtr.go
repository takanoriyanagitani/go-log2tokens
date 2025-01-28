package least

import (
	"context"
	"io"

	lt "github.com/takanoriyanagitani/go-log2tokens"
	tw "github.com/takanoriyanagitani/go-log2tokens/token/writer"
	. "github.com/takanoriyanagitani/go-log2tokens/util"
)

type RawTokenWriter func(string) IO[Void]

func (r RawTokenWriter) ToTokenWriter() tw.TokenWriter {
	return func(t lt.Token) IO[Void] {
		var raw string = t.Text
		return r(raw)
	}
}

func WriterToRawTokenWriter(w io.Writer) RawTokenWriter {
	return func(raw string) IO[Void] {
		return func(_ context.Context) (Void, error) {
			_, e := io.WriteString(w, raw+"\n")
			return Empty, e
		}
	}
}
