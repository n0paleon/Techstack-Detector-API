package workerpool

import "github.com/panjf2000/ants/v2"

type AntsPool struct {
	pool *ants.Pool
}

func NewAntsPool(size int) *AntsPool {
	p, err := ants.NewPool(size, ants.WithNonblocking(false))
	if err != nil {
		panic(err)
	}

	return &AntsPool{
		pool: p,
	}
}

func (p *AntsPool) Submit(fn func()) error {
	return p.pool.Submit(fn)
}

func (p *AntsPool) Release() {
	p.pool.Release()
}
