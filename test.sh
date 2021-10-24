#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./crude "$input" > tmp.s
  gcc -static -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

go build crude.go

assert 3 '3'
assert 14 '14'

echo OK
