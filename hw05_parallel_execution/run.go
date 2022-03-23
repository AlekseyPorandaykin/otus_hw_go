package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorCounter struct {
	mt    sync.Mutex
	count int
}

func (err *ErrorCounter) Add() {
	defer err.mt.Unlock()
	err.mt.Lock()
	err.count++
}

type Executor struct {
	tasks             []Task
	quantityGoroutine int
	quantityErrors    int
	errorCounter      *ErrorCounter
	ch                chan Task
}

func (executor *Executor) isErrorLimitExceeded() bool {
	if executor.quantityErrors <= 0 || executor.errorCounter.count >= executor.quantityErrors {
		return true
	}

	return false
}

func NewExecutor(tasks []Task, quantityGoroutine int, quantityErrors int) *Executor {
	return &Executor{
		tasks:             tasks,
		quantityGoroutine: quantityGoroutine,
		quantityErrors:    quantityErrors,
		errorCounter:      &ErrorCounter{},
		ch:                make(chan Task),
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	executor := NewExecutor(tasks, n, m)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		consumer(executor, &wg)
	}
	err := producer(executor, &wg)
	wg.Wait()
	close(executor.ch)

	return err
}

func consumer(executor *Executor, wg *sync.WaitGroup) {
	go func() {
		for {
			task, ok := <-executor.ch
			if ok == false {
				return
			}
			if err := task(); err != nil {
				executor.errorCounter.Add()
			}
			wg.Done()
		}
	}()
}

func producer(executor *Executor, wg *sync.WaitGroup) error {
	for _, task := range executor.tasks {
		if executor.isErrorLimitExceeded() {
			return ErrErrorsLimitExceeded
		}
		wg.Add(1)
		executor.ch <- task
	}

	return nil
}
