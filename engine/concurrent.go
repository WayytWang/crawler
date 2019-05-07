package engine

import (
	"crawler/model"
	"log"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	Submit(Request)
	WorkerChan() chan Request
	ReadyNotifier
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParserResult) //所有worker共用一个输出
	e.Scheduler.Run()              //create requestChan and workerChan

	for i := 0; i < e.WorkerCount; i++ {
		//创建10个worker 每个worker共用一个输出
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			log.Printf("Dupicate request: %s", r.Url)
			continue
		}
		//向Scheduler提交任务 往in管道中放入request
		e.Scheduler.Submit(r)
	}

	profileCount := 0
	for {
		result := <-out
		for _, item := range result.Items {
			if _, ok := item.(model.Profile); ok {
				profileCount++
			}
		}

		//URL dedup
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				continue
			}
			//向Scheduler提交任务 往in管道中放入request
			e.Scheduler.Submit(request)
		}
	}
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}

func createWorker(in chan Request, out chan ParserResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

// 2
// type ConcurrentEngine struct {
// 	Scheduler   Scheduler
// 	WorkerCount int
// }

// type Scheduler interface {
// 	Submit(Request)
// 	ConfigureMasterWorkerChan(chan Request)
// }

// func (e *ConcurrentEngine) Run(seeds ...Request) {
// 	in := make(chan Request)       //所有worker共用一个输入
// 	out := make(chan ParserResult) //所有worker共用一个输出
// 	e.Scheduler.ConfigureMasterWorkerChan(in)

// 	for i := 0; i < e.WorkerCount; i++ {
// 		//创建10个worker 每个worker共用一个输入 和 一个输出
// 		createWorker(in, out)
// 	}

// 	for _, r := range seeds {
// 		//向Scheduler提交任务 往in管道中放入request
// 		e.Scheduler.Submit(r)
// 	}

// 	itemCount := 0
// 	for {
// 		result := <-out
// 		for _, item := range result.Items {
// 			log.Printf("Got item #%d: %v", itemCount, item)
// 			itemCount++
// 		}

// 		for _, request := range result.Requests {
// 			//向Scheduler提交任务 往in管道中放入request
// 			e.Scheduler.Submit(request)
// 		}
// 	}
// }

// func createWorker(in chan Request, out chan ParserResult) {
// 	go func() {
// 		for {
// 			//tell scheduler i'm ready
// 			request := <-in
// 			result, err := worker(request)
// 			if err != nil {
// 				continue
// 			}
// 			out <- result
// 		}
// 	}()
// }

//3
// type ConcurrentEngine struct {
// 	Scheduler   Scheduler
// 	WorkerCount int
// }

// type Scheduler interface {
// 	Submit(Request)
// 	ConfigureMasterWorkerChan(chan Request)
// 	WorkerReady(chan Request)
// 	Run()
// }

// func (e *ConcurrentEngine) Run(seeds ...Request) {
// 	out := make(chan ParserResult) //所有worker共用一个输出
// 	e.Scheduler.Run()              //create requestChan and workerChan

// 	for i := 0; i < e.WorkerCount; i++ {
// 		//创建10个worker 每个worker共用一个输出
// 		createWorker(out, e.Scheduler)
// 	}

// 	for _, r := range seeds {
// 		//向Scheduler提交任务 往in管道中放入request
// 		e.Scheduler.Submit(r)
// 	}

// 	itemCount := 0
// 	for {
// 		result := <-out
// 		for _, item := range result.Items {
// 			log.Printf("Got item #%d: %v", itemCount, item)
// 			itemCount++
// 		}

// 		for _, request := range result.Requests {
// 			//向Scheduler提交任务 往in管道中放入request
// 			e.Scheduler.Submit(request)
// 		}
// 	}
// }

// func createWorker(out chan ParserResult, s Scheduler) {
// 	//每一个worker都有一个chan
// 	in := make(chan Request)
// 	go func() {
// 		for {
// 			//tell scheduler i'm ready
// 			s.WorkerReady(in)
// 			request := <-in
// 			result, err := worker(request)
// 			if err != nil {
// 				continue
// 			}
// 			out <- result
// 		}
// 	}()
// }
