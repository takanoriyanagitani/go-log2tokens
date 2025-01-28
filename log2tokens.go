package log2tokens

import (
	"bufio"
	"io"
	"iter"
	"os"
	"text/scanner"
)

type ScanOption string

const (
	ScanOptionIdents     ScanOption = "ScanIdents"
	ScanOptionInts       ScanOption = "ScanInts"
	ScanOptionFloats     ScanOption = "ScanFloats"
	ScanOptionChars      ScanOption = "ScanChars"
	ScanOptionStrings    ScanOption = "ScanStrings"
	ScanOptionRawStrings ScanOption = "ScanRawStrings"
	ScanOptionComments   ScanOption = "ScanComments"

	ScanOptionSkipComments ScanOption = "SkipComments"
)

var ScanOptionMap map[ScanOption]uint = map[ScanOption]uint{
	ScanOptionIdents:       scanner.ScanIdents,
	ScanOptionInts:         scanner.ScanInts,
	ScanOptionFloats:       scanner.ScanFloats,
	ScanOptionChars:        scanner.ScanChars,
	ScanOptionStrings:      scanner.ScanStrings,
	ScanOptionRawStrings:   scanner.ScanRawStrings,
	ScanOptionComments:     scanner.ScanComments,
	ScanOptionSkipComments: scanner.SkipComments,
}

func (o ScanOption) ToModeBit() uint {
	return ScanOptionMap[o]
}

type ScanOptions []ScanOption

func (o ScanOptions) ToConfigFunc() func(*scanner.Scanner) {
	return func(s *scanner.Scanner) {
		s.Mode = 0

		for _, opt := range o {
			var bit uint = opt.ToModeBit()
			s.Mode |= bit
		}
	}
}

func (o ScanOptions) ReaderToTokens(rdr io.Reader) iter.Seq[Token] {
	return ReaderToTokens(
		rdr,
		o.ToConfigFunc(),
	)
}

type RawScanOptions []string

func (r RawScanOptions) ToScanOptions(
	onInvalid func(string) error,
) (ScanOptions, error) {
	ret := make([]ScanOption, 0, len(r))
	for _, s := range r {
		_, found := ScanOptionMap[ScanOption(s)]
		if !found {
			e := onInvalid(s)
			if nil != e {
				return nil, e
			}
			continue
		}

		ret = append(ret, ScanOption(s))
	}
	return ret, nil
}

type Token struct {
	scanner.Position
	Text string
}

func ScannerToTokens(s *scanner.Scanner) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var buf Token
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
			buf.Position = s.Position
			buf.Text = s.TokenText()

			if !yield(buf) {
				return
			}
		}
	}
}

func ReaderToTokens(
	rdr io.Reader,
	config func(*scanner.Scanner),
) iter.Seq[Token] {
	var s scanner.Scanner

	s.Init(rdr)

	config(&s)

	return ScannerToTokens(&s)
}

type Config struct {
	Filename string
}

func (c Config) ToConfigFunc() func(*scanner.Scanner) {
	return func(s *scanner.Scanner) {
		s.Position.Filename = c.Filename
	}
}

var ConfigDefault Config

func StdinToTokens(config func(*scanner.Scanner)) iter.Seq[Token] {
	var br io.Reader = bufio.NewReader(os.Stdin)
	return ReaderToTokens(br, config)
}
