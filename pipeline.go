package gutil

import (
	"errors"
	"reflect"
	"sync"
)

//pipeline is inspired by this blog http://blog.golang.org/pipelines

//流水线模式.
//一开始是接受输入,最后一个是产生输出,中间通过通道连接.
//注意:goroutine或者channel没有关闭都会造成leak.
//流水线的模式:

//1.  每一步关闭输出通道,当所有发送动作完成的时候.
//2.  接受输入的通道的内容直到通道关闭.

//构造方法
// 每个阶段关闭输出通道,当发送操作完成的时候.
// 每个阶段不断接受输入通道的值,直到这些通道关闭,或者发送没有阻塞(也就是可以从done里读东西,或者可以从errc里面读东西)

//有两种,一种是fan-int,多个channel合并到一个,一种是fan-out,一个channel输出到多个.

//还有一种是一个channel开n个goroutine共享,用waitGroup.wait在一个goroutine里面等待关闭输出.
//最主要的关键是接受输入直到关闭,产生输出之后close通道.

//最后一个是显示关闭,就是用select防止一些发送者发不出去导致goroutine leak的情况,
//在range里面加一个select{case <-done:return},然后关闭对应的输出.

//错误处理,对于每个阶段的单利error和消息绑定成新的结构出输出下去.
//对于这个阶段的整体错误,给出一个error通道来接受.
//下个阶段得到任意的error都直接return导致,输出关闭或者是引起done关闭.

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

//sum the numbers and send to output.
func Sum(sum chan []int) (output chan []int) {
	output = make(chan []int)
	go func() {
		var total []int
		count := 0
		for slc := range sum {
			count++
			if len(total) < len(slc) {
				tmp := make([]int, len(slc))
				copy(tmp, total)
				total = tmp
			}
			for i, v := range slc {
				total[i] += v
			}
		}
		if count > 0 {
			output <- total
		}
		close(output)
	}()
	return
}

//sum numbers parallel
func ParallelSum(slcs ...[]int) []int {
	input := make(chan []int, len(slcs))
	output := make(chan []int)
	var result []int
	go func(input chan []int) {
		for _, slc := range slcs {
			input <- slc
		}
		close(input)
	}(input)

	for {
		var wg sync.WaitGroup
		wg.Add(cap(input) / 2)
		for i := 0; i < cap(input)/2; i++ {
			out := Sum(input)
			go func() {
				defer wg.Done()
				for o := range out {
					output <- o
				}
			}()
		}
		go func(output chan []int) {
			wg.Wait()
			close(output)
		}(output)

		input = make(chan []int, cap(input)/2)
		if cap(input) < 2 {
			result = <-output
			break
		}
		for o := range output {
			input <- o
		}
		output = make(chan []int)
		close(input)
	}
	return result
}
