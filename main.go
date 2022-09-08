package main

import "utils/logging"

func main() {
	conf := logging.Conf{
		Path:    "log/test.log",
		Encoder: "json",
	}
	logging.Init(conf)
	logging.Errorf("bug 2:%s", "bug")
}
