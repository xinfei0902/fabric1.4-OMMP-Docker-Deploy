package dcmd

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dconfig"
	"deploy-server/app/dlog"
	"deploy-server/app/dmysql"
	"deploy-server/app/dstore"
	"deploy-server/app/dtask"
	"deploy-server/app/tools"
	"io/ioutil"

	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	keyFlagsDebug         = "debug"
	keyFlagsHelp          = "help"
	keyFlagsConfig        = "config"
	keyFlagsLog           = "log"
	keyFlagsLogLevel      = "loglevel"
	keyFlagsLogTime       = "logtime"
	keyFlagsLogCount      = "logcount"
	keyFlagsMysqlUser     = "userName"
	keyFlagsMysqlPassWord = "passWord"
	keyFlagsMysqlIP       = "ip"
	keyFlagsMysqlPort     = "port"
	keyFlagsMysqlDBName   = "dbName"
)

//addBaseFlag 添加主要把 key value值写进globalContainer Map里
func addBaseFlag(name string) {
	dconfig.SetFileNameKey(name, "c", keyFlagsConfig, tools.STDPath("config.json"))
	dconfig.Register(name, "h", keyFlagsHelp, false, "Show this Help")
	dconfig.Register(name, "d", keyFlagsDebug, false, "Open Debug Model")
	// //	dconfig.Register(name, "l", KeyFlagsLog, dlog.GlobalConsoleMark, "Absolute path of log file")
	// dconfig.Register(name, "e", KeyFlagsLogLevel, "warning", "Output Log level: [debug, info, warning, error, fatal, panic]")
	// dconfig.Register(name, "", KeyFlagsLogTime, "24h", "Log split per time")
	// dconfig.Register(name, "", KeyFlagsLogCount, 7, "Log split files max count")
}

//isDebug 检测globalValues Map是否存在
func isDebug() bool {
	v, ok := dconfig.Get(keyFlagsDebug)
	if !ok {
		return false
	}
	return dconfig.GetBool(v)
}

//isHelp 检测globalValues Map是否存在
func isHelp() bool {
	v, ok := dconfig.Get(keyFlagsDebug)
	if !ok {
		return false
	}
	return dconfig.GetBool(v)
}

//开始初始化日志
func startLog() error {
	var logfilePath, logLevel string
	var logtime time.Duration
	var logcount int
	logLevel = "warning"
	//获取config.json文件 "log"=内容
	v, _ := dconfig.Get(keyFlagsLog)
	//如果v 为空 那么默认时第二个参数的值
	logfilePath = dconfig.GetString(v, dlog.GlobalConsoleMark)

	v, _ = dconfig.Get(keyFlagsLogLevel)
	logLevel = dconfig.GetString(v, "warning")

	v, _ = dconfig.Get(keyFlagsLogTime)
	logtime, _ = time.ParseDuration(dconfig.GetString(v, "24h"))

	v, _ = dconfig.Get(keyFlagsLogTime)
	logcount = dconfig.GetInt(v)

	return dlog.InitLog(logfilePath, logLevel,
		logfilePath != dlog.GlobalConsoleMark,
		logtime,
		uint64(logcount),
	)
}

//decorateRunE 在服务运行之前做一些初始化准备等等
func decorateRunE(name string, cb func(cmd *cobra.Command, args []string) (err error)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := dconfig.ReadFileAndHoldCommand(name, cmd)
		if err != nil {
			fmt.Println("Quit error: ", err)
			return
		}
		if isHelp() {
			cmd.Help()
			return
		}

		err = startLog()
		if err != nil {
			fmt.Println("Log error: ", err)
			return
		}
		err = dmysql.StartDB()
		if err != nil {
			fmt.Println("db error: ", err)
			return
		}
		err = cb(cmd, args)
		if err != nil {
			fmt.Println("Quit error: ", err)
		}
	}
}

//makeCommand  name=server
func makeCommand(name, short, long string,
	configCB func(name string) error,
	cb func(cmd *cobra.Command, args []string) (err error)) *cobra.Command {
	addBaseFlag(name)
	ret := &cobra.Command{
		Use:   name,
		Short: short,
		Long:  long,

		Run: decorateRunE(name, cb),
	}
	//调用上一级
	err := configCB(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	dconfig.AddFlags(name, ret)

	return ret
}

//addLocalConfigFlags 把 key=local vlaue=./local.josn 写入全局变量globalContainer
func addLocalConfigFlags(name string) error {
	err := dconfig.Register(name, "j", "local", "./local.json", "Local config filename")
	return err
}

//loadLocalConfig  加载local.json文件信息 加载缓存信息 初始化任务管理
func loadLocalConfig() (err error) {
	name := dconfig.GetStringByKey("local")
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		err = errors.WithMessage(err, "Read local config file error")
		return
	}

	one, err := dstore.NewLocalTypeFromJSON(buff)
	if err != nil {
		return
	}

	err = dcache.InitLocalConfig(one)
	if err != nil {
		err = errors.WithMessage(err, "Init local object failed")
		return
	}

	err = dtask.InitTaskManager()

	if err != nil {
		return errors.WithMessage(err, "Init task manager error")
	}

	buff, err = dstore.LocalObejctIntoJSONBuff(one)
	if err == nil && len(buff) > 0 {
		err = ioutil.WriteFile(name, buff, 0644)
	}
	if err != nil {
		dlog.Warn("Reflash local config files err: " + err.Error())
	}

	return nil
}
