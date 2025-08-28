package worker

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тест: одна задача выполняется
func TestPoolSingleTask(t *testing.T) {
	var called int32

	p := NewPool(1)
	p.Submit(func() {
		atomic.AddInt32(&called, 1)
	})
	p.Close()

	assert.Equal(t, int32(1), called, "task should be executed once")
}

// Тест: несколько задач выполняются
func TestPoolMultipleTasks(t *testing.T) {
	var counter int32
	p := NewPool(3) // несколько воркеров

	for i := 0; i < 10; i++ {
		p.Submit(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	p.Close()

	assert.Equal(t, int32(10), counter, "all tasks should be executed")
}

// Тест: Close ждёт выполнения всех задач
func TestPoolCloseWaits(t *testing.T) {
	var counter int32
	p := NewPool(2)

	p.Submit(func() {
		time.Sleep(200 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	p.Submit(func() {
		time.Sleep(200 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	start := time.Now()
	p.Close()
	elapsed := time.Since(start)

	assert.Equal(t, int32(2), counter, "both tasks should finish before Close returns")
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(200), "Close should wait for tasks")
}

// Тест: задачи выполняются параллельно
func TestPoolParallelExecution(t *testing.T) {
	var counter int32
	p := NewPool(2) // два воркера

	start := time.Now()

	// Две долгие задачи
	p.Submit(func() {
		time.Sleep(300 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})
	p.Submit(func() {
		time.Sleep(300 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	p.Close()
	elapsed := time.Since(start)

	assert.Equal(t, int32(2), counter)
	// Если задачи реально параллельные, общее время < 600 мс
	assert.Less(t, elapsed.Milliseconds(), int64(600), "tasks should run in parallel")
}
