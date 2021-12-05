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

msg="OK"
code=0

assert() {
    expected="$1"
    input="$2"

    ./rouphc -i "$input" > tmp.s
    gcc -static -o tmp tmp.s tmp2.o
    ./tmp
    actual="$?"

    if [ "$actual" = "$expected" ]; then
        echo -e "\e[1;32mPASSED\e[0m : '$input' => $actual"
    else
        echo -e "\e[1;31mFAILED\e[0m : '$input' => $expected expected, but got $actual"
        msg="NG"
        code=-1
    fi
}

assert_input() {
    expected=$(eval $1)
    input="$2"

    ./rouphc -i "$input" > tmp.s
    gcc -static -o tmp tmp.s
    ./tmp > res
    actual=$(cat res)

    if [ "$actual" = "$expected" ]; then
        echo -e "\e[1;32mPASSED\e[0m : '$input' : '$actual'"
    else
        echo -e "\e[1;31mFAILED\e[0m : '$input' : '$expected' => '$actual'"
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
#assert_input 'cat expected/imms' 'func main():int{put("abc");}'
#assert_input 'cat expected/imms-vnum' 'func main():int{put("abc");let a:int=1;put(a);}'
#assert_input 'cat expected/vnum-imms' 'func main():int{let a:int=1;put(a);put("abc");}'
#assert_input 'cat expected/imms-vnum-imms' 'func main():int{put("abc");let a:int=1;put(a);put("abc");}'
#assert_input 'cat expected/vnum' 'func main():int{let a:int=1;put(a);}'
#assert_input 'cat expected/vs' 'func main():int{let a:char*="abc";put(a);}'
assert 3 'func main():int{return 3;}'
assert 14 'func main():int{return 14;}'
assert 3 'func main():int{return 1+2;}'
assert 4 'func main():int{return 5-1;}'
assert 6 'func main():int{return 2+5-1;}'
assert 6 'func main():int{return 1+2+3;}'
assert 41 'func main():int{return  12 + 34 - 5;}'
assert 6 'func main():int{return 2*3;}'
assert 7 'func main():int{return 2*3+1;}'
assert 8 'func main():int{return (3-1)*4;}'
assert 4 'func main():int{return 8/2;}'
assert 4 'func main():int{return (3+5)/2;}'
assert 8 'func main():int{return  24 - 20 + (6- 4)*2;}'
assert 8 'func main():int{return   24 -20 + ( 6- 4)*2 ;}'
# Skipped because Go flag.Parse Cant recieve 'int main(){-' although it works.
# assert -3 'func main():int{-3'

echo -e "\e[1;34m< relational/equality >\e[0m"
assert 1 'func main():int{return 1==1;}'
assert 0 'func main():int{return 2==1;}'
assert 0 'func main():int{return 10==1;}'
assert 0 'func main():int{return 1==10;}'
assert 1 'func main():int{return 10==10;}'
assert 1 'func main():int{return (14-2)==(4*3);}'
assert 0 'func main():int{return 1!=1;}'
assert 1 'func main():int{return 2!=1;}'
assert 1 'func main():int{return 2!=10;}'
assert 1 'func main():int{return 20!=1;}'
assert 1 'func main():int{return 20!=10;}'
assert 1 'func main():int{return (3-1)!=(2*2);}'
assert 1 'func main():int{return 1=<2;}'
assert 1 'func main():int{return 2=<2;}'
assert 0 'func main():int{return 3=<2;}'
assert 0 'func main():int{return 10=<2;}'
assert 1 'func main():int{return 1=<20;}'
assert 1 'func main():int{return 10=<20;}'
assert 1 'func main():int{return 10=<20;}'
assert 1 'func main():int{return (1-1)=<20;}'
assert 0 'func main():int{return (2-1)=<(3-3);}'
assert 1 'func main():int{return 1<2;}'
assert 0 'func main():int{return 2<2;}'
assert 0 'func main():int{return 3<2;}'
assert 0 'func main():int{return 30<2;}'
assert 1 'func main():int{return 3<20;}'
assert 0 'func main():int{return 30<20;}'
assert 0 'func main():int{return (5*13)<(2*2);}'
assert 1 'func main():int{return 2>1;}'
assert 0 'func main():int{return 2>3;}'
assert 1 'func main():int{return 12>3;}'
assert 0 'func main():int{return 2>13;}'
assert 0 'func main():int{return 12>13;}'
assert 0 'func main():int{return (2-1)>13;}'
assert 1 'func main():int{return 1>(12-12);}'
assert 1 'func main():int{return 2*2>(12-12);}'
assert 1 'func main():int{return (3-2)>(12-12);}'
assert 1 'func main():int{return 2>=1;}'
assert 0 'func main():int{return 2>=3;}'
assert 1 'func main():int{return 2>=2;}'
assert 1 'func main():int{return 12>=2;}'
assert 0 'func main():int{return 2>=12;}'
assert 0 'func main():int{return 1>=12;}'
assert 1 'func main():int{return (4-1)>=(5-2);}'
assert 121 'func main():int{return 121;144;169;}'
assert 3 'func main():int{1+1;return 6/2;6-2;}'
assert 4 'func main():int{1+1;6/2;return 6-2;}'
assert 4 'func main():int{ return 6-2;}'

echo -e "\e[1;34m< identifier >\e[0m"
assert 100 'func main():int{let a:int=100;return a;}'
assert 10 'func main():int{let a:int=2;let b:int=8;return a+b;}'
assert 4 'func main():int{let a:int=2;let b:int=4;return b;}'
assert 9 'func main():int{let a:int=2;let b:int=4;return ((a+b)/a)*((a+b)/a);}'
assert 7 'func main():int{let a:int=2;let b:int=4;return (a/a+a)+b;}'
assert 2 'func main():int{let hoge:int=2;return hoge;}'
assert 5 'func main():int{let hoge_1:int=2;let fuga:int=3;return hoge_1+fuga;}'
assert 4 'func main():int{ let al :int= 2; let be :int= 3; let ga:int=1;return al*be-al*ga;}'
assert 6 'func main():int{let al:int=2;let be:int=3;let ga:int=1;let de:int=5;return (de-al)*ga+be;}'

echo -e "\e[1;34m< if >\e[0m"
assert 10 'func main():int{if 1 {return 10;return 20;}}'
assert 20 'func main():int{if 0{ return 10;}return 20;}'
assert 20 'func main():int{if 1==0 {return 10;}return 20;}'
assert 10 'func main():int{let hoge:int=1;if hoge {return 10;}return 20;}'
assert 30 'func main():int{let hoge:int=2;if hoge==0 {return 10;} else {return 30;}}'
assert 20 'func main():int{let hoge:int=2;if hoge<1 {return 10;}  else {return 20;}}'
assert 10 'func main():int{let hoge:int=1;if hoge == 1 {return 10 ;} else if hoge == 2 {return 20;} else {return 30;}}'
assert 20 'func main():int{let hoge:int=2;if hoge == 1 {return 10 ;} else if hoge == 2 {return 20;} else {return 30;}}'
assert 30 'func main():int{let hoge:int=3;if hoge == 1 {return 10 ;} else if hoge == 2 {return 20;} else {return 30;}}'

echo -e "\e[1;34m< for >\e[0m"
assert 55 'func main():int{let i:int=0; let j:int=0; for i=0; i=<10; i=i+1{ j=i+j;} return j;}'
assert 3 'func main():int{for ;; {return 3;} return 5;}'
assert 10 'func main():int{let i:int=0; for i=0;; i=i+1 { if(i==10) return i;}}'
assert 55 'func main():int{let i:int=0; let j:int=0; for i=0; i=<10; {j=i+j;  i=i+1;} return j;}'

echo -e "\e[1;34m< function call >\e[0m"
assert 10 'func main():int{return ret10();}'
assert 100 'func main():int{return ret100();}'
assert 8 'func main():int{return add(3, 5);}'
assert 2 'func main():int{return sub(5, 3);}'
assert 21 'func main():int{return add6(1,2,3,4,5,6);}'
assert 32 'func main():int { return ret32();  } func ret32():int { return 32;  }'
assert 9 'func main():int { return myadd(5,4); } func myadd(x:int,y:int):int { return x+y; }'
assert 9 'func myadd(x:int,y:int):int { return x+y; } func main():int{ return myadd(5,4); }'
assert 10 'func main():int { return myadd2(4,6); } func myadd2(x:int,y:int):int { return x+y; }'
assert 11 'func main():int { return myadd3(5,6); } func myadd3(x:int,y:int):int { return x+y; }'
assert 7 'func main():int { return myadd4(1,6); } func myadd4(x:int,y:int):int { return x+y; }'
assert 17 'func main():int { return myadd5(11,6); } func myadd5(x:int,y:int):int { return x+y; }'
assert 89 'func main():int{return fib(11);} func fib(n:int):int{if n==1 {return 1;} if n==2 {return 1;} return fib(n-1)+fib(n-2);}'

echo -e "\e[1;34m< pointer >\e[0m"
assert 3 'func main():int { let x:int=3; return *&x;  }'
assert 3 'func main():int { let x:int=3; let y:int*=&x; let z:int**=&y; return **z;  }'
assert 5 'func main():int { let x:int=3; let y:int=5; return *(&x+1);  }'
assert 5 'func main():int { let x:int=3; let y:int=5; return *(1+&x);  }'
assert 3 'func main():int { let x:int=3; let y:int=5; return *(&y-1);  }'
assert 5 'func main():int { let x:int=3; let y:int=5; let z:int*=&x; return *(z+1);  }'
assert 3 'func main():int { let x:int=3; let y:int=5; let z:int*=&y; return *(z-1);  }'
assert 5 'func main():int { let x:int=3; let y:int*=&x; *y=5; return x;  }'
assert 7 'func main():int { let x:int=3; let y:int=5; *(&x+1)=7; return y;  }'
assert 7 'func main():int { let x:int=3; let y:int=5; *(&y-1)=7; return x;  }'
assert 8 'func main():int { let x:int=3; let y:int=5; return foo(&x, y);  } func foo(x:int*, y:int):int { return *x + y;  }'

echo -e "\e[1;34m< array >\e[0m"
assert 2 'func main():int { let x:int[1]; *x=2; return *x;  }'
assert 3 'func main():int { let x:int[2]; let y:int*=&x; *y=3; return *x;  }'
assert 3 'func main():int { let x:int[3]; *x=3; *(x+1)=4; *(x+2)=5; return *x;  }'
assert 4 'func main():int { let x:int[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+1);  }'
assert 5 'func main():int { let x:int[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+2);  }'

assert 5 'func main():int { let x:int[2][3]; let y:int*=x; *y=5; return **x;  }'
assert 1 'func main():int { let x:int[2][3]; let y:int*=x; *(y+1)=1; return *(*x+1);  }'
assert 2 'func main():int { let x:int[2][3]; let y:int*=x; *(y+2)=2; return *(*x+2);  }'
assert 3 'func main():int { let x:int[2][3]; let y:int*=x; *(y+3)=3; return **(x+1);  }'
assert 4 'func main():int { let x:int[2][3]; let y:int*=x; *(y+4)=4; return *(*(x+1)+1);  }'
assert 5 'func main():int { let x:int[2][3]; let y:int*=x; *(y+5)=5; return *(*(x+1)+2);  }'
assert 6 'func main():int { let x:int[2][3]; let y:int*=x; *(y+6)=6; return **(x+2);  }'
assert 3 'func main():int { let x:int[3]; *x=3; x[1]=4; x[2]=5; return *x;  }'
assert 4 'func main():int { let x:int[3]; *x=3; x[1]=4; x[2]=5; return *(x+1);  }'
assert 5 'func main():int { let x:int[3]; *x=3; x[1]=4; x[2]=5; return *(x+2);  }'
assert 5 'func main():int { let x:int[3]; *x=3; x[1]=4; x[2]=5; return *(x+2);  }'
assert 5 'func main():int { let x:int[3]; *x=3; x[1]=4; 2[x]=5; return *(x+2);  }'

assert 0 'func main():int { let x:int[2][3]; let y:int*=x; y[0]=0; return x[0][0];  }'
assert 1 'func main():int { let x:int[2][3]; let y:int*=x; y[1]=1; return x[0][1];  }'
assert 2 'func main():int { let x:int[2][3]; let y:int*=x; y[2]=2; return x[0][2];  }'
assert 3 'func main():int { let x:int[2][3]; let y:int*=x; y[3]=3; return x[1][0];  }'
assert 4 'func main():int { let x:int[2][3]; let y:int*=x; y[4]=4; return x[1][1];  }'
assert 5 'func main():int { let x:int[2][3]; let y:int*=x; y[5]=5; return x[1][2];  }'
assert 6 'func main():int { let x:int[2][3]; let y:int*=x; y[6]=6; return x[2][0];  }'
assert 2 'func main():int { let x:int[3]; x[0]=2; x[1]=4; x[2]=5; return *x;  }'
assert 5 'func main():int { let x:int[5]; let i:int=0; let j:int=0;for i=0; i<5; i=i+1 {x[i]=0;} x[4]=5;j=x[4]; return j; }'

echo -e "\e[1;34m< sizeof >\e[0m"
assert 8 'func main():int { let x:int; return sizeof(x);  }'
assert 8 'func main():int { let x:int; return sizeof x;  }'
assert 8 'func main():int { let x:int*; return sizeof(x);  }'
assert 32 'func main():int { let x:int[4]; return sizeof(x);  }'
assert 96 'func main():int { let x:int[3][4]; return sizeof(x);  }'
assert 32 'func main():int { let x:int[3][4]; return sizeof(*x);  }'
assert 8 'func main():int { let x:int[3][4]; return sizeof(**x);  }'
assert 9 'func main():int { let x:int[3][4]; return sizeof(**x) + 1;  }'
assert 9 'func main():int { let x:int[3][4]; return sizeof **x + 1;  }'
assert 8 'func main():int { let x:int[3][4]; return sizeof(**x + 1);  }'

echo -e "\e[1;34m< global variables >\e[0m"
assert 0 'let x:int; func main():int { return x;  }'
assert 3 'let x:int; func main():int { x=3; return x;  }'
assert 7 'let x:int;let y:int; func main():int { x=3;y=4; return x+y;  }'
assert 7 'let x:int;let y:int; func main():int { x=3;y=4; let z:int = x+y; return z;  }'
assert 0 'let x:int[4]; func main():int { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[0];  }'
assert 1 'let x:int[4]; func main():int { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[1];  }'
assert 2 'let x:int[4]; func main():int { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[2];  }'
assert 3 'let x:int[4]; func main():int { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[3];  }'
assert 8 'let x:int; func main():int { return sizeof(x);  }'
assert 32 'let x:int[4]; func main():int { return sizeof(x);  }'


echo -e "\e[1;34m< char type >\e[0m"
assert 1 'func main():int { let x:char=1; return x;  }'
assert 1 'func main():int { let x:char=1; let y:char=2; return x;  }'
assert 2 'func main():int { let x:char=1; let y:char=2; return y;  }'

assert 1 'func main():int { let x:char; return sizeof(x);  }'
assert 10 'func main():int { let x:char[10]; return sizeof(x);  }'
assert 1 'func main():int { return sub_char(7, 3, 3);  } func sub_char(a:char, b:char, c:char):int { return a-b-c;  }'
assert 2 'func main():int { let c:char = f(); return c;  } func f():char*{ return 2; }'
echo -e "\e[1;34m< string type >\e[0m"
assert 97 'func main():int { return "abc"[0];  }'
assert 97 'func main():int { let s:char* = f(); return s[0];  } func f():char*{ return "abc"; }'
assert 98 'func main():int { return "abc"[1];  }'
assert 99 'func main():int { return "abc"[2];  }'
assert 0 'func main():int { return "abc"[3];  }'
assert 97 'func main():int { let s:char* = "abc"; return s[0]; }'
assert 4 'func main():int { return sizeof("abc");  }'

echo -e "\e[1;34m< escape sequence >\e[0m"
assert 7 'func main():int { return "\a"[0];  }'
assert 8 'func main():int { return "\b"[0];  }'
assert 9 'func main():int { return "\t"[0];  }'
assert 10 'func main():int { return "\n"[0];  }'
assert 11 'func main():int { return "\v"[0];  }'
assert 12 'func main():int { return "\f"[0];  }'
assert 13 'func main():int { return "\r"[0];  }'
assert 27 'func main():int { return "\e"[0];  }'
assert 0 'func main():int { return "\0"[0];  }'
assert 106 'func main():int { return "\j"[0];  }'
assert 107 'func main():int { return "\k"[0];  }'
assert 108 'func main():int { return "\l"[0];  }'


echo -e "\e[1;34m< remainder >\e[0m"
assert 1 'func main():int { return 3%2;  }'
assert 0 'func main():int { return 4%2;  }'
assert 4 'func main():int { return 4%5;  }'
assert 4 'func main():int { let a:int = 60; let b:int = 7; let c:int = a%b; return c;  }'
assert 1 'func main():int { return (7%4) == 3;  }'
assert 0 'func main():int { return (9%4) == 3;  }'




rm -f tmp.s
rm -f rouphc

echo $msg
exit $code
