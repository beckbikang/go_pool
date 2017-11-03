package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	fmt.Println("start")
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
	time.Sleep(time.Second * 5)
}
