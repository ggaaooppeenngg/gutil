package util

import (
	"errors"
	"reflect"
)

//pipeline is inspired by this blog http://blog.golang.org/pipelines

//流水线模式.
//一开始是接受输入,最后一个是产生输出,中间通过通道连接.
//注意:goroutine或者channel没有关闭都会造成leak.
//流水线的模式:

//1.  每一步关闭输出通道,当所有发送动作完成的时候.
//2.  接受输入的通道的内容直到通道关闭.
//需要注意的地方
//一个channel,每个工作者拿来range,等range结束就可以关闭输出,但是产生多个channel还要来合并.

//还有一种是一个channel开n个goroutine共享,用waitGroup.wait在一个goroutine里面等待关闭输出.

//最主要的关键是接受输入直到关闭,产生输出之后close通道.

//最后一个是显示关闭,就是用select防止一些发送者发不出去导致goroutine leak的情况,

//在range里面加一个select{case <-done:return},然后关闭对应的输出.

type pf func(inputs chan interface{}) chan interface{}

type PL struct {
	err      error
	workers  []interface{}
	channels []interface{}
}

func (p *PL) Input(inpus interface{}) {

}

func (p *PL) Outputs() {

}
func (p *PL) Pipe(worker interface{}, max int) {
	fv := reflect.ValueOf(worker)
	if fv.Kind() != reflect.Func {
		p.err = errors.New("worker type is not func")
		return
	}
	//TODO:type check
	p.workers = append(p.workers, worker)
	//TODO:get func argument
	p.channels = append(p.channels, make(chan int))
}
