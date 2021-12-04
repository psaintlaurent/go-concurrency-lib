package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"os"
	l "saint-laurent.us/m/v2/src/mps/lib"
	"strconv"
	"time"
)


/*
	This takes advantage of a FIFO worker queue pattern to ensure parallel ordered execution.
*/

const (
	Queued      = "Queued"
	Complete    = "Complete"
	InProgress  = "In Progress"
	QueueLength = 100
)

func main() {

	var queuedJobs []l.Queue
	var cnt int64

	for cnt = 0; cnt < QueueLength; cnt++ {

		queuedJobs = append(queuedJobs, l.Queue{ChJobs: make(chan *l.Job, QueueLength)})
		go sequentialWorker(queuedJobs[cnt], cnt)
	}

	r := gin.Default()

	r.GET("/id/:id",
		func(c *gin.Context) {

			var count int64

			vuuid, err := uuid.NewRandom()
			queueId, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil || queueId > QueueLength-1 {
				c.AbortWithError(400, err)
			}

			job := l.Job{Id: vuuid.String(), Message: make(chan int64, 3), Status: make(chan l.Status, 3)}
			queuedJobs[queueId].Jobs = append(queuedJobs[queueId].Jobs, &job)
			queuedJobs[queueId].ChJobs<- &job
			count = <-job.Message

			c.JSON(200, gin.H{"job_id": job.Id, "queue_id": queueId, "count": count})
		})

	r.GET("/status/queue/:queue_id/job/:job_id",
		func(c *gin.Context) {

			var ok bool
			var output string
			var status l.Status
			var job *l.Job

			queueId, err := strconv.ParseInt(c.Param("queue_id"), 10, 64)
			jobId := c.Param("job_id")
			if err != nil {
				c.AbortWithError(400, err)
			}

			for _, job = range queuedJobs[queueId].Jobs {

				if jobId == job.Id {

					select {
					case status, ok = <-job.Status:
						if ok {
							break
						}
					default:
						{
							status.Message = Queued
							status.TStamp = time.Now()
							break
						}
					}
					status = <-job.Status
				}
			}

			if status.Message != "" {
				output = fmt.Sprintf("%s Queue #%d is currently processing job %s with status: %s", status.TStamp.String(), queueId, jobId, status)
			} else {
				output = fmt.Sprintf("Invalid Queue Id or Job Id")
			}

			c.JSON(200, gin.H{"message": output})
		})

	err := r.Run()
	if err != nil {
		return
	}
}

func sequentialWorker(q l.Queue, queueId int64) {

	var count int64
	var status l.Status
	var job *l.Job

	fName := "./keep-count-" + strconv.FormatInt(queueId, 10) + ".txt"
	_, statErr := os.Stat(fName)
	f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if os.IsNotExist(statErr) {
		log.Printf("File %s does not exist\n", fName)
	} else {
		f.Write([]byte(fmt.Sprintf("%d\n", 0)))
	}

	if err != nil {
		log.Fatal(err)
	}

	for job = range q.ChJobs {

		status = l.Status{InProgress, time.Now()}

		select {
			case job.Status <- status:
				log.Printf("Status update sent to Queue# %d Job # %s\n", queueId, job.Id)
			default:
				log.Println("Failure to update status")
		}

		_, err = fmt.Fscanf(f, "%d\n", &count)
		count++
		if err != nil {
			log.Println(err)
		}
		_ = f.Truncate(0)
		_, err = f.Seek(0, 0)
		f.Write([]byte(fmt.Sprintf("%d\n", count)))

		job.Message<- count
		status = l.Status{Message: Complete, TStamp: time.Now()}
		job.Status<- status
		job.Complete = true
	}
}
