package worker

import "sync"

type Pool struct {
	wg     sync.WaitGroup
	tasks  chan func()
}

func NewPool(size int) *Pool {
	p := &Pool{
		tasks: make(chan func()),
	}

	for i := 0; i < size; i++ {
		go func() {
			for task := range p.tasks {
				task()
				p.wg.Done()
			}
		}()
	}

	return p
}

func (p *Pool) Submit(task func()) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *Pool) Close() {
	p.wg.Wait()
	close(p.tasks)
}
