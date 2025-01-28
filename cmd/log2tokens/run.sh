#!/bin/sh

opts=

opts="${opts} ScanIdents"
opts="${opts} ScanChars"
opts="${opts} ScanStrings"
opts="${opts} ScanRawStrings"
opts="${opts} ScanComments"

cat main.go |
	./log2tokens $opts |
	grep -e '..' |
	sort |
	uniq -c |
	sort -n
