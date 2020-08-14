package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/pengfeidai/go-snowflake/snowflake"
)

func TestLoad() {
	var wg sync.WaitGroup
	s, err := snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		log.Println(err)
		return
	}
	var check sync.Map
	t1 := time.Now()
	for i := 0; i < 200000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := s.NextVal()
			if _, ok := check.Load(val); ok {
				// id冲突检查
				log.Println(fmt.Errorf("error#unique: val:%v", val))
				return
			}
			check.Store(val, 0)
			if val == 0 {
				log.Println("error")
				return
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(t1)
	log.Printf("generate 20k ids elapsed: %v\n", elapsed)
}

func TestGenID() {
	s, err := snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < 5; i++ {
		val := s.NextVal()
		log.Printf("id: %v, time:%v\n", val, snowflake.GetGenTime(val))
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	TestGenID()
	// 测试生成20万id的速度
	TestLoad()
	// 获取时间戳字段已经使用的占比（0.0 - 1.0）
	// 默认开始时间为：2020年01月01日 00:00:00
	log.Printf("Timestamp status: %f", snowflake.GetTimestampStatus())
}
