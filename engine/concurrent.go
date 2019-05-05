package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)       //所有worker共用一个输入
	out := make(chan ParserResult) //所有worker共用一个输出
	e.Scheduler.ConfigureMasterWorkerChan(in)

	for i := 0; i < e.WorkerCount; i++ {
		//创建10个worker 每个worker共用一个输入 和 一个输出
		createWorker(in, out)
	}

	for _, r := range seeds {
		//向Scheduler提交任务 往in管道中放入request
		e.Scheduler.Submit(r)
	}

	itemCount := 0
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item #%d: %v", itemCount, item)
			itemCount++
		}

		for _, request := range result.Requests {
			//向Scheduler提交任务 往in管道中放入request
			e.Scheduler.Submit(request)
		}
	}
}

func createWorker(in chan Request, out chan ParserResult) {
	go func() {
		for {
			//tell scheduler i'm ready
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
