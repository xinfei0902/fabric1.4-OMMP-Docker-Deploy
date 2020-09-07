package dtask

import "deploy-server/app/objectdefine"

//TaskHandleMaker 后面匿名函数别名 用来执行任务
type TaskHandleMaker func(name string) objectdefine.RunTaskHandle

//MakeTaskPair 把每一步名称和接口对应 具体到每个任务中的每一步 例如MakeTaskPair("generalCreateChannelWriteDB",daction.MakeStepGeneralInsertCreateChannelIndent)
func MakeTaskPair(name string, maker TaskHandleMaker) objectdefine.SubTaskPair {
	return objectdefine.SubTaskPair{
		Name:   name,
		Handle: maker(name),
	}
}

type taskContainerInterface interface {
	GetStatus() *objectdefine.IndentStatus
	GetIndent() *objectdefine.Indent
}

type taskContainer struct {
	Status *objectdefine.IndentStatus
	Indent *objectdefine.Indent
}

func (container *taskContainer) GetStatus() *objectdefine.IndentStatus {
	return container.Status
}

func (container *taskContainer) GetIndent() *objectdefine.Indent {
	return container.Indent
}
