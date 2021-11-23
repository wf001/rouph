.intel_syntax noprefix
.data
.L.data.0:
  .byte 104
  .byte 111
  .byte 103
  .byte 101
  .byte 0
.text
.global main
main:
  push rbp
  mov rbp, rsp
  sub rsp, 0
  call .puts
.puts:
  push offset .L.data.0
  mov rsi, [rsp]
  mov rdx, 8
  mov rax, 1
  mov rdi, 1
  syscall
  pop rax
  push 0xa
  mov rsi, rsp
  mov rax, 1
  syscall
  pop rax
  ret
  add rsp, 8
  push 0
  pop rax
  jmp .Lreturn.main
.Lreturn.main:
  mov rsp, rbp
  pop rbp
  ret
