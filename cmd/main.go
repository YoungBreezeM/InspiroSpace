package main

import (
	"easygin/internal/boot"
	"flag"
	"fmt"
)

func main() {
	configPath := ""
	flag.StringVar(&configPath, "c", "", "config path")
	flag.Parse()
	//
	a, err := boot.BootStrap(configPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	a.Run()
}
