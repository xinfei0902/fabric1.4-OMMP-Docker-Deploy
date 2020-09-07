package dservice

import (
	"deploy-server/app/dtask"
	"errors"
)

//CheckID 检测任务id 是否存在
func CheckID(id string, existErr bool) error {

	if len(id) == 0 {
		return errors.New("Empty buff for ID")
	}

	// if false == tools.IsAlpha(id) || len(id) > 128 {
	// 	return errors.New("ID format error")
	// }

	if dtask.IsIndentAlive(id) == true {
		return errors.New("ID is alive")
	}

	if existErr == dtask.IsIndentIDExist(id) {
		if existErr {
			return errors.New("ID [" + id + "] is exist")
		}
		return errors.New("ID [" + id + "] is not exist")
	}

	return nil
}
