// func Add(int32,int32) int32
TEXT ·Add(SB),$0
	MOVL arg1+0(FP),SI
	ADDL arg2+4(FP),SI
	MOVL SI,ret+8(FP)
	RET
