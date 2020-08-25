package threadpool

import (
	"errors"
)

type Pool struct {
	Args    chan interface{}
	Threads int
	Func    func(interface{})
}

func NewWithFunc(numberOfThreads int, maxQueueLength int, poolFunc func(interface{})) *Pool {
	if maxQueueLength < numberOfThreads {
		panic(errors.New("the maximum queue length must be greater than or equal to the number of threads"))
	}
	var pool = &Pool{
		Args:    make(chan interface{}, maxQueueLength),
		Threads: numberOfThreads,
		Func:    poolFunc,
	}
	pool.start()
	return pool
}

func (p *Pool) start() {
	var i int
	var ok bool
	var args interface{}
	for i = 0; i < p.Threads; i++ {
		go func() {
			for {
				args, ok = <-p.Args
				if ok {
					p.Func(args)
				}
			}
		}()
	}
}

func (p *Pool) Invoke(args interface{}) {
	p.Args <- args
}

func (p *Pool) Close() {
	close(p.Args)
}
