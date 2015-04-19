// func Sum(xs []int64) int64
TEXT ·Sum(SB),$0
    MOVQ    $0,SI // SI 寄存器放入 0
    MOVQ    xs+0(FP),BX // 切片八个字节的指针
    MOVQ    BX,autotmp_0000+-16(SP) 
    MOVL    xs+8(FP),BX // 切片八个字节的长度
    MOVL    BX,autotmp_0000+-8(SP)
    MOVL    xs+12(FP),BX // 切片八个字节的容量
    MOVL    BX,autotmp_0000+-4(SP)
    MOVL    $0,AX // 迭代寄存器
    MOVL    autotmp_0000+-8(SP),DI // 长度存入 DI
    LEAQ    autotmp_0000+-16(SP),BX // 指针存入 BX
    MOVQ    (BX),CX // 把 BX 内容移入 CX
    JMP     L2
L1: INCL    AX    // i++
L2: CMPL    AX,DI // 如果 i  == len(xs) ,跳到L3.
    JGE     L3
    MOVQ    (CX),BP // *(CX) = BP 访存,取出内容赋值到BP
    ADDQ    $8,CX   // CX++ 移动8个字节(sizeof int64).
    ADDQ    BP,SI   // SI+=BP,累加结果
    JMP     L1      // 小循环
L3: MOVQ    SI,.noname+16(FP) // 把累加的结果存入返回值里面.
    RET
