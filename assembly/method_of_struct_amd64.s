// fucn (t *T)Get()int
TEXT ·T·Get(SB), 7, $0
	MOVQ	p+0(FP), AX
	MOVQ	AX, ret+8(FP)
	RET
