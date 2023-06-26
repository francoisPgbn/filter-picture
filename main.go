package main

import (
	"flag"
	"fmt"
	"time"

	filter "photo/filter"
	task "photo/task"
)

func main() {
	var srcDir = flag.String("src", "./imgs", "Input directory")
	var dstDir = flag.String("dst", "./output", "Output directory")
	var filterType = flag.String("filter", "grayscale", "grayscale/blur")
	var taskType = flag.String("task", "channel", "waitgrp/channel")
	var poolSize = flag.Int("poolsize", 4, "Workers pool sizefor the channel task")
	flag.Parse()

	var f filter.Filter

	switch *filterType {
	case "grayscale":
		f = filter.Grayscale{}
	case "blur":
		f = filter.Blur{}
	}

	var t task.Tasker
	switch *taskType {
	case "waitgrp":
		t = task.NewWaitGrpTask(srcDir, dstDir, f)
	case "channel":
		t = task.NewChannelTask(srcDir, dstDir, f, poolSize)
	}

	start := time.Now()
	t.Process()
	elapsed := time.Since(start)
	fmt.Printf("Image processing took %s\n", elapsed)

}
