package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"workpool"
)

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

func main() {
	//tpool1()
	//tpool2()
	tpool3()
}

//创建一个connect
var connectId int32 = 0

type DbConnect struct {
	cid int32
}

func (db *DbConnect) getConnectId() int32 {
	return db.cid
}

func (db *DbConnect) Close() error {
	return nil
}

//创建连接
type DbConnectCreater struct{}

func (dbc *DbConnectCreater) CreateConnect() (io.Closer, error) {

	id := atomic.AddInt32(&connectId, 1)

	return &DbConnect{cid: id}, nil
}

func tpool3() {
	fmt.Println("test connect pool3")
	creater := &DbConnectCreater{}
	pool3, err := workpool.NewConnectPool(10, creater)
	if err != nil {
		log.Println("create faild")
	}
	defer pool3.Close()
	//放入10个连接
	for i := 0; i < 10; i++ {
		conect, err := creater.CreateConnect()
		if err == nil {
			pool3.PutConnect(conect)
		}
	}

	for i := 0; i < 12; i++ {
		connect, err := pool3.GetConnect()
		if err == nil {
			log.Println("connect_id=", connect.(*DbConnect).getConnectId())

		}
	}

}

func tpool2() {
	//test pool 2
	var run1, run2 workpool.Runner
	run1 = &Arunner{}
	run2 = &Arunner{}

	pool2 := workpool.NewKworkpoolChan(8, 20)

	w1 := workpool.NewWork(run1, 1)
	w2 := workpool.NewWork(run2, 2)
	pool2.AddRunner(w1)
	pool2.AddRunner(w2)
	fmt.Println(pool2)

	time.Sleep(time.Second * 5)
}

func tpool1() {
	//test pool 1
	pool := workpool.NewKworkpool(20)
	pool.Start()
	var run1, run2 workpool.Runner
	run1 = &Arunner{}
	run2 = &Arunner{}

	w1 := workpool.NewWork(run1, 1)
	w2 := workpool.NewWork(run2, 2)
	pool.AddRunner(w1)
	pool.AddRunner(w2)
	fmt.Print(pool)
}
