.intel_syntax noprefix
.data
.L.data.0:
  .byte 97
  .byte 103
  .byte 99
  .byte 0
.text
.global main
main:
  push rbp
  mov rbp, rsp
  sub rsp, 0
  push offset .L.data.0
  mov r8, 3
  mov rsi, [rsp]
  call .puts
  add rsp, 8
  push 0
  pop rax
  jmp .Lreturn.main
.Lreturn.main:
  mov rsp, rbp
  pop rbp
  ret
.puts:
  mov rdx, 1
  mov rax, 1
  mov rdi, 1
  syscall
  sub r8, 1
  add rsi, 1
  cmp r8, 0
  jnz .puts
  push 0xa
  mov rsi, rsp
  mov rax, 1
  syscall
  pop rax
  mov rax, 1
  ret
