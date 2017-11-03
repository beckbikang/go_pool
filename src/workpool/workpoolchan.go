package workpool

import (
	"log"
	"sync"
)

//工作池
type KworkpoolChan struct {
	agroup      sync.WaitGroup
	works       chan *Work //一个chan
	processSize int        //启动的work pool的个数
	flag        bool       //是否关闭
}

//新建一个池子
func NewKworkpoolChan(poolSize, processSize int) *KworkpoolChan {
	pool := &KworkpoolChan{
		works:       make(chan *Work, poolSize),
		processSize: processSize,
		flag:        true,
	}

	pool.Start()
	return pool
}

//添加数据
func (kl *KworkpoolChan) AddRunner(w *Work) {
	//发送数据到chan
	if kl.flag {
		kl.works <- w
	}
	log.Println("works_len =", len(kl.works))
}

func (kl *KworkpoolChan) Start() {
	if !kl.flag {
		kl.Close()
		return
	}
	kl.agroup.Add(kl.processSize)

	for i := 0; i < kl.processSize; i++ {
		go func() {
			for runWorker := range kl.works {
				runWorker.Runner.Run(runWorker.Args)
			}
			kl.agroup.Done()
		}()
	}

}

func (kl *KworkpoolChan) Close() {
	if kl.flag != false {
		kl.flag = false
	}
	close(kl.works)
	//等待所有线程完成任务
	kl.agroup.Wait()
}
