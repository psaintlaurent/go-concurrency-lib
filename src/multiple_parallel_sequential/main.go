package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

/*
	This takes advantage of a FIFO worker queue pattern to ensure parallel ordered execution.
*/


type Job struct {

	id string
	fn chan int64
	message chan int64
	status chan Status
	complete bool
}

type Status struct {

	message string
	tStamp time.Time
}

type Queue struct {

	chJobs chan Job
	jobs []*Job
	ttl int64
}

const (

	Queued     = "Queued"
	Complete   = "Complete"
	InProgress = "In Progress"
)

func main() {

	var queuedJobs []Queue
	var cnt int64

	for cnt=0;cnt<100;cnt++ {

		queuedJobs = append(queuedJobs, Queue{})
		go sequentialWorker(queuedJobs[cnt], cnt)
	}

	r := gin.Default()

	r.GET("/id/:id",
		func(c *gin.Context) {

			var count int64
			var chMessage chan int64
			var chStatus chan Status
			uuid, err := uuid.NewRandom()
			queueId, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil { c.AbortWithError(400, err) }


			job := Job{ id: uuid.String(), message: chMessage,  status: chStatus }
			queuedJobs[queueId].chJobs <- job
			queuedJobs[queueId].jobs = append(queuedJobs[queueId].jobs, &job)
			count = <-chMessage

			c.JSON(200, gin.H{"job_id": job.id, "queue_id": queueId, "count": count})
		})

	r.GET("/status/queue/:queue_id/job/:job_id",
		func(c *gin.Context) {

			var ok bool
			var output string
			var status Status
			var job *Job

			queueId, err := strconv.ParseInt(c.Param("queue_id"), 10, 64)
			jobId := c.Param("job_id")
			if err != nil { c.AbortWithError(400, err) }

			for _, job = range queuedJobs[queueId].jobs {

				if jobId == job.id {

					select {
						case status, ok = <-job.status: if ok { break }
						default: {
							status.message = Queued
							status.tStamp = time.Now()
							break
						}
					}
					status = <-job.status
				}
			}

			if status.message != "" {
				output = fmt.Sprintf("%s Queue #%d is currently processing job %s with status: %s", status.tStamp.String(), queueId, jobId, status)
			} else {
				output = fmt.Sprintf("Invalid Queue Id or Job Id")
			}

			c.JSON(200, gin.H{"message": output})
		})

	r.Run()
}

func sequentialWorker(q Queue, queueId int64) {

	var count int64
	var status Status
	fName := "./keep-count-" + strconv.FormatInt(queueId, 10) + ".txt"
	_, statErr := os.Stat(fName)
	f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if os.IsNotExist(statErr) {
		f.Write([]byte(fmt.Sprintf("%d\n", 0)))
		return
	}

	if err != nil { log.Fatal(err) }

	for job := range q.chJobs {

		status = Status{InProgress, time.Now() }
		job.status <- status
		_, err = fmt.Fscanf(f, "%d\n", &count)
		count++
		if err != nil { log.Println(err) }
		_ = f.Truncate(0)
		_, err = f.Seek(0, 0)
		f.Write([]byte(fmt.Sprintf("%d\n", count)))
		job.message<- count
		status = Status{Complete, time.Now() }
		job.status<- status
		job.complete = true
	}
}