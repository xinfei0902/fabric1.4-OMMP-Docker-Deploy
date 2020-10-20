package objectdefine

import (
	"sync"
)

//SubTaskPair 订单和日志
type SubTaskPair struct {
	Name   string
	Handle RunTaskHandle
}

//TaskType 任务结构
type TaskType struct {
	ID string

	Indent *Indent
	Status *IndentStatus

	LogChan chan []*StepHistory
	Exit    chan struct{}

	StepName string

	LastError      error
	StartWaitGroup sync.WaitGroup

	StartHandle WholeTaskHandle
	Step        []SubTaskPair
	EndHandle   WholeTaskHandle
}

//WholeTaskHandle 匿名函数别名
type WholeTaskHandle func(*TaskType) error

//RunTaskHandle 匿名函数别名
type RunTaskHandle func(*Indent, *TaskNode) error

//TaskNode 任务节点
type TaskNode struct {
	Log chan<- []*StepHistory

	Exit <-chan struct{}

	LogStore []*StepHistory
}

//TaskNodeInterface
type TaskNodeInterface interface {
	AppendLog(logs ...*StepHistory) (bool, error)
	FlushAll() (bool, error)
	IsQuit() bool
}
