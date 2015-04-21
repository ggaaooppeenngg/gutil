// func stringToByteSlice(s string) []byte
TEXT ·stringToByteSlice(SB),7,$0
    MOVQ            bs+0(FP),    AX            //; addr
    MOVQ            bs+8(FP),    BX            //; len

    MOVQ            AX,            return+16(FP)    //;addr
    MOVQ            BX,            return+24(FP)    //;len
    MOVQ            BX,            return+32(FP)    //;cap
    RET
