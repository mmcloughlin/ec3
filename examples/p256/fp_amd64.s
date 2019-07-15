// Code generated by ec3. DO NOT EDIT.

#include "textflag.h"

// func Add(x *Elt, y *Elt)
TEXT ·Add(SB), NOSPLIT, $0-16
	MOVQ    x+0(FP), AX
	MOVQ    y+8(FP), CX
	MOVQ    (AX), DX
	MOVQ    8(AX), BX
	MOVQ    16(AX), BP
	MOVQ    24(AX), SI
	MOVQ    (CX), DI
	MOVQ    8(CX), R8
	MOVQ    16(CX), R9
	MOVQ    24(CX), CX
	XORQ    R10, R10
	ADDQ    DI, DX
	ADCQ    R8, BX
	ADCQ    R9, BP
	ADCQ    CX, SI
	ADCQ    $0x00000000, R10
	MOVQ    DX, CX
	MOVQ    BX, DI
	MOVQ    BP, R8
	MOVQ    SI, R9
	SUBQ    p<>+0(SB), CX
	SBBQ    p<>+8(SB), DI
	SBBQ    p<>+16(SB), R8
	SBBQ    p<>+24(SB), R9
	SBBQ    $0x00000000, R10
	CMOVQCC CX, DX
	CMOVQCC DI, BX
	CMOVQCC R8, BP
	CMOVQCC R9, SI
	MOVQ    DX, (AX)
	MOVQ    BX, 8(AX)
	MOVQ    BP, 16(AX)
	MOVQ    SI, 24(AX)
	RET

DATA p<>+0(SB)/8, $0xffffffffffffffff
DATA p<>+8(SB)/8, $0x00000000ffffffff
DATA p<>+16(SB)/8, $0x0000000000000000
DATA p<>+24(SB)/8, $0xffffffff00000001
GLOBL p<>(SB), RODATA|NOPTR, $32

// func Sub(x *Elt, y *Elt)
TEXT ·Sub(SB), NOSPLIT, $0-16
	MOVQ    x+0(FP), AX
	MOVQ    y+8(FP), CX
	MOVQ    (AX), DX
	MOVQ    8(AX), BX
	MOVQ    16(AX), BP
	MOVQ    24(AX), SI
	MOVQ    (CX), DI
	MOVQ    8(CX), R8
	MOVQ    16(CX), R9
	MOVQ    24(CX), CX
	XORQ    R10, R10
	SUBQ    DI, DX
	SBBQ    R8, BX
	SBBQ    R9, BP
	SBBQ    CX, SI
	SBBQ    $0x00000000, R10
	MOVQ    DX, CX
	MOVQ    BX, DI
	MOVQ    BP, R8
	MOVQ    SI, R9
	ADDQ    p<>+0(SB), CX
	ADCQ    p<>+8(SB), DI
	ADCQ    p<>+16(SB), R8
	ADCQ    p<>+24(SB), R9
	ANDQ    $0x00000001, R10
	CMOVQNE CX, DX
	CMOVQNE DI, BX
	CMOVQNE R8, BP
	CMOVQNE R9, SI
	MOVQ    DX, (AX)
	MOVQ    BX, 8(AX)
	MOVQ    BP, 16(AX)
	MOVQ    SI, 24(AX)
	RET

// func Mul(z *Elt, x *Elt, y *Elt)
TEXT ·Mul(SB), NOSPLIT, $64-24
	MOVQ z+0(FP), AX
	MOVQ x+8(FP), CX
	MOVQ y+16(FP), BX

	// y[0]
	MOVQ (BX), DX
	XORQ BP, BP

	// x[0] * y[0] -> z[0]
	MULXQ (CX), SI, DI

	// x[1] * y[0] -> z[1]
	MULXQ 8(CX), R8, R9
	ADCXQ R8, DI

	// x[2] * y[0] -> z[2]
	MULXQ 16(CX), R8, R11
	ADCXQ R8, R9

	// x[3] * y[0] -> z[3]
	MULXQ 24(CX), DX, R8
	ADCXQ DX, R11
	ADCXQ BP, R8
	MOVQ  SI, (SP)

	// y[1]
	MOVQ 8(BX), DX
	XORQ BP, BP

	// x[0] * y[1] -> z[1]
	MULXQ (CX), SI, R12
	ADCXQ SI, DI
	ADOXQ R12, R9

	// x[1] * y[1] -> z[2]
	MULXQ 8(CX), SI, R12
	ADCXQ SI, R9
	ADOXQ R12, R11

	// x[2] * y[1] -> z[3]
	MULXQ 16(CX), SI, R12
	ADCXQ SI, R11
	ADOXQ R12, R8

	// x[3] * y[1] -> z[4]
	MULXQ 24(CX), DX, SI
	ADCXQ DX, R8
	ADCXQ BP, SI
	ADOXQ BP, SI
	MOVQ  DI, 8(SP)

	// y[2]
	MOVQ 16(BX), DX
	XORQ BP, BP

	// x[0] * y[2] -> z[2]
	MULXQ (CX), DI, R12
	ADCXQ DI, R9
	ADOXQ R12, R11

	// x[1] * y[2] -> z[3]
	MULXQ 8(CX), DI, R12
	ADCXQ DI, R11
	ADOXQ R12, R8

	// x[2] * y[2] -> z[4]
	MULXQ 16(CX), DI, R12
	ADCXQ DI, R8
	ADOXQ R12, SI

	// x[3] * y[2] -> z[5]
	MULXQ 24(CX), DX, DI
	ADCXQ DX, SI
	ADCXQ BP, DI
	ADOXQ BP, DI
	MOVQ  R9, 16(SP)

	// y[3]
	MOVQ 24(BX), DX
	XORQ BP, BP

	// x[0] * y[3] -> z[3]
	MULXQ (CX), BX, R9
	ADCXQ BX, R11
	ADOXQ R9, R8

	// x[1] * y[3] -> z[4]
	MULXQ 8(CX), BX, R9
	ADCXQ BX, R8
	ADOXQ R9, SI

	// x[2] * y[3] -> z[5]
	MULXQ 16(CX), BX, R9
	ADCXQ BX, SI
	ADOXQ R9, DI

	// x[3] * y[3] -> z[6]
	MULXQ 24(CX), CX, DX
	ADCXQ CX, DI
	ADCXQ BP, DX
	ADOXQ BP, DX
	MOVQ  R11, 24(SP)
	MOVQ  R8, 32(SP)
	MOVQ  SI, 40(SP)
	MOVQ  DI, 48(SP)
	MOVQ  DX, 56(SP)

	// Reduction.
	MOVQ    (SP), CX
	MOVQ    8(SP), BX
	MOVQ    16(SP), BP
	MOVQ    24(SP), SI
	MOVQ    32(SP), DI
	MOVQ    40(SP), R8
	MOVQ    CX, DX
	XORQ    R10, R10
	MULXQ   p<>+0(SB), R9, R11
	ADCXQ   R9, CX
	ADOXQ   R11, BX
	MULXQ   p<>+8(SB), CX, R9
	ADCXQ   CX, BX
	ADOXQ   R9, BP
	MULXQ   p<>+16(SB), CX, R9
	ADCXQ   CX, BP
	ADOXQ   R9, SI
	MULXQ   p<>+24(SB), CX, DX
	ADCXQ   CX, SI
	ADOXQ   DX, DI
	ADCXQ   R10, DI
	ADOXQ   R10, R8
	MOVQ    48(SP), CX
	MOVQ    BX, DX
	XORQ    R10, R10
	MULXQ   p<>+0(SB), R9, R11
	ADCXQ   R9, BX
	ADOXQ   R11, BP
	MULXQ   p<>+8(SB), BX, R9
	ADCXQ   BX, BP
	ADOXQ   R9, SI
	MULXQ   p<>+16(SB), BX, R9
	ADCXQ   BX, SI
	ADOXQ   R9, DI
	MULXQ   p<>+24(SB), DX, BX
	ADCXQ   DX, DI
	ADOXQ   BX, R8
	ADCXQ   R10, R8
	ADOXQ   R10, CX
	MOVQ    56(SP), BX
	MOVQ    BP, DX
	XORQ    R10, R10
	MULXQ   p<>+0(SB), R9, R11
	ADCXQ   R9, BP
	ADOXQ   R11, SI
	MULXQ   p<>+8(SB), BP, R9
	ADCXQ   BP, SI
	ADOXQ   R9, DI
	MULXQ   p<>+16(SB), BP, R9
	ADCXQ   BP, DI
	ADOXQ   R9, R8
	MULXQ   p<>+24(SB), DX, BP
	ADCXQ   DX, R8
	ADOXQ   BP, CX
	ADCXQ   R10, CX
	ADOXQ   R10, BX
	MOVQ    $0x00000000, BP
	MOVQ    SI, DX
	XORQ    R10, R10
	MULXQ   p<>+0(SB), R9, R11
	ADCXQ   R9, SI
	ADOXQ   R11, DI
	MULXQ   p<>+8(SB), SI, R9
	ADCXQ   SI, DI
	ADOXQ   R9, R8
	MULXQ   p<>+16(SB), SI, R9
	ADCXQ   SI, R8
	ADOXQ   R9, CX
	MULXQ   p<>+24(SB), DX, SI
	ADCXQ   DX, CX
	ADOXQ   SI, BX
	ADCXQ   R10, BX
	ADOXQ   R10, BP
	MOVQ    DI, DX
	MOVQ    R8, SI
	MOVQ    CX, R9
	MOVQ    BX, R10
	SUBQ    p<>+0(SB), DX
	SBBQ    p<>+8(SB), SI
	SBBQ    p<>+16(SB), R9
	SBBQ    p<>+24(SB), R10
	SBBQ    $0x00000000, BP
	CMOVQCC DX, DI
	CMOVQCC SI, R8
	CMOVQCC R9, CX
	CMOVQCC R10, BX
	MOVQ    DI, (AX)
	MOVQ    R8, 8(AX)
	MOVQ    CX, 16(AX)
	MOVQ    BX, 24(AX)
	RET