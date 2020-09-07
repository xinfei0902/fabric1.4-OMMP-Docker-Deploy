package objectdefine

//CommandList 工具: 生产证书 创世区块工具命令列表
type CommandList struct {
	OS string `json:"os,omitempty"`

	Call []ProcessPair `json:"call,omitempyt"`
}

//ProcessPair 执行列表里面参数信息等
type ProcessPair struct {
	Exec        string   `json:"exec,omitempty"`
	Args        []string `json:"args,omitempty"`
	Dir         string   `json:"dir,omitempty"`
	Environment []string `json:"env,omitempty"`
}
