package worker

import "sync"

// Pool — интерфейс пула, удобный для моков
type Pool interface {
	Submit(func())
	Close()
}

// реализация
type pool struct {
	wg    sync.WaitGroup
	tasks chan func()
}

// NewPool создаёт пул воркеров
func NewPool(size int) Pool {
	p := &pool{
		tasks: make(chan func(), size),
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

func (p *pool) Submit(task func()) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *pool) Close() {
	p.wg.Wait()
	close(p.tasks)
}
