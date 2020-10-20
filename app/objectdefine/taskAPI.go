package objectdefine

import (
	"deploy-server/app/dlog"

	"deploy-server/app/tools"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

//NewTaskNode  没啥可说的
func NewTaskNode(logChan chan<- []*StepHistory, exitChan <-chan struct{}) *TaskNode {
	return &TaskNode{
		Log:  logChan,
		Exit: exitChan,
	}
}

//AppendLog 日志
func (node *TaskNode) AppendLog(logs ...*StepHistory) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			// discard panic of write closed chan
			dlog.Warn(fmt.Sprintf("AppendLog chan closed: %v", err))
		}
	}()

	if len(node.LogStore) > 0 {
		node.LogStore = append(node.LogStore, logs...)
	} else {
		node.LogStore = logs
	}

	select {
	case node.Log <- node.LogStore:
		node.LogStore = nil
		return true, nil
	case _, ok := <-node.Exit:
		if false == ok {
			return false, errors.New("Exit Task")
		}
	case <-time.After(1 * time.Second):
		return false, nil
	}
	return false, nil
}

//FlushAll 置空
func (node *TaskNode) FlushAll() (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			// discard panic of write closed chan
			dlog.Warn(fmt.Sprintf("FlushAll chan closed: %v", err))
		}
	}()

	// block it
	if len(node.LogStore) > 0 {
		select {
		case node.Log <- node.LogStore:
			node.LogStore = nil
			return true, nil
		case _, ok := <-node.Exit:
			if false == ok {
				return false, errors.New("Exit Task")
			}
		}
	}

	return true, nil
}

//IsQuit 是否退出
func (node *TaskNode) IsQuit() bool {
	select {
	case _, ok := <-node.Exit:
		if false == ok {
			return true
		}
	default:
	}
	return false
}

//NewTaskType  new TaskType{} 建立一个新的任务
func NewTaskType(step []SubTaskPair, start, end WholeTaskHandle) *TaskType {
	return &TaskType{
		LogChan:     make(chan []*StepHistory),
		Exit:        make(chan struct{}),
		Step:        step,
		StartHandle: start,
		EndHandle:   end,
	}
}

func (task *TaskType) createNode() *TaskNode {
	return NewTaskNode(task.LogChan, task.Exit)
}

//createStatus 把获取的执行步骤名称 返回 success done 都默认为false 所以调用接口返回的都是false false
func (task *TaskType) createStatus() *IndentStatus {
	ret := &IndentStatus{
		ID:   task.Indent.ID,
		Desc: task.Indent.Desc,

		Plains:  make([]StepPlains, 0, len(task.Step)),
		History: make([]*StepHistory, 0, 30*len(task.Step)),
	}

	for _, one := range task.Step {
		ret.Plains = append(ret.Plains, StepPlains{Name: one.Name})
	}
	return ret
}

func (task *TaskType) isEmpty() bool {
	return len(task.ID) == 0 || nil == task.Indent
}

//GetStatus 获取任务状态
func (task *TaskType) GetStatus() *IndentStatus {
	if task.isEmpty() {
		return &IndentStatus{
			ID:   "",
			Desc: "(no task)",
		}
	}

	return task.Status
}

//GetIndent 获取任务订单信息
func (task *TaskType) GetIndent() *Indent {
	return task.Indent
}

//Run 开启线程 实时返回状态日志
func (task *TaskType) Run(quitMark chan<- string) error {
	nodeWriter := task.createNode()

	// start mainWork
	task.StartWaitGroup.Add(1)
	go task.mainWork(quitMark)
	task.StartWaitGroup.Wait()

	if task.LastError != nil {
		return errors.WithMessage(task.LastError, "Start main task error")
	}

	// start sub work
	task.StartWaitGroup.Add(1)
	go task.subWorkStepByStep(nodeWriter)
	task.StartWaitGroup.Wait()

	if task.LastError != nil {
		return errors.WithMessage(task.LastError, "Run sub tasks error")
	}

	return nil
}

//Prepare 准备返回结果 把任务步骤返回
func (task *TaskType) Prepare(input *Indent) (*IndentStatus, error) {
	if false == task.isEmpty() {
		return nil, errors.New("Task is empty")
	}

	if len(input.BaseOutput) == 0 {
		return nil, errors.New("Outputpath is empty")
	}

	task.ID = input.ID
	task.Indent = input

	// create object
	task.Status = task.createStatus()
	return task.Status, nil
}

//SaveStatus 保存状态
func (task *TaskType) SaveStatus() error {
	// mainWork write log
	buff, err := json.Marshal(task.Status)
	if err != nil {
		return errors.WithMessage(err, "Task ["+task.ID+"] done. Parse logs into json error")
	}

	if len(task.Indent.BaseOutput) == 0 {
		return errors.New("Empty task folder")
	}

	name := filepath.Join(task.Indent.BaseOutput, ConstHistoryTaskStatusFileName)

	err = ioutil.WriteFile(name, buff, 0644)
	if err != nil {
		return errors.WithMessage(err, "Task ["+task.ID+"] done. Save logs into file error")
	}
	return nil
}

//SaveIndent 保存订单 接口没有用到
func (task *TaskType) SaveIndent() error {
	//tmp := task.Indent.Secret
	//task.Indent.Secret = nil
	// mainWork write log
	buff, err := json.Marshal(task.Indent)
	if err != nil {
		return errors.WithMessage(err, "Parse indent ["+task.ID+"] into json error")
	}

	//task.Indent.Secret = tmp

	if len(task.Indent.BaseOutput) == 0 {
		return errors.New("Empty task folder")
	}

	name := filepath.Join(task.Indent.BaseOutput, ConstHistoryTaskIndentFileName)

	err = ioutil.WriteFile(name, buff, 0644)
	if err != nil {
		return errors.WithMessage(err, "Save logs ["+task.ID+"] into file error")
	}
	return nil
}

//LoadStatus 加载状态文件
func (task *TaskType) LoadStatus() error {
	if task.Indent == nil || len(task.Indent.BaseOutput) == 0 {
		return errors.New("No output path")
	}
	name := filepath.Join(task.Indent.BaseOutput, ConstHistoryTaskStatusFileName)

	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return errors.WithMessage(err, "Read logs file error")
	}

	one := &IndentStatus{}
	err = json.Unmarshal(buff, one)
	if err != nil {
		return errors.WithMessage(err, "Parse indent file error")
	}

	if one.ID != task.Indent.ID {
		return errors.WithMessage(err, "Indent ["+task.ID+"] & File ["+one.ID+"] is not match")
	}

	task.Status = one
	return nil
}

//LoadIndent 加载订单 没有用到
func (task *TaskType) LoadIndent() error {
	if task.Indent == nil || len(task.Indent.BaseOutput) == 0 {
		return errors.New("No output path")
	}

	//secret := task.Indent.Secret
	output := task.Indent.BaseOutput

	name := filepath.Join(task.Indent.BaseOutput, ConstHistoryTaskIndentFileName)

	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return errors.WithMessage(err, "Read indent file error")
	}

	one := &Indent{}
	err = json.Unmarshal(buff, one)
	if err != nil {
		return errors.WithMessage(err, "Parse indent file error")
	}

	if one.ID != task.Indent.ID {
		return errors.WithMessage(err, "Indent ["+task.ID+"] & File ["+one.ID+"] is not match")
	}

	//one.Secret = secret
	one.BaseOutput = output
	task.Indent = one

	return nil
}

func (task *TaskType) mainWork(quitMark chan<- string) {
	//判断任务输出目录是否生成成功
	if task.StartHandle != nil {
		err := task.StartHandle(task)
		if err != nil {
			task.LastError = err
			task.StartWaitGroup.Done()
			return
		}
	}
	task.StartWaitGroup.Done()

	logOffset := uint64(1)
Main:
	for {
		select {
		case _, ok := <-task.Exit:
			if false == ok {
				break Main
			}
		case lines, ok := <-task.LogChan:
			if false == ok {
				break Main
			}
			if len(lines) == 0 {
				continue
			}
			for i := range lines {
				lines[i].ID = logOffset
				logOffset++
				task.Status.History = append(task.Status.History, lines[i])
			}
		}
	}

	for i := range task.Status.Plains {
		task.Status.Plains[i].Done = true
	}

LogMain:
	for {
		select {
		case lines, ok := <-task.LogChan:
			if false == ok {
				break LogMain
			}
			if len(lines) == 0 {
				continue
			}
			for i := range lines {
				lines[i].ID = logOffset
				logOffset++
				task.Status.History = append(task.Status.History, lines[i])
			}
		default:
			close(task.LogChan)
		}
	}

	// mainWork get exit
	if task.EndHandle != nil {
		err := task.EndHandle(task)
		if err != nil {
			task.Status.History = append(task.Status.History, &StepHistory{
				ID:    logOffset,
				Name:  "close",
				Error: []string{fmt.Sprintf("Close task error: %v", err)},
			})
		}
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				// log
				dlog.Warn(fmt.Sprintf("quitMark chan closed: %v", err))
			}
		}()

		quitMark <- task.ID
	}()

	// mainWork quit
	return
}

func (task *TaskType) subWorkStepByStep(output *TaskNode) {
	//
	task.StartWaitGroup.Done()
	// sub work quit / set exit

	for i, pair := range task.Step {

		if output.IsQuit() {
			break
		}
		//接受 任务每一步骤返回的日志 如果错误截至 没有错误把success置true
		err := pair.Handle(task.Indent, output)
		task.Status.Plains[i].Done = true
		if err != nil {
			output.AppendLog(&StepHistory{
				Name:  pair.Name,
				Error: []string{"Done by error: " + err.Error(), "Task is done"},
			})
			break
		}

		output.AppendLog(&StepHistory{
			Name: pair.Name,
			Log:  []string{"Task is success"},
		})
		task.Status.Plains[i].Success = true
	}

	// safe close
	tools.SafeCloseQuit(task.Exit)
}

//SaveGeneralIndent 不用了
func (task *TaskType) SaveGeneralIndent() error {

	return nil
}

//AppendDNS 用来保存IP和域名键值对的
func AppendDNS(container map[string][]string, ip, domain string) {
	pair, ok := container[ip]
	if !ok {
		pair = make([]string, 0, 16)
	}
	for _, one := range pair {
		if one == domain {
			return
		}
	}
	pair = append(pair, domain)

	container[ip] = pair
}

//DeleteDNS 用来删除IP和域名键值对的
func DeleteDNS(container map[string][]string, ip, domain string) {
	pair, _ := container[ip]
	// if !ok {
	// 	pair = make([]string, 0, 16)
	// }
	temp := make([]string, 0, 16)
	for _, one := range pair {
		if one != domain {
			temp = append(temp, one)
		}
	}
	pair = temp
	container[ip] = pair
}

//AppendPort 添加端口
func AppendPort(container map[string]map[int]bool, ip string, port *int) bool {
	list, ok := container[ip]
	if !ok {
		list = make(map[int]bool, 16)
	}

	list, ok = CheckPort(list, port)
	if ok {
		container[ip] = list
	}
	return ok
}

//CheckPort 检测 自加
func CheckPort(port map[int]bool, one *int) (map[int]bool, bool) {
	v := port[*one]
	for v {
		*one += 1000
		if *one > 65000 {
			return nil, false
		}

		v = port[*one]
	}

	port[*one] = true
	return port, true
}

//DeletePort 删除端口
func DeletePort(container map[string]map[int]bool, ip string, port *int) {
	list, _ := container[ip]
	// if !ok {
	// 	list = make(map[int]bool, 16)
	// }
	delete(list, *port)
	container[ip] = list
}

//------------------------------------------------------
