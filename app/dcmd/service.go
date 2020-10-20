package dcmd

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dservice"
	"deploy-server/app/web"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//initService 初始
func initService() *cobra.Command {
	ret := makeCommand("service", "service", `Start service for deploy tasks`,
		serviceFlags,
		enterService,
	)
	return ret
}

//serviceFlags 设置k-v导入全局变量
func serviceFlags(name string) (err error) {
	err = addLocalConfigFlags(name)
	if err != nil {
		return err
	}
	err = addWebFlag(name)
	if err != nil {
		return err
	}

	return
}

//enterService 加载启动服务依赖项
func enterService(cmd *cobra.Command, args []string) (err error) {

	var debug, db bool
	debug = isDebug()
	err = loadLocalConfig()
	if err != nil {
		return
	}
	err = dservice.RegisterWebAPI(debug, db)
	if err != nil {
		err = errors.WithMessage(err, "Register Web APIs error")
		return
	}
	err = StartWeb()
	return err
}

//addWebFlag 添加启动默认值 把key 和默认value写入全局变量globalContainer  服务启动需要端口等
func addWebFlag(name string) error {
	GlobalDefaultAddress := ":7000"

	err := dconfig.Register(name, "", "cert", "", "Cert File name for TLS")
	if err != nil {
		return err
	}
	err = dconfig.Register(name, "", "key", "", "Key File name for TLS")
	err = dconfig.Register(name, "a", "address", GlobalDefaultAddress, "Bind Service on this Address. Default: "+GlobalDefaultAddress)
	return err
}

//StartWeb 获取配置启动服务
func StartWeb() (err error) {
	cert := dconfig.GetStringByKey("cert")
	key := dconfig.GetStringByKey("key")
	address := dconfig.GetStringByKey("address")
	if len(cert) > 0 && len(key) > 0 {
		fmt.Println("Start TLS service on", `"`+address+`"`)
		err = web.StartTLSService(address, cert, key)
	} else {
		fmt.Println("Start service on", `"`+address+`"`)
		err = web.StartService(address)
	}
	return
}
