package task

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	filter "photo/filter"
)

type Tasker interface {
	Process() error
}

type TaskWaitGroup struct {
	srcDir  *string
	destDir *string
	filter  filter.Filter
}

type TaskChannel struct {
	srcDir   *string
	destDir  *string
	poolSize *int
	filter   filter.Filter
}

func NewWaitGrpTask(srcDir *string, destDir *string, f filter.Filter) Tasker {
	return &TaskWaitGroup{
		srcDir:  srcDir,
		destDir: destDir,
		filter:  f,
	}
}

func NewChannelTask(srcDir *string, destDir *string, f filter.Filter, poolsize *int) Tasker {
	return &TaskChannel{
		srcDir:   srcDir,
		destDir:  destDir,
		filter:   f,
		poolSize: poolsize,
	}
}

func (task *TaskWaitGroup) Process() error {
	fmt.Println("WaitGroup");

	files, err := ioutil.ReadDir(*task.srcDir)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du dossier source %s: %s\n", task.srcDir, err.Error())
		return err
	}

	var wg sync.WaitGroup

	for _, file := range files {
		if !file.IsDir() {
			inputPath := filepath.Join(*task.srcDir, file.Name())
			outputPath := filepath.Join(*task.destDir, file.Name())

			wg.Add(1)
			go func(inputPath, outputPath string) {
				defer wg.Done()

				task.filter.Process(inputPath, outputPath)

			}(inputPath, outputPath)
		}
	}

	wg.Wait()
	return nil
}

func (task *TaskChannel) Process() error {
	fmt.Println("TaskChannel");

	files, err := ioutil.ReadDir(*task.srcDir)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du dossier source %s: %s\n", *task.srcDir, err.Error())
		return nil
	}

	jobs := make(chan string)
	done := make(chan bool)
	fileNbr := 0

	for _, file := range files {
		if !file.IsDir() {
			fileNbr++
		}
	}

	nbrDone := make(chan bool, fileNbr)

	// Fonction pour traiter les images en utilisant les canaux
	worker := func() {
		for fileName := range jobs {
			inputPath := filepath.Join(*task.srcDir, fileName)
			outputPath := filepath.Join(*task.destDir, fileName)
			task.filter.Process(inputPath, outputPath)
			nbrDone <- true
		}
		done <- true
	}

	// Lancement des goroutines pour le traitement des images
	for i := 0; i < *task.poolSize; i++ {
		go worker()
	}

	for _, file := range files {
		if !file.IsDir() {
			jobs <- file.Name()
		}
	}

	for index := 0; index < fileNbr; index++ {
		<-nbrDone
	}

	close(jobs)
	return nil
}
