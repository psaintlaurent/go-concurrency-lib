package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

/*
	This takes advantage of a FIFO worker queue pattern to ensure ordered execution.
*/

func main() {

	jobCh := make(chan chan int64, 100)
	go sequentialWorker(jobCh)
	r := gin.Default()

	r.GET("/",
		func(c *gin.Context) {

			retChan := make(chan int64)
			jobCh <- retChan
			count := <-retChan
			c.JSON(200, gin.H{"message": fmt.Sprintf("Holla!, we have hit %d times.", count)})
		})

	r.Run()
}

func sequentialWorker(jobCh <-chan chan int64) {

	f, err := os.OpenFile("./keep_count", os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var count int64
	for job := range jobCh {

		_, err1 := fmt.Fscanf(f, "%d\n", &count)
		count++
		if err1 != nil {
			log.Println(err1)
		}
		_ = f.Truncate(0)
		_, err = f.Seek(0, 0)
		f.Write([]byte(fmt.Sprintf("%d\n", count)))
		job <- count
	}
}
