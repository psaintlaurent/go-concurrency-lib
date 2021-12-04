package lib

type Job struct {
	Id       string
	Fn       chan int64
	Message  chan int64
	Status   chan Status
	Complete bool
}
