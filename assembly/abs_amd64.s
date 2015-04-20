// func Abs( a int64) int64
TEXT ·Abs(SB),$0 
	MOVQ arg1+0(FP),AX
	CMPQ AX,$0
	JGE ret
	NEGQ AX
ret:
	MOVQ AX,abs+8(FP)
	RET	
	
