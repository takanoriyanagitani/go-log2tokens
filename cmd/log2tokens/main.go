package main

import (
	"bufio"
	"context"
	"io"
	"iter"
	"log"
	"os"

	lt "github.com/takanoriyanagitani/go-log2tokens"
	tw "github.com/takanoriyanagitani/go-log2tokens/token/writer"
	wl "github.com/takanoriyanagitani/go-log2tokens/token/writer/least"
	. "github.com/takanoriyanagitani/go-log2tokens/util"
)

var rawScanOptions IO[[]string] = Of(os.Args[1:])

var scanOpts IO[lt.ScanOptions] = Bind(
	rawScanOptions,
	Lift(func(raw []string) (lt.ScanOptions, error) {
		var opts lt.RawScanOptions = raw
		return opts.ToScanOptions(
			func(rejected string) error {
				log.Printf("unknown scan option: %s\n", rejected)
				return nil
			},
		)
	}),
)

var stdin2tokens IO[iter.Seq[lt.Token]] = Bind(
	scanOpts,
	Lift(func(opt lt.ScanOptions) (iter.Seq[lt.Token], error) {
		return opt.ReaderToTokens(bufio.NewReader(os.Stdin)), nil
	}),
)

func stdin2writer(w io.Writer) IO[Void] {
	return Bind(
		stdin2tokens,
		func(tokens iter.Seq[lt.Token]) IO[Void] {
			var bw *bufio.Writer = bufio.NewWriter(w)
			defer bw.Flush()

			var rtw wl.RawTokenWriter = wl.WriterToRawTokenWriter(bw)
			var tw tw.TokenWriter = rtw.ToTokenWriter()

			return Bind(
				tw.WriteAll(tokens),
				Lift(func(_ Void) (Void, error) {
					bw.Flush()
					return Empty, nil
				}),
			)
		},
	)
}

var stdin2stdout IO[Void] = stdin2writer(os.Stdout)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
