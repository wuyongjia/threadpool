package threadpool

import (
	"errors"
	"sync"
)

type Pool struct {
	Args        chan interface{}
	Threads     int
	QueueLength int
	Func        func(interface{})
}

type PoolSequence struct {
	ArgsList    []chan interface{}
	Id          int
	Threads     int
	QueueLength int
	Func        func(interface{})
	Lock        *sync.RWMutex
}

func NewWithFunc(threads int, queueLength int, poolFunc func(interface{})) *Pool {
	if queueLength < threads {
		panic(errors.New("the maximum queue length must be greater than or equal to the number of threads"))
	}
	var pool = &Pool{
		Args:        make(chan interface{}, queueLength),
		Threads:     threads,
		QueueLength: queueLength,
		Func:        poolFunc,
	}
	pool.start()
	return pool
}

func NewSequenceWithFunc(threads int, queueLength int, poolFunc func(interface{})) *PoolSequence {
	var pool = &PoolSequence{
		ArgsList:    make([]chan interface{}, threads),
		Id:          0,
		Threads:     threads,
		QueueLength: queueLength,
		Func:        poolFunc,
		Lock:        &sync.RWMutex{},
	}
	var i int
	for i = 0; i < threads; i++ {
		pool.ArgsList[i] = make(chan interface{}, queueLength)
	}
	pool.start()
	return pool
}

func (p *Pool) start() {
	var i int
	for i = 0; i < p.Threads; i++ {
		go func() {
			var ok bool
			var args interface{}
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

func (p *PoolSequence) start() {
	var i int
	for i = 0; i < p.Threads; i++ {
		var ch = p.ArgsList[i]
		go func() {
			var ok bool
			var args interface{}
			for {
				args, ok = <-ch
				if ok {
					p.Func(args)
				}
			}
		}()
	}
}

func (p *PoolSequence) GetSequenceId() int {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	var id = p.Id
	p.Id++
	if p.Id >= p.Threads {
		p.Id = 0
	}
	return id
}

func (p *PoolSequence) Invoke(id int, args interface{}) {
	if id >= 0 {
		p.ArgsList[id] <- args
	} else {
		p.ArgsList[p.GetSequenceId()] <- args
	}
}

func (p *PoolSequence) Close() {
	var c chan interface{}
	for _, c = range p.ArgsList {
		close(c)
	}
}
