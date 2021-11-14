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
assert 3 'int main(){return 3;}'
assert 14 'int main(){return 14;}'
assert 3 'int main(){return 1+2;}'
assert 4 'int main(){return 5-1;}'
assert 6 'int main(){return 2+5-1;}'
assert 6 'int main(){return 1+2+3;}'
assert 41 'int main(){return  12 + 34 - 5;}'
assert 6 'int main(){return 2*3;}'
assert 7 'int main(){return 2*3+1;}'
assert 8 'int main(){return (3-1)*4;}'
assert 4 'int main(){return 8/2;}'
assert 4 'int main(){return (3+5)/2;}'
assert 8 'int main(){return  24 - 20 + (6- 4)*2;}'
assert 8 'int main(){return   24 -20 + ( 6- 4)*2 ;}'
# Skipped because Go flag.Parse Cant recieve 'int main(){-' although it works.
# assert -3 'int main(){-3'

# relational/equality
assert 1 'int main(){return 1==1;}'
assert 0 'int main(){return 2==1;}'
assert 0 'int main(){return 10==1;}'
assert 0 'int main(){return 1==10;}'
assert 1 'int main(){return 10==10;}'
assert 1 'int main(){return (14-2)==(4*3);}'
assert 0 'int main(){return 1!=1;}'
assert 1 'int main(){return 2!=1;}'
assert 1 'int main(){return 2!=10;}'
assert 1 'int main(){return 20!=1;}'
assert 1 'int main(){return 20!=10;}'
assert 1 'int main(){return (3-1)!=(2*2);}'
assert 1 'int main(){return 1=<2;}'
assert 1 'int main(){return 2=<2;}'
assert 0 'int main(){return 3=<2;}'
assert 0 'int main(){return 10=<2;}'
assert 1 'int main(){return 1=<20;}'
assert 1 'int main(){return 10=<20;}'
assert 1 'int main(){return 10=<20;}'
assert 1 'int main(){return (1-1)=<20;}'
assert 0 'int main(){return (2-1)=<(3-3);}'
assert 1 'int main(){return 1<2;}'
assert 0 'int main(){return 2<2;}'
assert 0 'int main(){return 3<2;}'
assert 0 'int main(){return 30<2;}'
assert 1 'int main(){return 3<20;}'
assert 0 'int main(){return 30<20;}'
assert 0 'int main(){return (5*13)<(2*2);}'
assert 1 'int main(){return 2>1;}'
assert 0 'int main(){return 2>3;}'
assert 1 'int main(){return 12>3;}'
assert 0 'int main(){return 2>13;}'
assert 0 'int main(){return 12>13;}'
assert 0 'int main(){return (2-1)>13;}'
assert 1 'int main(){return 1>(12-12);}'
assert 1 'int main(){return 2*2>(12-12);}'
assert 1 'int main(){return (3-2)>(12-12);}'
assert 1 'int main(){return 2>=1;}'
assert 0 'int main(){return 2>=3;}'
assert 1 'int main(){return 2>=2;}'
assert 1 'int main(){return 12>=2;}'
assert 0 'int main(){return 2>=12;}'
assert 0 'int main(){return 1>=12;}'
assert 1 'int main(){return (4-1)>=(5-2);}'
assert 121 'int main(){return 121;144;169;}'
assert 3 'int main(){1+1;return 6/2;6-2;}'
assert 4 'int main(){1+1;6/2;return 6-2;}'
assert 4 'int main(){ return 6-2;}'
# identifier
assert 100 'int main(){int a=100;return a;}'
assert 10 'int main(){int a=2;int b=8;return a+b;}'
assert 4 'int main(){int a=2;int b=4;return b;}'
assert 9 'int main(){int a=2;int b=4;return ((a+b)/a)*((a+b)/a);}'
assert 7 'int main(){int a=2;int b=4;return (a/a+a)+b;}'
assert 2 'int main(){int hoge=2;return hoge;}'
assert 5 'int main(){int hoge_1=2;int fuga=3;return hoge_1+fuga;}'
assert 4 'int main(){ int al = 2; int be = 3; int ga=1;return al*be-al*ga;}'
assert 6 'int main(){int al=2;int be=3;int ga=1;int de=5;return (de-al)*ga+be;}'
# if
assert 10 'int main(){if (1) return 10;return 20;}'
assert 20 'int main(){if (0) return 10;return 20;}'
assert 20 'int main(){if (1==0) return 10;return 20;}'
assert 10 'int main(){int hoge=1;if (hoge) return 10;return 20;}'
assert 30 'int main(){int hoge=2;if (hoge==0) return 10; if(hoge==1) return 20; else return 30;}'
assert 20 'int main(){int hoge=2;if (hoge<1)return 10 ; else return 20;}'
# for
assert 55 'int main(){int i=0; int j=0; for (i=0; i=<10; i=i+1) j=i+j; return j;}'
assert 3 'int main(){for (;;) return 3; return 5;}'
assert 10 'int main(){int i=0; for(i=0; ; i=i+1) if(i==10) return i;}'
assert 55 'int main(){int i=0; int j=0; for (i=0; i=<10; ) {j=i+j;  i=i+1;} return j;}'
#function
assert 10 'int main(){return ret10();}'
assert 100 'int main(){return ret100();}'
assert 8 'int main(){return add(3, 5);}'
assert 2 'int main(){return sub(5, 3);}'
assert 21 'int main(){return add6(1,2,3,4,5,6);}'
assert 32 'int main() { return ret32();  } int ret32() { return 32;  }'
assert 9 'int main() { return myadd(5,4); } int myadd(int x,int y) { return x+y; }'
assert 10 'int main() { return myadd2(4,6); } int myadd2(int x,int y) { return x+y; }'
assert 11 'int main() { return myadd3(5,6); } int myadd3(int x,int y) { return x+y; }'
assert 7 'int main() { return myadd4(1,6); } int myadd4(int x,int y) { return x+y; }'
assert 17 'int main() { return myadd5(11,6); } int myadd5(int x,int y) { return x+y; }'
assert 89 'int main(){return fib(11);} int fib(int n){if(n==1)return 1;if(n==2) return 1; return fib(n-1)+fib(n-2);}'

# poiter
assert 3 'int main() { int x=3; return *&x;  }'
assert 3 'int main() { int x=3; int *y=&x; int **z=&y; return **z;  }'
assert 5 'int main() { int x=3; int y=5; return *(&x+1);  }'
assert 5 'int main() { int x=3; int y=5; return *(1+&x);  }'
assert 3 'int main() { int x=3; int y=5; return *(&y-1);  }'
assert 5 'int main() { int x=3; int y=5; int *z=&x; return *(z+1);  }'
assert 3 'int main() { int x=3; int y=5; int *z=&y; return *(z-1);  }'
assert 5 'int main() { int x=3; int *y=&x; *y=5; return x;  }'
assert 7 'int main() { int x=3; int y=5; *(&x+1)=7; return y;  }'
assert 7 'int main() { int x=3; int y=5; *(&y-1)=7; return x;  }'
assert 8 'int main() { int x=3; int y=5; return foo(&x, y);  } int foo(int *x, int y) { return *x + y;  }'
# array
assert 2 'int main() { int x[1]; *x=2; return *x;  }'
assert 3 'int main() { int x[2]; int *y=&x; *y=3; return *x;  }'
assert 3 'int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *x;  }'
assert 4 'int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+1);  }'
assert 5 'int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+2);  }'

assert 5 'int main() { int x[2][3]; int *y=x; *y=5; return **x;  }'
assert 1 'int main() { int x[2][3]; int *y=x; *(y+1)=1; return *(*x+1);  }'
assert 2 'int main() { int x[2][3]; int *y=x; *(y+2)=2; return *(*x+2);  }'
assert 3 'int main() { int x[2][3]; int *y=x; *(y+3)=3; return **(x+1);  }'
assert 4 'int main() { int x[2][3]; int *y=x; *(y+4)=4; return *(*(x+1)+1);  }'
assert 5 'int main() { int x[2][3]; int *y=x; *(y+5)=5; return *(*(x+1)+2);  }'
assert 6 'int main() { int x[2][3]; int *y=x; *(y+6)=6; return **(x+2);  }'
assert 3 'int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *x;  }'
assert 4 'int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+1);  }'
assert 5 'int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+2);  }'
assert 5 'int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+2);  }'
assert 5 'int main() { int x[3]; *x=3; x[1]=4; 2[x]=5; return *(x+2);  }'

assert 0 'int main() { int x[2][3]; int *y=x; y[0]=0; return x[0][0];  }'
assert 1 'int main() { int x[2][3]; int *y=x; y[1]=1; return x[0][1];  }'
assert 2 'int main() { int x[2][3]; int *y=x; y[2]=2; return x[0][2];  }'
assert 3 'int main() { int x[2][3]; int *y=x; y[3]=3; return x[1][0];  }'
assert 4 'int main() { int x[2][3]; int *y=x; y[4]=4; return x[1][1];  }'
assert 5 'int main() { int x[2][3]; int *y=x; y[5]=5; return x[1][2];  }'
assert 6 'int main() { int x[2][3]; int *y=x; y[6]=6; return x[2][0];  }'
assert 2 'int main() { int x[3]; x[0]=2; x[1]=4; x[2]=5; return *x;  }'
assert 5 'int main() { int x[5]; int i=0; int j=0;for(i=0; i<5; i=i+1){x[i]=0;} x[4]=5;j=x[4]; return j; }'

# sizeof
assert 8 'int main() { int x; return sizeof(x);  }'
assert 8 'int main() { int x; return sizeof x;  }'
assert 8 'int main() { int *x; return sizeof(x);  }'
assert 32 'int main() { int x[4]; return sizeof(x);  }'
assert 96 'int main() { int x[3][4]; return sizeof(x);  }'
assert 32 'int main() { int x[3][4]; return sizeof(*x);  }'
assert 8 'int main() { int x[3][4]; return sizeof(**x);  }'
assert 9 'int main() { int x[3][4]; return sizeof(**x) + 1;  }'
assert 9 'int main() { int x[3][4]; return sizeof **x + 1;  }'
assert 8 'int main() { int x[3][4]; return sizeof(**x + 1);  }'
# global variables
assert 0 'int x; int main() { return x;  }'
assert 3 'int x; int main() { x=3; return x;  }'
assert 7 'int x;int y; int main() { x=3;y=4; return x+y;  }'
assert 7 'int x;int y; int main() { x=3;y=4; int z = x+y; return z;  }'
assert 0 'int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[0];  }'
assert 1 'int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[1];  }'
assert 2 'int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[2];  }'
assert 3 'int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[3];  }'
assert 8 'int x; int main() { return sizeof(x);  }'
assert 32 'int x[4]; int main() { return sizeof(x);  }'

rm -f tmp.s

echo $msg
exit $code
