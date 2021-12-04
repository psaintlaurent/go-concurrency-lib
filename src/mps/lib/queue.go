package lib

type Queue struct {
	ChJobs chan *Job
	Jobs   []*Job
	Ttl    int64
}