package dtask

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dlog"
	"deploy-server/app/objectdefine"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

//TaskStartHandle 任务开始
type TaskStartHandle func(input *objectdefine.Indent) (*objectdefine.IndentStatus, error)

//const 变量
const (
	TaskGeneralComplete         = "generalcompletedeploy"
	TaskGeneralCompleteDelete         = "generalcompletedeletedeploy"
	TaskGeneralCreateChannel    = "generalcreatechannel"
	TaskGeneralAddOrg           = "generaladdorg"
	TaskGeneralDeleteOrg        = "generaldeleteorg"
	TaskGeneralAddPeer          = "generaladdpeer"
	TaskGeneralDeletePeer       = "generaldeletepeer"
	TaskGeneralDisablePeer      = "generaldisablepeer"
	TaskGeneralEnablePeer       = "generalenablepeer"
	TaskGeneralModiflyPeer      = "generalmodiflypeer"
	TaskGeneralChainCodeUpload  = "generalchaincodeupload"
	TaskGeneralChainCodeBulid  = "generalchaincodebulid"
	TaskGeneralChainCodeBulidUpload = "generalchaincodebulidupload"
	TaskGeneralChainCodeList    = "generalchaincodelist"
	TaskGeneralChainCodeAdd     = "generalchaincodeadd"
	TaskGeneralChainCodeDelete  = "generalchaincodedelete"
	TaskGeneralChainCodeUpgrade = "generalchaincodeupgrade"
	TaskGeneralChainCodeDisable = "generalchaincodedisable"
	TaskGeneralChainCodeEnable  = "generalchaincodeenable"
	TaskGeneralAddServerInfo  = "generaladdserverinfo"
)

//MakePushTaskHandle 获取初始化的全局变量 任务管理器 然后开启任务
func MakePushTaskHandle(kind string) TaskStartHandle {
	return func(input *objectdefine.Indent) (*objectdefine.IndentStatus, error) {
		if nil == input {
			return nil, errors.New("Empty Params")
		}
		core := GetTaskManager()
		if true == core.IsIndentAlive(input.ID) {
			return nil, errors.New("Task [" + input.ID + "] is exist")
		}
		return core.StartTask(input, kind)
	}
}

//TaskManager 任务管理结构
type TaskManager struct {
	lockcore sync.RWMutex

	ActiveTask map[string]taskContainerInterface

	HistoryTask map[string]taskContainerInterface

	Quit chan string

	Exit chan struct{}
}

var globalTaskManager *TaskManager
var globalTaskManagerOnce sync.Once

//GetTaskManager 获取初始化之后任务管理 里已经加载历史任务 或者新的任务
func GetTaskManager() *TaskManager {
	return globalTaskManager
}

//newTaskManager new 一个新的任务管理结构 开辟空间
func newTaskManager() *TaskManager {
	return &TaskManager{
		ActiveTask:  make(map[string]taskContainerInterface),
		HistoryTask: make(map[string]taskContainerInterface),

		Quit: make(chan string),
		Exit: make(chan struct{}),
	}
}

//InitTaskManager 初始开启任务管理
func InitTaskManager() (err error) {
	one := newTaskManager()
	err = one.LoadFromFiles()
	if err != nil {
		return
	}
	globalTaskManagerOnce.Do(func() {
		//把初始化历史任务 赋值 开启线程
		globalTaskManager = one
		go globalTaskManager.MainWork()
	})
	return
}

//Lock 锁
func (manager *TaskManager) Lock() {
	manager.lockcore.Lock()
}

//Unlock 开锁
func (manager *TaskManager) Unlock() {
	manager.lockcore.Unlock()
}

//RLock 锁
func (manager *TaskManager) RLock() {
	manager.lockcore.RLock()
}

//RUnlock 开锁
func (manager *TaskManager) RUnlock() {
	manager.lockcore.RUnlock()
}

//LoadFromFiles 重启任务 需要重新加载原先的任务id 把任务放到历史任务中
func (manager *TaskManager) LoadFromFiles() error {

	pair, err := manager.getLocalTasks()
	if err != nil {
		return err
	}
	manager.setHistory(pair)

	return nil
}

//setHistory 把状态和参数存放任务管理历史任务中
func (manager *TaskManager) setHistory(input map[string]taskContainerInterface) {
	manager.Lock()
	defer manager.Unlock()
	manager.HistoryTask = input

}

//getLocalTasks 获取历史任务id 以及id执行过程中状态等
func (manager *TaskManager) getLocalTasks() (map[string]taskContainerInterface, error) {
	root := dcache.GetStoreRoot()
	namelist := make([]string, 0, 32)
	//变量所有存放以任务id为目录的路径
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		if info.IsDir() == false {
			return nil
		}
		_, id := filepath.Split(path)
		namelist = append(namelist, id)
		return filepath.SkipDir
	})

	if err != nil {
		return nil, err
	}
	ret := make(map[string]taskContainerInterface)

	for _, one := range namelist {
		if manager.IsIndentAlive(one) {
			continue
		}
		filename := filepath.Join(root, one, objectdefine.ConstHistoryTaskStatusFileName)
		status, err := GetLocalStatusByFile(filename)
		if err != nil {
			if err != nil {
				dlog.Warn(fmt.Sprintf("GetLocalStatusByFile [%s] error: %v", one, err))
			}
			continue
		}

		if len(status.ID) == 0 || len(status.Plains) == 0 || status.ID != one {
			continue
		}

		filename = filepath.Join(root, one, objectdefine.ConstHistoryTaskIndentFileName)

		indent, err := GetLocalIndentByFile(filename)
		if err != nil || len(indent.ID) == 0 || len(indent.Version) == 0 || indent.ID != one {
			indent = nil
		}

		ret[one] = &taskContainer{
			Indent: indent,
			Status: status,
		}

	}
	return ret, nil
}

//IsIndentAlive 检测任务id 是否正在活跃（使用）中
func (manager *TaskManager) IsIndentAlive(id string) bool {
	manager.RLock()
	defer manager.RUnlock()
	_, ok := manager.ActiveTask[id]
	return ok
}

//RegisterActive 把任务id 注册到活跃中
func (manager *TaskManager) RegisterActive(id string, container taskContainerInterface) bool {
	if nil == container {
		return false
	}

	manager.Lock()
	defer manager.Unlock()
	_, ok := manager.ActiveTask[id]
	if ok {
		return false
	}

	manager.ActiveTask[id] = container
	return true
}

//RemoveActive 移除任务id
func (manager *TaskManager) RemoveActive(id string) {
	manager.Lock()
	defer manager.Unlock()
	_, ok := manager.ActiveTask[id]
	if !ok {
		return
	}

	delete(manager.ActiveTask, id)
	return
}

//判断是否有正在运行的任务id
// func (manager *TaskManager) isRunActive() bool {
// 	manager.RLock()
// 	defer manager.RUnlock()
// 	for id, _ := range manager.ActiveTask {
// 		if len(id) != 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

//GetLocalStatusByFile 读取任务id目录下面的status.json文件  此文件存放此任务执行过程中状态 以及要点日志
func GetLocalStatusByFile(name string) (*objectdefine.IndentStatus, error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if len(buff) == 0 {
		return nil, os.ErrInvalid
	}

	status := &objectdefine.IndentStatus{}
	err = json.Unmarshal(buff, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

//GetLocalIndentByFile 读取任务id目录下面的indent.json文件  此文件存放此任务执行传进的参数
//目前不使用
func GetLocalIndentByFile(name string) (*objectdefine.Indent, error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if len(buff) == 0 {
		return nil, os.ErrInvalid
	}

	ret := &objectdefine.Indent{}
	err = json.Unmarshal(buff, ret)
	if err != nil {
		return nil, err
	}

	parent := filepath.Dir(name)
	if filepath.IsAbs(parent) {
		ret.BaseOutput = parent
	}
	return ret, nil
}

//MainWork 不断检测任务id 状态 一旦执行完毕 由活跃转历史
func (manager *TaskManager) MainWork() {

	wait := time.Second * 5
Main:
	for {
		select {
		case id, ok := <-manager.Quit:
			if !ok {
				break Main
			}

			manager.MoveActiveToHistory(id)

		case _, ok := <-manager.Exit:
			if !ok {
				break Main
			}
		//超时处理
		case <-time.After(wait):
			manager.syncHistory()
		}
	}

	select {
	case _, ok := <-manager.Quit:
		if ok {
			close(manager.Quit)
		}
	default:
		close(manager.Quit)
	}
}

//MoveActiveToHistory 由活跃任务转成历史任务
func (manager *TaskManager) MoveActiveToHistory(id string) {
	manager.Lock()
	defer manager.Unlock()

	one, ok := manager.ActiveTask[id]
	if !ok {
		return
	}

	delete(manager.ActiveTask, id)

	_, ok = manager.HistoryTask[id]
	if !ok {
		manager.HistoryTask[id] = one
	}
}

//syncHistory 同步一下历史任务
func (manager *TaskManager) syncHistory() {
	pair, err := manager.getLocalTasks()
	if err != nil {
		dlog.Warn(fmt.Sprintf("Sync History Task from files error: %v", err))
		return
	}
	manager.setHistory(pair)
}

//StartTask 开始执行任务
func (manager *TaskManager) StartTask(indent *objectdefine.Indent, kind string) (*objectdefine.IndentStatus, error) {
	var taskOpt *objectdefine.TaskType
	switch kind {
	case TaskGeneralComplete:
		taskOpt = NewGeneralCompleteDeploy()
	case TaskGeneralCompleteDelete:
		taskOpt = NewGeneralCompleteDeleteDeploy()	
	case TaskGeneralCreateChannel:
		taskOpt = NewGeneralCreateChannel()
	case TaskGeneralAddOrg:
		taskOpt = NewGeneralAddOrgTask()
	case TaskGeneralDeleteOrg:
		taskOpt = NewGeneralDeleteOrgTask()
	case TaskGeneralAddPeer:
		taskOpt = NewGeneralAddPeerTask()
	case TaskGeneralDeletePeer:
		taskOpt = NewGeneralDeletePeerTask()
	case TaskGeneralDisablePeer:
		taskOpt = NewGeneralDisablePeerTask()
	case TaskGeneralEnablePeer:
		taskOpt = NewGeneralEnablePeerTask()
	case TaskGeneralModiflyPeer:
		taskOpt = NewGeneralModiflyPeerTask()
	case TaskGeneralChainCodeAdd:
		taskOpt = NewGeneralChaincCodeAdd()
	case TaskGeneralChainCodeDelete:
		taskOpt = NewGeneralChaincCodeDelete()
	case TaskGeneralChainCodeUpgrade:
		taskOpt = NewGeneralChainCodeUpgrade()
	case TaskGeneralChainCodeDisable:
		taskOpt = NewGeneralChainCodeDisable()
	case TaskGeneralChainCodeEnable:
		taskOpt = NewGeneralChainCodeEnable()
	case TaskGeneralChainCodeBulid:
		taskOpt = NewGeneralChainCodeBulid() 	
	case TaskGeneralAddServerInfo:
		taskOpt = NewGeneralAddServerInfo()
	default:
		return nil, errors.New("Unkown kind [" + kind + "]")
	}
	if taskOpt == nil {
		return nil, errors.New("Unsupported Task [" + kind + "]")
	}
	//先返回操作任务步骤
	one, err := taskOpt.Prepare(indent)
	if err != nil {
		return nil, errors.WithMessage(err, "Task ["+indent.ID+"] Prepare error")
	}
	if false == manager.RegisterActive(indent.ID, taskOpt) {
		return nil, errors.New("Same Task [" + indent.ID + "] is exist")
	}

	err = taskOpt.Run(manager.Quit)

	if err != nil {
		manager.RemoveActive(indent.ID)
		return nil, errors.WithMessage(err, "Run Task ["+indent.ID+"] is error")
	}
	return one, nil
}

//IsIndentAlive 检测任务id 是否正在活跃（使用）中
func IsIndentAlive(id string) bool {
	core := GetTaskManager()
	return core.IsIndentAlive(id)
}

//IsIndentIDExist 检测任务id 是否存在
func IsIndentIDExist(id string) bool {
	core := GetTaskManager()
	return core.IsIndentIDExist(id)
}

//IsIndentIDExist 检测任务id 是否存在
func (manager *TaskManager) IsIndentIDExist(id string) bool {
	manager.RLock()
	defer manager.RUnlock()
	_, ok := manager.ActiveTask[id]
	if ok {
		return true
	}
	_, ok = manager.HistoryTask[id]
	return ok
}

//GetTaskStatus 获取任务执行状态
func GetTaskStatus(id string) *objectdefine.IndentStatus {
	core := GetTaskManager()
	return core.GetTaskStatus(id)
}

//GetTaskStatus 获取任务执行状态
func (manager *TaskManager) GetTaskStatus(id string) *objectdefine.IndentStatus {
	manager.RLock()
	defer manager.RUnlock()
	one, ok := manager.ActiveTask[id]
	if ok && one != nil && one.GetStatus() != nil {
		ret := *one.GetStatus()
		return &ret
	}
	one, ok = manager.HistoryTask[id]
	if ok && one != nil && one.GetStatus() != nil {
		ret := *one.GetStatus()
		return &ret
	}
	return nil
}
