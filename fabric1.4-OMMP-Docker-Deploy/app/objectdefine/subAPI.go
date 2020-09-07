package objectdefine

type InterfaceNodeKind interface {
	GetIP() string
	GetDomain() string

	GetFolder() string

	GetSetupPath() string
}

func (opt *KafkaType) GetIP() string     { return opt.IP }
func (opt *ZooKeeperType) GetIP() string { return opt.IP }
func (opt *PeerType) GetIP() string      { return opt.IP }

//func (opt *BaaSType) GetIP() string      { return opt.IP }
func (opt *CAType) GetIP() string    { return opt.IP }
func (opt *OrderType) GetIP() string { return opt.IP }

func (opt *KafkaType) GetFolder() string     { return opt.Folder }
func (opt *ZooKeeperType) GetFolder() string { return opt.Folder }
func (opt *PeerType) GetFolder() string      { return opt.Folder }

//func (opt *BaaSType) GetFolder() string      { return opt.Folder }
// func (opt *CAType) GetFolder() string    { return opt.Folder }
// func (opt *OrderType) GetFolder() string { return opt.Folder }

// func (opt *KafkaType) GetSetupPath() string     { return opt.SetupPath }
// func (opt *ZooKeeperType) GetSetupPath() string { return opt.SetupPath }
// func (opt *PeerType) GetSetupPath() string      { return opt.SetupPath }

// //func (opt *BaaSType) GetSetupPath() string      { return opt.SetupPath }
// func (opt *CAType) GetSetupPath() string    { return opt.SetupPath }
// func (opt *OrderType) GetSetupPath() string { return opt.SetupPath }

func (opt *KafkaType) GetDomain() string     { return opt.Domain }
func (opt *ZooKeeperType) GetDomain() string { return opt.Domain }
func (opt *PeerType) GetDomain() string      { return opt.Domain }

//func (opt *BaaSType) GetDomain() string      { return opt.Domain }
func (opt *CAType) GetDomain() string    { return opt.Domain }
func (opt *OrderType) GetDomain() string { return opt.Domain }
