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
assert 41 ' 12 + 34 - 5 '
assert 6 '2*3 '
assert 7 '2*3+1'
assert 8 '(3-1)*4'
assert 4 '8/2'
assert 4 '(3+5)/2'
assert 8 ' 24 - 20 + (6- 4)*2 '
# Skipped because Go flag.Parse Cant recieve '-' although it works.
# assert -3 '-3'
assert 1 '1==1'
assert 0 '2==1'
assert 0 '10==1'
assert 0 '1==10'
assert 1 '10==10'
assert 1 '(14-2)==(4*3)'
assert 0 '1!=1'
assert 1 '2!=1'
assert 1 '2!=10'
assert 1 '20!=1'
assert 1 '20!=10'
assert 1 '(3-1)!=(2*2)'
assert 1 '1=<2'
assert 1 '2=<2'
assert 0 '3=<2'
assert 0 '10=<2'
assert 1 '1=<20'
assert 1 '10=<20'
assert 1 '10=<20'
assert 1 '(1-1)=<20'
assert 0 '(2-1)=<(3-3)'
assert 1 '1<2'
assert 0 '2<2'
assert 0 '3<2'
assert 0 '30<2'
assert 1 '3<20'
assert 0 '30<20'
assert 0 '(5*13)<(2*2)'
assert 1 '2>1'
assert 0 '2>3'
assert 1 '12>3'
assert 0 '2>13'
assert 0 '12>13'
assert 0 '(2-1)>13'
assert 1 '1>(12-12)'
assert 1 '2*2>(12-12)'
assert 1 '(3-2)>(12-12)'
assert 1 '2>=1'
assert 0 '2>=3'
assert 1 '2>=2'
assert 1 '12>=2'
assert 0 '2>=12'
assert 0 'a>=12'
assert 1 '(4-1)>=(5-2)'

rm -f tmp.s

echo $msg
exit $code
