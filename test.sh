#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./crude-lang-go "$input" > tmp.s
  gcc -static -o tmp tmp.s
  ./tmp
  actual="$?"
  msg="OK"
  code=0

  if [ "$actual" = "$expected" ]; then
    echo -e "\e[1;32mPASSED\e[0m : $input => $actual"
  else
    echo -e "\e[1;31mFAILED\e[0m : $input => $expected expected, but got $actual"
    msg="NG"
    code=-1
  fi
}

go build .

#############
# Assertion
#############
assert 3 '3'
assert 14 '14'
assert 3 '1+2'
assert 4 '5-1'
assert 6 '2+5-1'
assert 6 '1+2+3'
assert 41 " 12 + 34 - 5 "

rm -f tmp.s

echo $msg
exit $code
