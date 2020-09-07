package dtask

import (
	"deploy-server/app/objectdefine"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

//MakeGeneralStartStep 检测是否存在目录
func MakeGeneralStartStep() objectdefine.WholeTaskHandle {
	return func(task *objectdefine.TaskType) error {
		// check folder
		fmt.Println("start os.start")
		info, err := os.Stat(task.Indent.BaseOutput)
		if err == nil && info != nil {
			return errors.New("Files is exist")
		}
		// build folder
		err = os.MkdirAll(task.Indent.BaseOutput, 0777)
		if err != nil {
			return errors.WithMessage(err, "Build folder error")
		}

		return nil
	}
}

//MakeGeneralEndStep 保存订单以及状态
func MakeGeneralEndStep() objectdefine.WholeTaskHandle {
	return func(task *objectdefine.TaskType) error {
		err := task.SaveStatus()
		if err != nil {
			return errors.WithMessage(err, "Write status error")
		}
		return nil
	}
}
