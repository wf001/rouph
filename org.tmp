.intel_syntax noprefix
.data
.L.data.0:
  .byte 72
  .byte 101
  .byte 108
  .byte 108
  .byte 111
  .byte 0
.text
.global main
main:
  push rbp
  mov rbp, rsp
  sub rsp, 0
  push offset .L.data.0
  push 0
  pop rdi
  pop rax
  imul rdi, 1
  add rax, rdi
  push rax
  pop rax
  movsx rax, byte ptr [rax]
  push rax
  pop rax
  call .puts
  jmp .Lreturn.main

.puts:
  push offset .L.data.0
  mov rsi, [rsp]
  mov rdx, 5
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

.Lreturn.main:
  mov rsp, rbp
  pop rbp
  ret
