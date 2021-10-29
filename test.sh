#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./crude-lang-go "$input" > tmp.s
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

go build .

assert 3 '3'
assert 14 '14'
assert 3 '1+2'
assert 4 '5-1'
assert 6 '2+5-1'
assert 6 '1+2+3'
assert 41 " 12 + 34 - 5 "

echo OK
