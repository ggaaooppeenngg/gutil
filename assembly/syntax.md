一下内容均针对的是AMD64架构的Plan 9 汇编.

这个[库](https://github.com/ziutek/blas)是基础线性代数库, 全部用汇编重写，对于写汇编是一个非常重要的参考.

汇编手册 https://9p.io/sys/doc/asm.html
简易教程 https://golang.org/doc/asm

Go 的汇编里面, 参数的命名是下划线式的.

FP: 帧指针 参数和本地变量
SP: 栈指针 比如arg1+0(SP)是函数的第一个参数.
PC: 指令计数器.
SB: 静态地址 可以认为是内存的首地址.`foo(SB)`可以当做是foo在内存当中的地址.foo<>(SB)表示只在当前文件可见,类似于C的static.

通过 0(FP) 可以获取参数的的一个参数. first\_arg+0(FP) 可以当作第一个参数的一个名字.
SP指向的是栈顶, 所以要用x-8(SP)这样的形式访问变量.

TLS: Thread Local Storage,可以获取线程的栈大小用来扩展栈.
AMD64里面,DI(%rdi),SI(%rsi),DX(%rdx),CX(%rcx),R8(%r8),R9(%r9)用作函数参数,但是Plan 9里面也是一样的，被调用者拿参数就使用栈帧.

函数定义里的$0-32 表示frame size 0 参数的大小是32.
The frame size $24-8 states that the function has a 24-byte frame and is called with 8 bytes of argument, which live on the caller's frame.
If NOSPLIT is not specified for the TEXT, the argument size must be provided. For assembly functions with Go prototypes, go vet will check that the argument size is correct.

DATA    symbol+offset(SB)/width, value
初始化地址上width长度的内容.

FUNCDATA $2, $sym(SB)
表示fundata列表的index 2的位置用sym这个符号.
PCDATA $3 $45
表示PCDATA列表对应的index 4的值是45

对应的文档在  https://docs.google.com/document/d/1lyPIbmsYbXnpNj57a261hgOYVpNRcgydurVQIyZOz_o/pub

http://blog.altoros.com/golang-part-4-object-files-and-function-metadata.html
这个thread 和GC有关, 等看懂GC再来处理继续挖掘这个问题.
https://groups.google.com/forum/#!topic/golang-nuts/a6NKBbL9fX0


AX(%rax)在AMD64里面是用来函数返回值，但是plan9里面是用栈帧或者AX,比如bytes.Equal
里面用的就是AX.

BX(%rbx),BP(%rpx),R12(%r12),R13(%r13),R14(%r14),R15(%r15),是数据寄存器，随便存东西，在Plan 9里面也是对应的.

如果是切片作参数，依次是  地址|长度|容量 每个都是8字节，总共24个字节.
如果字符串是参数，依次是  地址|长度| 字符串没有容量.

函数调用的时候:
x86的规则是这样的:寄存器分为调用者保存寄存器和被调用者保存寄存器。按
照惯例,eax,edx,ecx寄存器是调用者保存,ebx,esi,edi,ebp等寄存器是被调用者
负责保存。

附录:
常数特征:
64 位的最小数,−9,223,372,036,854,775,808,对应的补码是1000000...
   最大数后面不一样的是807.
-1 是全1.
指令:
SCAS:比较AX和DI当中的值,然后设置flags.
