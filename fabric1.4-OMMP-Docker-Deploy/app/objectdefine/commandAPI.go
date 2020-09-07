package objectdefine

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//Run 工具执行命令
func (call *CommandList) Run() ([]string, error) {
	logs := make([]string, 0, len(call.Call))
	for _, one := range call.Call {
		opt := exec.Command(one.Exec, one.Args...)
		opt.Env = append(os.Environ(), one.Environment...)
		opt.Dir = one.Dir
		output, err := opt.CombinedOutput()
		if len(output) > 0 {
			lines := ConsoleStringToString(output)
			pairs := strings.Split(lines, "\n")
			logs = append(logs, pairs...)
		}

		if err != nil {
			err = errors.WithMessage(err, "Run "+one.Exec+" Process Error")
			return logs, err
		}
	}
	return logs, nil
}

//ConsoleStringToString 字符串转换
func ConsoleStringToString(buff []byte) string {
	return Utf16ToString(buff)
}

//Utf16ToString 格式转换
func Utf16ToString(input []byte) string {
	//win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	// utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(input), simplifiedchinese.GBK.NewDecoder())

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		return ""
	}
	return string(decoded)
}
