package tokwtr

import (
	"context"
	"iter"

	lt "github.com/takanoriyanagitani/go-log2tokens"
	. "github.com/takanoriyanagitani/go-log2tokens/util"
)

type TokenWriter func(lt.Token) IO[Void]

func (w TokenWriter) WriteAll(tokens iter.Seq[lt.Token]) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		for tok := range tokens {
			select {
			case <-ctx.Done():
				return Empty, ctx.Err()
			default:
			}

			_, e := w(tok)(ctx)
			if nil != e {
				return Empty, e
			}
		}

		return Empty, nil
	}
}
