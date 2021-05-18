package main

import (
	"context"
	logx "demo/etcdwatch/logx"
	"github.com/etcd-io/etcd/clientv3"
	"golang.org/x/time/rate"
	"log"
	"time"
)

func rateTester() {
	//limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)
	limiter := rate.NewLimiter(10, 1)
	var err error
	err = limiter.Wait(context.Background())
	if err != nil {
		return
	}

}

func main() {
	logx.InitLogger("data/logs/etcdwatch.log", "debug")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.0.84.174:2379", "10.0.84.174:23279", "10.0.84.174:33279"}, //etcd集群三个实例的端口
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		log.Println("connect failed, err:", err)
		return
	}

	logx.X.Info("connect succ")

	defer cli.Close()

	for {
		rch := cli.Watch(context.Background(), "/develop_staging", clientv3.WithPrefix()) //阻塞在这里，如果没有key里没有变化，就一直停留在这里
		for wresp := range rch {
			for _, ev := range wresp.Events {
				log.Printf("%v %q:%q\n", ev, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
