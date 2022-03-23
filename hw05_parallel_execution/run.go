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

func Run(tasks []Task, n, m int) error {
	chTask := make(chan Task)
	errorCounter := &ErrorCounter{}
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		go func() {
			for {
				task, ok := <-chTask
				if ok == false {
					return
				}
				if err := task(); err != nil {
					errorCounter.Add()
				}
				wg.Done()
			}
		}()
	}
	var err error
	for _, task := range tasks {
		if m <= 0 || errorCounter.count >= m {
			err = ErrErrorsLimitExceeded
			break
		}
		wg.Add(1)
		chTask <- task
	}
	wg.Wait()
	close(chTask)

	return err
}
