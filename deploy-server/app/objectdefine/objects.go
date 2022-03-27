package objectdefine

// BaseReponse for API
type BaseReponse struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

///

// ConsensusList for API
type ConsensusList struct {
	Consensus []string `json:"consensus,omitempty"`
}

///

// IndentID for API
type IndentID struct {
	ID string `json:"id"`
}

////

// IndentSecret for API
type IndentSecret struct {
	ID     string     `json:"id,omitempty"`
	//Secret []IPSecret `json:"secret,omitempty"`
}

// IPSecret for IndentSecret
type IndentServer struct {
	ServerName        string   `json:"servername,omitempty"`
	ServerDes        string `json:"serverdes,omitempty"`
	ServerExtIp         string   `json:"serverextip,omitempty"`
	ServerIntIp  string `json:"serverintip,omitempty"`
	ServerUser  string  `json:"serveruser,omitempty"`
	ServerPassword string  `json:"serverpassword,omitempty"`
	ServerNum int `json:"servernum,omitempty"`
	ServerStatus int  `json:"serverstatus,omitempty"`
}

//IndentHTTP : add new http request
type IndentHTTP struct {
	IP       string `json:"ip,omitempty"`
	HomePage string `json:"homepage,omitempty"`
	//URL  string `json:"url,omitempty"`
}

//RequestBody 请求信息
type RequestBody struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

//ReponseBody 响应
type ReponseBody struct {
	Code *CodeType `json:"code,omitempty"`
	Data string    `json:"data,omitempty"`
}

//CodeType  响应信息里面错误码已经信息
type CodeType struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

///

// IndentStatus for API
type IndentStatus struct {
	ID      string         `json:"id,omitempty"`
	Desc    string         `json:"desc,omitempty"`
	Plains  []StepPlains   `json:"plains,omitempty"`
	History []*StepHistory `json:"history,omitempty"`
}

// StepPlains for IndentStatus
type StepPlains struct {
	Name    string `json:"name,omitempty"`
	Success bool   `json:"success"`
	Done    bool   `json:"done"`
}

// StepHistory for IndentStatus
type StepHistory struct {
	ID    uint64   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Log   []string `json:"log,omitempty"`
	Error []string `json:"error,omitempty"`
}

///

// ChaincodeInstantiate for API
type ChaincodeInstantiate struct {
	ID        string                   `json:"id,omitempty"`
	ChainCode map[string]ChainCodeType `json:"chaincode,omitempty"`
}

///

// VersionType for API
type VersionType struct {
	Version     string                   `json:"version,omitempty"`
	VersionRoot string                   `json:"versionroot,omitempty"`
	BinRoot     string                   `json:"binroot,omitempty"`
	ChainCode   map[string]ChainCodeType `json:"chaincode,omitempty"`

	Build map[string]VersionResource `json:"build,omitempty"`

	Sub map[string]string `json:"sub,omitempty"`
}

// VersionResource for VersionType
type VersionResource struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	BinName string `json:"bin,omitempty"`
}
