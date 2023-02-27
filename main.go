//package main
//
//import "utils/logging"
//
//func main() {
//	conf := logging.Conf{
//		Path:    "log/test.log",
//		Encoder: "json",
//	}
//	logging.Init(conf)
//	logging.Errorf("bug 2:%s", "bug")
//}

package main

import (
	"fmt"
	"os"
	"utils/routers"
)

func main() {
	//tokenBucket := middleware.NewTokenBucket(2, time.Second)
	//for i := 0; i < 3; i++ {
	//	fmt.Println(tokenBucket.IsLimited())
	//}
	//time.Sleep(1 * time.Second)
	//fmt.Println(tokenBucket.IsLimited())

	r := routers.IndexInit()
	fmt.Println("start")
	if err := r.Run(":9092"); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
