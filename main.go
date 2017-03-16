package main

import (
	"github.com/reechou/robot-train/config"
	"github.com/reechou/robot-train/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
