	.text
	.def	 @feat.00;
	.scl	3;
	.type	0;
	.endef
	.globl	@feat.00
.set @feat.00, 0
	.file	"test.ll"
	.def	 main;
	.scl	2;
	.type	32;
	.endef
	.globl	main                    # -- Begin function main
	.p2align	4, 0x90
main:                                   # @main
.seh_proc main
# %bb.0:                                # %main_body
	pushq	%rbp
	.seh_pushreg %rbp
	subq	$48, %rsp
	.seh_stackalloc 48
	leaq	48(%rsp), %rbp
	.seh_setframe %rbp, 48
	.seh_endprologue
	callq	__main
	movabsq	$8031924123371070824, %rax # imm = 0x6F77206F6C6C6568
	movq	%rax, -13(%rbp)
	movb	$10, -1(%rbp)
	movl	$6581362, -5(%rbp)      # imm = 0x646C72
	leaq	-13(%rbp), %rcx
	callq	puts
	nop
	addq	$48, %rsp
	popq	%rbp
	retq
	.seh_handlerdata
	.text
	.seh_endproc
                                        # -- End function
