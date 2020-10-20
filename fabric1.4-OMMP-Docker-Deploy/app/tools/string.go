package tools

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

//STDString 去掉空白 转小写
func STDString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

//TryParseStringToObj 当接口完成之后 有数据返回这里做处理
func TryParseStringToObj(input []byte) (ret interface{}, err error) {
	if len(input) == 0 {
		ret = string(input)
		return
	}
	switch input[0] {
	case '{':
		obj := make(map[string]interface{})
		err = json.Unmarshal(input, &obj)
		if err != nil {
			break
		}
		ret = obj
	case '[':
		obj := make([]interface{}, 0, 1)
		err = json.Unmarshal(input, &obj)
		if err != nil {
			break
		}
		ret = obj
	default:
		ret = string(input)
	}
	return
}

//IsAlpha 用来检测 小写和数字
func IsAlpha(input string) bool {
	for _, c := range input {
		if 'a' <= c && 'z' >= c {
			continue
		}
		if 'A' <= c && 'Z' >= c {
			continue
		}
		if '0' <= c && '9' >= c {
			continue
		}
		return false
	}
	return true
}

//InitialsToUpper 首字母检测是否大写 不是改成大写
func InitialsToUpper(s string) string {
	isASCII := true

	c := s[0]
	if c >= utf8.RuneSelf {
		isASCII = false
	}
	var by []rune
	if isASCII { // optimize for ASCII-only strings.
		for i, b := range s {
			if i == 0 {
				c := b
				if c >= 'a' && c <= 'z' {
					c -= 'a' - 'A'
				}
				by = append(by, c)
			} else {
				by = append(by, b)
			}
		}
		return string(by)
	}
	return s
}

//ReplaceTemplateBuff  用来检测 脚本文件$[]标识 替换整体内容
func ReplaceTemplateBuff(buff []byte, mapList map[string]string, sep []byte, parameterMark, openingBrace, closingBrace byte) []byte {
	if sep == nil {
		sep = []byte{}
	}
	key := make([]byte, 0, 50)

	copyBuff := make([]byte, 0, len(buff)+256)

	step := 0

	for _, c := range buff {

		switch step {
		case 0:
			if c == parameterMark {
				step = 1
				continue
			}
			copyBuff = append(copyBuff, c)
		case 1:
			if c == openingBrace {
				step = 2

				continue
			}
			step = 0
			copyBuff = append(copyBuff, parameterMark)
			copyBuff = append(copyBuff, c)
		case 2:
			if c == closingBrace {
				value, ok := mapList[string(key)]

				copyBuff = append(copyBuff, sep...)
				if ok && len(value) > 0 {
					copyBuff = append(copyBuff, []byte(value)...)
				}
				copyBuff = append(copyBuff, sep...)

				key = make([]byte, 0, 50)
				step = 0
				continue
			}

			key = append(key, c)
		}

	}

	return copyBuff
}
