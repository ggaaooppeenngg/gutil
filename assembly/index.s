// func Index(a []int64,c)int64
TEXT ·Index(SB),$0-32
	MOVQ a+0(FP),SI // data
	MOVQ a_len+8(FP),BX // len
	MOVQ c+24(FP),AX // int64 sought
	MOVQ SI,DI	
	MOVQ BX,CX //   重复CX次SCASQ
	REPN; SCASQ 
	//  这是两个命令的组合，不是注释.
	//  每次和AX比较,如果不对就递减CX.
	JZ success
	MOVQ $-1,AX
	RET
success:
	SUBQ SI,DI 
	SUBL $1,DI
	MOVQ DI,AX
	RET
