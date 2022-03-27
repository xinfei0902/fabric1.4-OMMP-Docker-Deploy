package daction


// import (

// 	"deploy-server/app/objectdefine"
// 	"deploy-server/app/dssh"

// 	"github.com/pkg/errors"
// )

// func MakeStepGeneralCheckChaincodeIsInstantiated(name string) objectdefine.RunTaskHandle {
// 	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
// 		//远端ssh
// 		isInstantied,err := dssh.CheckChaincodeIsInstantied(general)
// 		if err != nil{
// 			return err
// 		}
// 		//执行命令

// 		//获取输出文件

// 		//解析文件内容判断是否实例化
// 	}
// }