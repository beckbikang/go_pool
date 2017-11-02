package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Runner interface {
	Run(interface{}) (interface{}, error)
}

//定义一个有趣的函数类型，这个类型实现了接口 代理模式
type RunnerFunc func(interface{}) (interface{}, error)

func (r RunnerFunc) Run(avar interface{}) (interface{}, error) {
	return r(avar)
}

type Work struct {
	Runner Runner
	Args   interface{}
}

func NewWork(runner Runner, args interface{}) *Work {
	return &Work{runner, args}
}

//工作池
type Kworkpool struct {
	rmutex     sync.RWMutex
	mutex      sync.Mutex //写锁
	runnerList *list.List //任务列表
	poolSize   int        //启动的pool的个数
	flag       bool       //是否关闭
}

//新建一个池子
func NewKworkpool(size int) *Kworkpool {
	l := list.New()
	return &Kworkpool{
		poolSize:   size,
		flag:       true,
		runnerList: l,
	}
}

//添加数据
func (kl *Kworkpool) AddRunner(w *Work) {
	if kl.flag {
		kl.runnerList.PushFront(w)
	}
}

func (kl *Kworkpool) Start() {
	if !kl.flag {
		kl.Close()
		return
	}
	kl.run()
}

//启动多个goroutine
func (kl *Kworkpool) run() {
	for i := 0; i < kl.poolSize; i++ {
		go kl.work()
	}
}

//实际的工作脚本运行pool
func (kl *Kworkpool) work() {
	for {
		//检测数据的长度
		kl.rmutex.RLock()
		listLen := kl.runnerList.Len()
		kl.rmutex.RUnlock()
		if listLen == 0 { //休眠100毫秒
			//已经关闭就结束程序,需要判断是否已经没有任务
			if !kl.flag {
				break
			}
			time.Sleep(time.Millisecond * 100)
			continue
		}
		kl.mutex.Lock()
		elem := kl.runnerList.Back()
		if elem == nil {
			kl.mutex.Unlock()
			continue
		}
		worker := kl.runnerList.Remove(elem).(*Work)
		kl.mutex.Unlock()
		log.Println(" start running ...")
		worker.Runner.Run(worker.Args)
	}
}

func (kl *Kworkpool) Close() {
	kl.flag = false
}

//测试runner
type Arunner struct {
}

func (a *Arunner) Run(avar interface{}) (interface{}, error) {
	i := avar.(int)
	filename := "./log/" + strconv.Itoa(time.Now().Nanosecond())
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("go", i)
		return nil, err
	}
	fp.WriteString("go" + strconv.Itoa(i) + "\n")
	fp.Close()
	return nil, nil
}

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	fmt.Println("start")
	pool := NewKworkpool(20)
	pool.Start()

	var run1, run2 Runner
	run1 = &Arunner{}
	run2 = &Arunner{}

	w1 := NewWork(run1, 1)
	w2 := NewWork(run2, 2)
	pool.AddRunner(w1)
	pool.AddRunner(w2)
	fmt.Print(pool)
	time.Sleep(time.Second * 5)
}
