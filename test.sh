#!/bin/bash
cat <<EOF | gcc -xc -c -o tmp2.o -
int ret10() { return 10;  }
int ret100() { return 100;  }
int add(int x, int y) { return x+y;  }
int sub(int x, int y) { return x-y;  }

int add6(int a, int b, int c, int d, int e, int f) {
  return a+b+c+d+e+f;
}
EOF

assert() {
    expected="$1"
    input="$2"

    ./crude-lang-go "$input" > tmp.s
    gcc -static -o tmp tmp.s tmp2.o
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
#nasm -f elf64 tmp.s
#ld -s -o tmp tmp.o

############
# Assertion
#############
assert 3 'main(){return 3;}'
assert 14 'main(){return 14;}'
assert 3 'main(){return 1+2;}'
assert 4 'main(){return 5-1;}'
assert 6 'main(){return 2+5-1;}'
assert 6 'main(){return 1+2+3;}'
assert 41 'main(){return  12 + 34 - 5;}'
assert 6 'main(){return 2*3;}'
assert 7 'main(){return 2*3+1;}'
assert 8 'main(){return (3-1)*4;}'
assert 4 'main(){return 8/2;}'
assert 4 'main(){return (3+5)/2;}'
assert 8 'main(){return  24 - 20 + (6- 4)*2;}'
assert 8 'main(){return   24 -20 + ( 6- 4)*2 ;}'
# Skipped because Go flag.Parse Cant recieve 'main(){-' although it works.
# assert -3 'main(){-3'
assert 1 'main(){return 1==1;}'
assert 0 'main(){return 2==1;}'
assert 0 'main(){return 10==1;}'
assert 0 'main(){return 1==10;}'
assert 1 'main(){return 10==10;}'
assert 1 'main(){return (14-2)==(4*3);}'
assert 0 'main(){return 1!=1;}'
assert 1 'main(){return 2!=1;}'
assert 1 'main(){return 2!=10;}'
assert 1 'main(){return 20!=1;}'
assert 1 'main(){return 20!=10;}'
assert 1 'main(){return (3-1)!=(2*2);}'
assert 1 'main(){return 1=<2;}'
assert 1 'main(){return 2=<2;}'
assert 0 'main(){return 3=<2;}'
assert 0 'main(){return 10=<2;}'
assert 1 'main(){return 1=<20;}'
assert 1 'main(){return 10=<20;}'
assert 1 'main(){return 10=<20;}'
assert 1 'main(){return (1-1)=<20;}'
assert 0 'main(){return (2-1)=<(3-3);}'
assert 1 'main(){return 1<2;}'
assert 0 'main(){return 2<2;}'
assert 0 'main(){return 3<2;}'
assert 0 'main(){return 30<2;}'
assert 1 'main(){return 3<20;}'
assert 0 'main(){return 30<20;}'
assert 0 'main(){return (5*13)<(2*2);}'
assert 1 'main(){return 2>1;}'
assert 0 'main(){return 2>3;}'
assert 1 'main(){return 12>3;}'
assert 0 'main(){return 2>13;}'
assert 0 'main(){return 12>13;}'
assert 0 'main(){return (2-1)>13;}'
assert 1 'main(){return 1>(12-12);}'
assert 1 'main(){return 2*2>(12-12);}'
assert 1 'main(){return (3-2)>(12-12);}'
assert 1 'main(){return 2>=1;}'
assert 0 'main(){return 2>=3;}'
assert 1 'main(){return 2>=2;}'
assert 1 'main(){return 12>=2;}'
assert 0 'main(){return 2>=12;}'
assert 0 'main(){return 1>=12;}'
assert 1 'main(){return (4-1)>=(5-2);}'
assert 121 'main(){return 121;144;169;}'
assert 3 'main(){1+1;return 6/2;6-2;}'
assert 4 'main(){1+1;6/2;return 6-2;}'
assert 4 'main(){ return 6-2;}'
assert 100 'main(){a=100;return a;}'
assert 10 'main(){a=2;b=8;return a+b;}'
assert 4 'main(){a=2;b=4;return b;}'
assert 9 'main(){a=2;b=4;return ((a+b)/a)*((a+b)/a);}'
assert 7 'main(){a=2;b=4;return (a/a+a)+b;}'
assert 2 'main(){hoge=2;return hoge;}'
assert 5 'main(){hoge_1=2;fuga=3;return hoge_1+fuga;}'
assert 4 'main(){ al = 2; be = 3; ga=1;return al*be-al*ga;}'
assert 6 'main(){al=2;be=3;ga=1;de=5;return (de-al)*ga+be;}'
assert 10 'main(){if (1) return 10;return 20;}'
assert 20 'main(){if (0) return 10;return 20;}'
assert 20 'main(){if (1==0) return 10;return 20;}'
assert 10 'main(){hoge=1;if (hoge) return 10;return 20;}'
assert 30 'main(){hoge=2;if (hoge==0) return 10; if(hoge==1) return 20; else return 30;}'
assert 20 'main(){hoge=2;if (hoge<1)return 10 ; else return 20;}'
assert 55 'main(){i=0; j=0; for (i=0; i=<10; i=i+1) j=i+j; return j;}'
assert 3 'main(){for (;;) return 3; return 5;}'
assert 10 'main(){i=0; for(i=0; ; i=i+1) if(i==10) return i;}'
assert 55 'main(){i=0; j=0; for (i=0; i=<10; ) {j=i+j;  i=i+1;} return j;}'
assert 10 'main(){return ret10();}'
assert 100 'main(){return ret100();}'
assert 8 'main(){return add(3, 5);}'
assert 2 'main(){return sub(5, 3);}'
assert 21 'main(){return add6(1,2,3,4,5,6);}'
assert 32 'main() { return ret32();  } ret32() { return 32;  }'
assert 9 'main() { return myadd(5,4); } myadd(x,y) { return x+y; }'
assert 10 'main() { return myadd2(4,6); } myadd2(x,y) { return x+y; }'
assert 11 'main() { return myadd3(5,6); } myadd3(x,y) { return x+y; }'
assert 7 'main() { return myadd4(1,6); } myadd4(x,y) { return x+y; }'
assert 17 'main() { return myadd5(11,6); } myadd5(x,y) { return x+y; }'
assert 89 'main(){return fib(11);} fib(n){if(n==1)return 1;if(n==2) return 1; return fib(n-1)+fib(n-2);}'
assert 3 'main() { x=3; return *&x;  }'
assert 4 'main() { x=3; y=&x; x=4; return *y;  }'
assert 3 'main() { x=3; y=&x; z=&y; return **z;  }'
assert 5 'main() { x=3; y=5; return *(&x+8);  }'
assert 3 'main() { x=3; y=5; return *(&y-8);  }'
assert 5 'main() { x=3; y=&x; *y=5; return x;  }'
assert 7 'main() { x=3; y=5; *(&x+8)=7; return y;  }'
assert 7 'main() { x=3; y=5; *(&y-8)=7; return x;  }'

rm -f tmp.s

echo $msg
exit $code
