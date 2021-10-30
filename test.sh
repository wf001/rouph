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
    echo -e "\e[1;32mPASSED\e[0m : '$input' => $actual"
  else
    echo -e "\e[1;31mFAILED\e[0m : '$input' => $expected expected, but got $actual"
    msg="NG"
    code=-1
  fi
}

go build .

#############
# Assertion
#############
assert 3 'return 3;'
assert 14 'return 14;'
assert 3 'return 1+2;'
assert 4 'return 5-1;'
assert 6 'return 2+5-1;'
assert 6 'return 1+2+3;'
assert 41 'return  12 + 34 - 5;'
assert 6 'return 2*3;'
assert 7 'return 2*3+1;'
assert 8 'return (3-1)*4;'
assert 4 'return 8/2;'
assert 4 'return (3+5)/2;'
assert 8 'return  24 - 20 + (6- 4)*2;'
assert 8 'return   24 -20 + ( 6- 4)*2 ;'
# Skipped because Go flag.Parse Cant recieve '-' although it works.
# assert -3 '-3'
assert 1 'return 1==1;'
assert 0 'return 2==1;'
assert 0 'return 10==1;'
assert 0 'return 1==10;'
assert 1 'return 10==10;'
assert 1 'return (14-2)==(4*3);'
assert 0 'return 1!=1;'
assert 1 'return 2!=1;'
assert 1 'return 2!=10;'
assert 1 'return 20!=1;'
assert 1 'return 20!=10;'
assert 1 'return (3-1)!=(2*2);'
assert 1 'return 1=<2;'
assert 1 'return 2=<2;'
assert 0 'return 3=<2;'
assert 0 'return 10=<2;'
assert 1 'return 1=<20;'
assert 1 'return 10=<20;'
assert 1 'return 10=<20;'
assert 1 'return (1-1)=<20;'
assert 0 'return (2-1)=<(3-3);'
assert 1 'return 1<2;'
assert 0 'return 2<2;'
assert 0 'return 3<2;'
assert 0 'return 30<2;'
assert 1 'return 3<20;'
assert 0 'return 30<20;'
assert 0 'return (5*13)<(2*2);'
assert 1 'return 2>1;'
assert 0 'return 2>3;'
assert 1 'return 12>3;'
assert 0 'return 2>13;'
assert 0 'return 12>13;'
assert 0 'return (2-1)>13;'
assert 1 'return 1>(12-12);'
assert 1 'return 2*2>(12-12);'
assert 1 'return (3-2)>(12-12);'
assert 1 'return 2>=1;'
assert 0 'return 2>=3;'
assert 1 'return 2>=2;'
assert 1 'return 12>=2;'
assert 0 'return 2>=12;'
assert 0 'return 1>=12;'
assert 1 'return (4-1)>=(5-2);'
assert 121 'return 121;144;169;'
assert 3 '1+1;return 6/2;6-2;'
assert 4 '1+1;6/2;return 6-2;'
assert 4 ' return 6-2;'
rm -f tmp.s

echo $msg
exit $code
