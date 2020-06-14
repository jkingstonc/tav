	.text
	.def	 @feat.00;
	.scl	3;
	.type	0;
	.endef
	.globl	@feat.00
.set @feat.00, 0
	.file	"test.ll"
	.def	 test;
	.scl	2;
	.type	32;
	.endef
	.section	.rdata,"dr"
	.p2align	2               # -- Begin function test
.LCPI0_0:
	.long	1067869798              # float 1.29999995
	.text
	.globl	test
	.p2align	4, 0x90
test:                                   # @test
# %bb.0:
	movss	.LCPI0_0(%rip), %xmm0   # xmm0 = mem[0],zero,zero,zero
	retq
                                        # -- End function
	.def	 main;
	.scl	2;
	.type	32;
	.endef
	.globl	main                    # -- Begin function main
	.p2align	4, 0x90
main:                                   # @main
.seh_proc main
# %bb.0:
	pushq	%rbp
	.seh_pushreg %rbp
	subq	$32, %rsp
	.seh_stackalloc 32
	leaq	32(%rsp), %rbp
	.seh_setframe %rbp, 32
	.seh_endprologue
	callq	__main
	movl	$1, %eax
	addq	$32, %rsp
	popq	%rbp
	retq
	.seh_handlerdata
	.text
	.seh_endproc
                                        # -- End function
