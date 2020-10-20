package dconfig

type flagsPair struct {
	Name  string
	Short string
	Value interface{}
	Usage string
}

type valuePair struct {
	Default interface{}
	File    interface{}
	Command interface{}
}

var (
	globalContainerWhole map[string]map[string]valuePair
	globalFlagsWhole     map[string]map[string]flagsPair

	globalValues  map[string]interface{}
	globalFileKey string
)

func initValues() {
	globalContainerWhole = make(map[string]map[string]valuePair)
	globalFlagsWhole = make(map[string]map[string]flagsPair)
	globalValues = make(map[string]interface{})
}

func (p *valuePair) GetInterface() (ret interface{}) {
	ret = GetInterface(p.Command)
	if ret != nil {
		return
	}
	ret = GetInterface(p.File)
	if ret != nil {
		return
	}
	ret = GetInterface(p.Default)
	return
}

func (p *valuePair) GetStringValue() (ret string) {
	ret = GetString(p.Command, "")
	if len(ret) > 0 {
		return ret
	}
	ret = GetString(p.File, "")
	if len(ret) > 0 {
		return ret
	}
	ret = GetString(p.Default, "")
	return
}

func (p *valuePair) GetIntValue() (ret int) {
	ret = GetInt(p.Command)
	if ret != 0 {
		return ret
	}
	ret = GetInt(p.File)
	if ret != 0 {
		return ret
	}
	ret = GetInt(p.Default)
	return
}

func (p *valuePair) GetBoolValue() (ret bool) {
	ret = GetBool(p.Command)
	if ret {
		return ret
	}
	ret = GetBool(p.File)
	if ret {
		return ret
	}
	ret = GetBool(p.Default)
	return
}

func GetString(value interface{}, emptyReplace string) string {
	if value == nil {
		return emptyReplace
	}
	switch value.(type) {
	case string:
		if len(value.(string)) == 0 {
			break
		}
		return value.(string)
	case *string:
		p := value.(*string)
		if p == nil || len(*p) == 0 {
			break
		}
		return *p
	}
	return emptyReplace
}

func GetInt(value interface{}) int {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case int:
		return value.(int)
	case *int:
		return *value.(*int)
	case float64, *float64:
		return int(GetFloat(value))
	}
	return 0
}

func GetFloat(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case float64:
		return value.(float64)
	case *float64:
		return *value.(*float64)
	}
	return 0
}

func GetBool(value interface{}) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case bool:
		return value.(bool)
	case *bool:
		return *value.(*bool)
	}
	return false
}

func GetInterface(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case bool, *bool:
		ret := GetBool(value)
		if ret {
			return ret
		}
	case int, *int:
		ret := GetInt(value)
		if ret != 0 {
			return ret
		}
	case float64, *float64:
		ret := GetFloat(value)
		if ret != 0 {
			return ret
		}
	case string, *string:
		ret := GetString(value, "")
		if len(ret) > 0 {
			return ret
		}
	default:
		panic("add code")
	}
	return nil
}
