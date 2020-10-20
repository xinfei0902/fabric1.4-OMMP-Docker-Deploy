package dmysql

import (
	"database/sql"
	"deploy-server/app/dcache"
	"deploy-server/app/dconfig"
	"deploy-server/app/objectdefine"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	userName string
	passWord string
	ip       string
	port     string
	dbName   string
	db       *sql.DB
)

type channelStruct struct {
	ChannelName   string `json:"channelname,omitempty"`
	OrgNum        int    `json:"orgnum,omitempty"`
	PeerNum       int    `json:"peernum,omitempty"`
	CCNum         int    `json:"ccnum,omitempty"`
	ChannelStatus int    `json:"status,omitempty"`
}

type OrgStruct struct {
	OrgName     string `json:"orgname,omitempty"`
	OrgDomain   string `json:"orgdomain,omitempty"`
	PeerNum     int64  `json:"peernum,omitempty"`
	OrgStatus   int    `json:"status,omitempty"`
	ChannelName string `json:"channelname,omitempty"`
}

type PeerStruct struct {
	IP          string `json:"ip,omitempty"`
	PeerName    string `json:"peername,omitempty"`
	PeerDomain  string `json:"peerdomain,omitempty"`
	PeerPort    int    `json:"peerport,omitempty"`
	CouchdbPort int    `json:"couchdbport,omitempty"`
	NickName    string `json:"nickname,omitempty"`
	AccessKey   string `json:"accesskey,omitempty"`
	OrgName     string `json:"orgname,omitempty"`
	OrgDomain   string `json:"orgdomain,omitempty"`
	RunStatus   int    `json:"runstatus,omitempty"`
	PeerStatus  int    `json:"status,omitempty"`
	ChannelName string `json:"channelname,omitempty"`
}

type ChainCodeStruct struct {
	CCName         string   `json:"ccname,omitempty"`
	CCVersion      string   `json:"ccversion,omitempty"`
	EndorsementOrg []string `json:"endors,omitempty"`
	Channel        string   `json:"channel,omitempty"`
	Policy         string   `json:"policy,omitempty"`
	Desc           string   `json:"desc,omitempty"`
	IsInstall      int      `json:"isinstall,omitempty"`
	Status         int      `json:"status,omitempty"`
}

//StartDB 初始化数据库
func StartDB() error {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"

	userName = dconfig.GetStringByKey("mysqlUser")
	passWord = dconfig.GetStringByKey("mysqlPass")
	ip = dconfig.GetStringByKey("mysqlIP")
	port = dconfig.GetStringByKey("mysqlPort")
	dbName = dconfig.GetStringByKey("mysqlDBName")
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	db, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//GetAllInfoformIndent 获取所有端口
func GetAllInfoformIndent() (map[string]map[int]bool, error) {
	sqld := "select distinct peer_ip,peer_port,couchdb_port,cc_port,ca_port from ommp_indent"
	row, err := db.Query(sqld)
	defer row.Close()
	allIPMap, err := GetRowsValues(row)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query deploy result rows to map fail")
	}
	//读取
	one, err := AllInfoListToMap(allIPMap)
	if err != nil {
		return nil, errors.WithMessage(err, "deploy info map to mapfail")
	}
	return one, nil
}

//GetCurrtOrgPeerNum 获取当期通道组织下节点个数
func GetCurrtOrgPeerNum(orgName, channelName string) (int, error) {
	var count int
	sql := "select count(*) from ommp_indent where channelname = \"" + channelName + "\" and org_name = \"" + orgName + "\""
	err := db.QueryRow(sql).Scan(&count)
	if err != nil {
		return -1, errors.Errorf("mysql query channel=" + channelName + " orgName = " + orgName + " fail")
	}
	return count, nil
}

//GetAllPeerNum 获取网络中所有节点个数
func GetAllPeerNum() (int, error) {
	var count int
	sql := "select distinct peer_domain,peer_port from ommp_indent"
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		return -1, errors.Errorf("mysql get all peer num fail")
	}
	for row.Next() {
		count++
	}
	return count, nil
}

//GetOrgNumFromDB 获取去重组织个数
func GetOrgNumFromDB(channelName string) (int, error) {
	sqlq := "select distinct org_name from ommp_indent where channelname = \"" + channelName + "\" "
	var orgNum int
	row, err := db.Query(sqlq)
	defer row.Close()
	if err != nil {
		return -1, err
	}
	for row.Next() {
		orgNum++
	}

	return orgNum, nil
}

//GetStartTaskBeforIndent 获取任务之前的订单信息
func GetStartTaskBeforIndent(input *objectdefine.Indent) (*objectdefine.Indent, error) {
	channelName := input.ChannelName
	//首先查询订单表节点信息
	indentSQL := "select * from ommp_indent where channelname = \"" + channelName + "\""
	row, err := db.Query(indentSQL)
	defer row.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}

	//读取map信息转struct
	one, err := IndentMspToStruct(indentArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "peer info map to indent struct fail")
	}

	//查询orderer信息
	//orderer不要分通道
	//indentSQL = "select * from  orderer where channelname = \"" + channelName + "\""
	indentSQL = "select * from  ommp_orderer"
	rows, err := db.Query(indentSQL)
	defer rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_orderer fail")
	}
	ordererArrayMap, err := GetRowsValues(rows)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_orderer result rows to map fail")
	}

	two, err := IndentOrderToStruct(one, ordererArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "orderer info map to indent struct fail")
	}

	//查询chaincode信息
	indentSQL = "select * from ommp_chaincode where channelname = \"" + channelName + "\""
	rowC, err := db.Query(indentSQL)
	defer rowC.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_chaincode fail")
	}
	ccArrayMap, err := GetRowsValues(rowC)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_chaincode result rows to map fail")
	}
	indent, err := IndentCCToStruct(two, ccArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "chaincode info map to indent struct fail")
	}
	return indent, nil
}

//GetDeleteTaskBeforIndent 获取删除任务之前的订单信息
func GetDeleteTaskBeforIndent(input *objectdefine.Indent) (*objectdefine.Indent, error) {
	channelName := input.ChannelName
	//首先查询订单表节点信息
	indentSQL := "select * from ommp_indent where channelname = \"" + channelName + "\""
	row, err := db.Query(indentSQL)
	defer row.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}

	//读取map信息转struct
	one, err := IndentDeleteMspToStruct(indentArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "peer info map to indent struct fail")
	}

	//查询orderer信息
	//orderer不要分通道
	//indentSQL = "select * from  orderer where channelname = \"" + channelName + "\""
	indentSQL = "select * from  ommp_orderer"
	rows, err := db.Query(indentSQL)
	defer rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_orderer fail")
	}
	ordererArrayMap, err := GetRowsValues(rows)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_orderer result rows to map fail")
	}

	two, err := IndentOrderToStruct(one, ordererArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "orderer info map to indent struct fail")
	}

	//查询chaincode信息
	indentSQL = "select * from ommp_chaincode where channelname = \"" + channelName + "\""
	rowC, err := db.Query(indentSQL)
	defer rowC.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_chaincode fail")
	}
	ccArrayMap, err := GetRowsValues(rowC)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_chaincode result rows to map fail")
	}
	indent, err := IndentCCToStruct(two, ccArrayMap)
	if err != nil {
		return nil, errors.WithMessage(err, "chaincode info map to indent struct fail")
	}
	return indent, nil
}

//GetIPListFromIndent 获取ip列表和节点个数
func GetIPListFromIndent() (map[string]int, error) {

	indentSQL := "select distinct peer_ip from ommp_indent"
	row, err := db.Query(indentSQL)
	defer row.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	var ipList []string
	for _, ipL := range indentArrayMap {
		v := ipL["peer_ip"]
		ipList = append(ipList, v)
	}
	ipNumPair := make(map[string]int, len(ipList))
	for _, ip := range ipList {
		indentSQL := "select count(*) from ommp_indent where peer_ip = \"" + ip + "\""
		rowN, err := db.Query(indentSQL)
		defer rowN.Close()
		if err != nil {
			return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
		}
		peerNumMap, err := GetRowsValues(rowN)
		if err != nil {
			return nil, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
		}
		for _, numM := range peerNumMap {
			n := numM["count(*)"]
			num, _ := strconv.Atoi(n)
			ipNumPair[ip] = num
		}
	}
	return ipNumPair, nil
}

//GetIPPortListFromIndent 获取当期ip所有使用端口
func GetIPPortListFromIndent(ip string) ([]int, error) {

	sqld := "select peer_port,couchdb_port,cc_port,ca_port from ommp_indent where peer_ip = \"" + ip + "\""
	row, err := db.Query(sqld)
	defer row.Close()
	allPortMap, err := GetRowsValues(row)
	if err != nil {
		return nil, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	//读取
	one, err := AllPortListToMap(allPortMap)
	if err != nil {
		return nil, errors.WithMessage(err, "indent info map to mapfail")
	}
	return one, nil
}

//GetChannelListFromIndent 获取通道列表
func GetChannelListFromIndent(pageNum, pageSize int) ([]*channelStruct, int, error) {
	var indentSQL string
	if pageNum == 0 {
		indentSQL = "select distinct channelname,channel_status from ommp_indent"
	} else {
		count := pageNum * pageSize
		pageNs := count - pageSize
		indentSQL = fmt.Sprintf("select distinct channelname,channel_status from ommp_indent order by channelname desc limit %d,%d", pageNs, pageSize)
	}
	row, err := db.Query(indentSQL)
	defer row.Close()
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	//var channelList []string
	channelList := make(map[string]int)
	for _, ipL := range indentArrayMap {
		s := ipL["channelname"]
		v, _ := strconv.Atoi(ipL["channel_status"])
		//channelList = append(channelList, v)
		channelList[s] = v
	}

	//获取总通道个数
	total, err := GetChannelTotalNum()
	if err != nil || total == -1 {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent channel total num fail")
	}
	var pageTotal int
	if total <= pageSize {
		pageTotal = 1
		fmt.Println(pageTotal)
	} else {
		if total%pageSize == 0 {
			pageTotal = total / pageSize
		} else {
			pageTotal = (total / pageSize) + 1
		}
	}

	//获取组织个数和节点个数
	channelListInfo, err := GetChannelOrgAndPeerNum(channelList)
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent result fail")
	}
	return channelListInfo, pageTotal, nil

}

//GetChannelTotalNum
func GetChannelTotalNum() (int, error) {
	var channelNum int
	indentSQL := "select count(distinct channelname)  from ommp_indent"
	err := db.QueryRow(indentSQL).Scan(&channelNum)
	if err != nil {
		return -1, errors.WithMessage(err, "mysql query ommp_indent channel total num fail")
	}
	return channelNum, nil
}

//GetChannelOrgAndPeerNum 获取组织个数和节点个数 合约个数
func GetChannelOrgAndPeerNum(channelList map[string]int) ([]*channelStruct, error) {
	channelListInfo := make([]*channelStruct, 0)
	for chanenl, channelStatus := range channelList {
		indentSQL := "select distinct org_name from ommp_indent where channelname = \"" + chanenl + "\" and org_status !=2"
		row, err := db.Query(indentSQL)
		defer row.Close()
		if err != nil {
			return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
		}
		var orgNum int
		for row.Next() {
			orgNum++
		}
		indentSQL = "select distinct peer_domain from ommp_indent where channelname = \"" + chanenl + "\" and peer_status !=2"
		rowP, err := db.Query(indentSQL)
		defer rowP.Close()
		if err != nil {
			return nil, errors.WithMessage(err, "mysql query ommp_indent fail")
		}
		var peerNum int
		for rowP.Next() {
			peerNum++
		}
		var ccNum int
		indentSQL = "select count(*) from ommp_chaincode where channelname = \"" + chanenl + "\" and status !=2 "
		err = db.QueryRow(indentSQL).Scan(&ccNum)
		if err != nil {
			return nil, errors.WithMessage(err, "mysql query ommp_chaincode fail")
		}

		result := &channelStruct{}
		result.ChannelName = chanenl
		result.OrgNum = orgNum
		result.PeerNum = peerNum
		result.CCNum = ccNum
		result.ChannelStatus = channelStatus
		channelListInfo = append(channelListInfo, result)
	}
	return channelListInfo, nil
}

//GetChannelOrgInfoListFromIndent 获取通道下组织信息
func GetChannelOrgInfoListFromIndent(channel string, pageNum, pageSize int) ([]*OrgStruct, int, error) {
	var sqlSelect string
	if pageNum == 0 {
		if len(channel) == 0 {
			sqlSelect = "select distinct channelname,org_name,org_domain,org_status from ommp_indent"
		} else {
			sqlSelect = "select distinct channelname,org_name,org_domain,org_status from ommp_indent where channelname = \"" + channel + "\""
		}
	} else {
		count := pageNum * pageSize
		pageNs := count - pageSize
		if len(channel) == 0 {
			sqlS := "select distinct channelname,org_name,org_domain,org_status from ommp_indent"
			sqlSelect = fmt.Sprintf("%s order by org_name asc limit %d,%d", sqlS, pageNs, pageSize)
		} else {
			sqlS := "select distinct channelname,org_name,org_domain,org_status from ommp_indent where channelname = \"" + channel + "\""
			sqlSelect = fmt.Sprintf("%s order by org_name asc limit %d,%d", sqlS, pageNs, pageSize)
		}
	}

	row, err := db.Query(sqlSelect)
	defer row.Close()
	allOrgMap, err := GetRowsValues(row)
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	//orgInfo := make(map[string]string, 0)
	orgInfoArray := make([]*OrgStruct, 0)
	for _, orgM := range allOrgMap {
		orgName := orgM["org_name"]
		orgDomain := orgM["org_domain"]
		orgStatus, _ := strconv.Atoi(orgM["org_status"])

		//for orgN, orgD := range orgInfo {
		sqld := "select distinct peer_name from ommp_indent where org_name = \"" + orgName + "\" and peer_status !=2 and channelname = \"" + channel + "\""
		rowN, err := db.Query(sqld)
		defer rowN.Close()
		if err != nil {
			return nil, -1, errors.Errorf("mysql exec query peer num from ommp_indent fail %s", err)
		}
		var peerNum int64
		for rowN.Next() {
			peerNum++
		}
		result := &OrgStruct{}
		result.OrgName = orgName
		result.OrgDomain = orgDomain
		result.PeerNum = peerNum
		result.OrgStatus = orgStatus
		result.ChannelName = orgM["channelname"]
		orgInfoArray = append(orgInfoArray, result)
		//	}
	}
	//获取总组织个数
	total, err := GetOrgTotalNum(channel)
	if err != nil || total == -1 {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent org total num fail")
	}
	var pageTotal int
	if total <= pageSize {
		pageTotal = 1
		fmt.Println(pageTotal)
	} else {
		if total%pageSize == 0 {
			pageTotal = total / pageSize
		} else {
			pageTotal = (total / pageSize) + 1
		}
	}

	return orgInfoArray, pageTotal, nil
}

//GetOrgTotalNum 获取通道下组织总个数
func GetOrgTotalNum(channel string) (int, error) {
	var orgNum int
	indentSQL := "select count(distinct org_name)  from ommp_indent where channelname=\"" + channel + "\""
	err := db.QueryRow(indentSQL).Scan(&orgNum)
	if err != nil {
		return -1, errors.WithMessage(err, "mysql query ommp_indent org total num fail")
	}
	return orgNum, nil
}

//GetChannelOrgPeerInfoListFromIndent 获取组织下节点信息 或者所有节点
func GetChannelOrgPeerInfoListFromIndent(channel, orgName string, pageNum, pageSize int) ([]*PeerStruct, int, error) {
	var sqld string
	var sqlTotal string
	if pageNum == 0 {
		if len(orgName) == 0 {
			if len(channel) == 0 {
				sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent"
				sqlTotal = "select count(distinct channelname,peer_domain) from ommp_indent"
			} else {
				sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent where channelname = \"" + channel + "\""
				sqlTotal = "select count(distinct peer_domain) from ommp_indent where channelname = \"" + channel + "\""
			}
		} else {
			if len(channel) == 0 {
				return nil, -1, errors.New("channel name not exist")
			}
			sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent where channelname = \"" + channel + "\" and org_name = \"" + orgName + "\""
			sqlTotal = "select count(distinct peer_domain) from ommp_indent where channelname = \"" + channel + "\" and org_name = \"" + orgName + "\""
		}

	} else {
		count := pageNum * pageSize
		pageNs := count - pageSize
		//var sqlS string
		if len(orgName) == 0 {
			if len(channel) == 0 {
				sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent"
				sqlTotal = "select count(distinct channelname,peer_domain) from ommp_indent"
			} else {
				sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent where channelname = \"" + channel + "\""
				sqlTotal = "select count(distinct peer_domain) from ommp_indent where channelname = \"" + channel + "\""
			}
		} else {
			if len(channel) == 0 {
				return nil, -1, errors.New("channel name not exist")
			}
			sqld = "select channelname,peer_ip,peer_name,nick_name,accesskey,peer_domain,peer_port,couchdb_port,org_name,org_domain,peer_run_status,peer_status from ommp_indent where channelname = \"" + channel + "\" and org_name = \"" + orgName + "\""
			sqlTotal = "select count(distinct peer_domain) from ommp_indent where channelname = \"" + channel + "\" and org_name = \"" + orgName + "\""
		}
		sqld = fmt.Sprintf("%s order by org_name asc limit %d,%d", sqld, pageNs, pageSize)
	}

	row, err := db.Query(sqld)
	defer row.Close()
	allPeerMap, err := GetRowsValues(row)
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	//获取节点总个数Peer
	total, err := GetPeerTotalNum(sqlTotal)
	if err != nil || total == -1 {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent org total num fail")
	}
	var pageTotal int
	if total <= pageSize {
		pageTotal = 1
	} else {
		if total%pageSize == 0 {
			pageTotal = total / pageSize
		} else {
			pageTotal = (total / pageSize) + 1
		}
	}

	peerInfoArray := make([]*PeerStruct, 0)
	//peerInfo := make(map[string]string, 0)
	for _, orgM := range allPeerMap {
		peerIP := orgM["peer_ip"]
		peerName := orgM["peer_name"]
		peerNickName := orgM["nick_name"]
		peerAccessKey := orgM["accesskey"]
		peerDomain := orgM["peer_domain"]
		peerPort, _ := strconv.Atoi(orgM["peer_port"])
		couchdbPort, _ := strconv.Atoi(orgM["couchdb_port"])
		peerRunStatus, _ := strconv.Atoi(orgM["peer_run_status"])
		peerStatus, _ := strconv.Atoi(orgM["peer_status"])
		orgName := orgM["org_name"]
		orgDomain := orgM["org_domain"]
		channelName := orgM["channelname"]

		result := &PeerStruct{}
		result.IP = peerIP
		result.PeerName = peerName
		result.NickName = peerNickName
		result.AccessKey = peerAccessKey
		result.PeerDomain = peerDomain
		result.PeerPort = peerPort
		result.CouchdbPort = couchdbPort
		result.OrgName = orgName
		result.OrgDomain = orgDomain
		result.RunStatus = peerRunStatus
		result.PeerStatus = peerStatus
		result.ChannelName = channelName
		peerInfoArray = append(peerInfoArray, result)
	}

	return peerInfoArray, pageTotal, nil
}

//GetPeerTotalNum 获取节点总个数
func GetPeerTotalNum(sql string) (int, error) {
	var peerNum int
	//indentSQL := "select count(distinct org_name)  from ommp_indent where channelname=\"" + channel + "\""
	err := db.QueryRow(sql).Scan(&peerNum)
	if err != nil {
		return -1, errors.WithMessage(err, "mysql query ommp_indent peer total num fail")
	}
	return peerNum, nil
}

//GetChainCodeInfoListFromIndent 获取合约信息
func GetChainCodeInfoListFromIndent(channelName string, pageNum, pageSize int) ([]*ChainCodeStruct, int, error) {
	if len(channelName) == 0 {
		var sqld string
		var sqlTotal string
		if pageNum == 0 {
			sqld = "select * from ommp_chaincode"
		} else {
			count := pageNum * pageSize
			pageNs := count - pageSize
			sqld = fmt.Sprintf("select * from ommp_chaincode order by id asc limit %d,%d", pageNs, pageSize)
		}

		row, err := db.Query(sqld)
		defer row.Close()
		allCCMap, err := GetRowsValues(row)
		if err != nil {
			return nil, -1, errors.WithMessage(err, "mysql query ommp_chaincode result rows to map fail")
		}
		//获取合约总个数
		sqlTotal = "select count(*) from ommp_chaincode"
		total, err := GetChaincodeTotalNum(sqlTotal)
		if err != nil || total == -1 {
			return nil, -1, errors.WithMessage(err, "mysql query ommp_indent org total num fail")
		}
		var pageTotal int
		if total <= pageSize {
			pageTotal = 1
			fmt.Println(pageTotal)
		} else {
			if total%pageSize == 0 {
				pageTotal = total / pageSize
			} else {
				pageTotal = (total / pageSize) + 1
			}
		}
		ccInfoArray := make([]*ChainCodeStruct, 0)
		for _, ccM := range allCCMap {
			ccName := ccM["cc_name"]
			ccVersion := ccM["cc_version"]
			ccPolicy := ccM["cc_policy"]
			ccDesc := ccM["detail"]
			if len(ccDesc) == 0 {
				ccDesc = ""
			}
			var endorsOrg []string
			json.Unmarshal([]byte(ccM["cc_org"]), &endorsOrg)
			// ccEndors := endorsOrg
			ccChannel := ccM["channelname"]
			ccIsInstall, _ := strconv.Atoi(ccM["is_install"])
			ccStatus, _ := strconv.Atoi(ccM["status"])
			result := &ChainCodeStruct{}
			result.CCName = ccName
			result.CCVersion = ccVersion
			result.Policy = ccPolicy
			result.Desc = ccDesc
			result.EndorsementOrg = endorsOrg
			result.Channel = ccChannel
			result.IsInstall = ccIsInstall
			result.Status = ccStatus
			//result.TotalPage = pageTotal
			ccInfoArray = append(ccInfoArray, result)
		}
		return ccInfoArray, pageTotal, nil
	}
	var sqld string
	var sqlTotal string
	if pageNum == 0 {
		sqld = "select * from ommp_chaincode where channelname=\"" + channelName + "\""
	} else {
		count := pageNum * pageSize
		pageNs := count - pageSize
		sqls := "select * from ommp_chaincode where channelname=\"" + channelName + "\""
		sqld = fmt.Sprintf("%s order by id asc limit %d,%d", sqls, pageNs, pageSize)
	}
	row, err := db.Query(sqld)
	defer row.Close()
	allCCMap, err := GetRowsValues(row)
	if err != nil {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_chaincode result rows to map fail")
	}
	//获取合约总个数
	sqlTotal = "select count(*) from ommp_chaincode where channelname=\"" + channelName + "\""
	total, err := GetChaincodeTotalNum(sqlTotal)
	if err != nil || total == -1 {
		return nil, -1, errors.WithMessage(err, "mysql query ommp_indent org total num fail")
	}
	var pageTotal int
	if total <= pageSize {
		pageTotal = 1
		fmt.Println(pageTotal)
	} else {
		if total%pageSize == 0 {
			pageTotal = total / pageSize
		} else {
			pageTotal = (total / pageSize) + 1
		}
	}
	ccInfoArray := make([]*ChainCodeStruct, 0)
	for _, ccM := range allCCMap {
		ccName := ccM["cc_name"]
		ccVersion := ccM["cc_version"]
		ccPolicy := ccM["cc_policy"]
		ccDesc := ccM["detail"]
		if len(ccDesc) == 0 {
			ccDesc = ""
		}

		var endorsOrg []string
		json.Unmarshal([]byte(ccM["cc_org"]), &endorsOrg)
		// ccEndors := endorsOrg
		ccChannel := ccM["channelname"]
		ccIsInstall, _ := strconv.Atoi(ccM["is_install"])
		ccStatus, _ := strconv.Atoi(ccM["status"])
		result := &ChainCodeStruct{}
		result.CCName = ccName
		result.CCVersion = ccVersion
		result.Policy = ccPolicy
		result.Desc = ccDesc
		result.EndorsementOrg = endorsOrg
		result.Channel = ccChannel
		result.IsInstall = ccIsInstall
		result.Status = ccStatus
		//result.TotalPage = pageTotal
		ccInfoArray = append(ccInfoArray, result)
	}
	return ccInfoArray, pageTotal, nil

}

//GetChaincodeTotalNum 获取合约总个数
func GetChaincodeTotalNum(sql string) (int, error) {
	var ccNum int
	//indentSQL := "select count(distinct org_name)  from ommp_indent where channelname=\"" + channel + "\""
	err := db.QueryRow(sql).Scan(&ccNum)
	if err != nil {
		return -1, errors.WithMessage(err, "mysql query ommp_chaincode  total num fail")
	}
	return ccNum, nil
}

//GetRowsValues 拿取返回的rows内容
func GetRowsValues(row *sql.Rows) ([]map[string]string, error) {
	// Get column names
	columns, err := row.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	indentM := make([]map[string]string, 0)
	for row.Next() {
		err = row.Scan(scanArgs...)
		var value string
		indent := make(map[string]string, len(columns))
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
				indent[columns[i]] = value
			}
		}
		indentM = append(indentM, indent)
	}
	return indentM, nil
}

//IndentMspToStruct 读取map信息转struct
func IndentMspToStruct(indentMsp []map[string]string) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	orgMapType := make(map[string]objectdefine.OrgType, 0)

	fireWall := make(map[string]map[int]bool, 0)
	portCheck := make(map[int]bool, 0)
	for i, indent := range indentMsp {
		if i == 0 {
			ret.ID = indent["source_id"]
			ret.ChannelName = indent["channelname"]
			ret.SourceID = indent["source_id"]
			if len(ret.SourceID) == 0 {
				ret.SourceID = "sourceid"
			}
			ret.BaseOutput = filepath.ToSlash(dcache.GetOutputSubPath(ret.ID, ""))
			ret.SourceBaseOutput = filepath.ToSlash(dcache.GetOutputSubPath(ret.SourceID, ""))
			ret.Consensus = indent["consensus"]
			ret.Version = indent["version"]
		}
		if org, ok := orgMapType[indent["org_name"]]; ok {
			orgA := org
			peerA := org.Peer
			peerType := objectdefine.PeerType{}
			peerStatus, _ := strconv.Atoi(indent["peer_status"])
			if peerStatus == 1 {
				peerType.IP = indent["peer_ip"]
				peerType.Name = indent["peer_name"]
				peerType.User = indent["peer_user"]
				peerType.Port, _ = strconv.Atoi(indent["peer_port"])
				peerType.PeerID, _ = strconv.Atoi(indent["peer_id"])
				peerType.Status, _ = strconv.Atoi(indent["peer_status"])
				peerType.RunStatus, _ = strconv.Atoi(indent["peer_run_status"])
				portCheck[peerType.Port] = true
				fireWall[peerType.IP] = portCheck
				peerType.CouchdbPort, _ = strconv.Atoi(indent["couchdb_port"])
				portCheck[peerType.CouchdbPort] = true
				fireWall[peerType.IP] = portCheck
				peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
				portCheck[peerType.ChaincodePort] = true
				fireWall[peerType.IP] = portCheck
				//peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
				peerType.Domain = indent["peer_domain"]
				peerType.CliName = indent["cli_name"]
				peerType.NickName = indent["nick_name"]
				peerType.AccessKey = indent["accesskey"]
				peerA = append(peerA, peerType)
				orgA.Peer = peerA
				orgMapType[indent["org_name"]] = orgA
			}
		} else {
			//组织结构
			orgType := objectdefine.OrgType{}
			orgStatus, _ := strconv.Atoi(indent["org_status"])
			if orgStatus == 1 {
				orgName := indent["org_name"]
				orgType.Name = indent["org_name"]
				orgType.OrgDomain = indent["org_domain"]
				//节点结构
				peerArray := make([]objectdefine.PeerType, 0)
				peerType := objectdefine.PeerType{}
				peerStatus, _ := strconv.Atoi(indent["peer_status"])
				if peerStatus == 1 {
					peerType.IP = indent["peer_ip"]
					peerType.Name = indent["peer_name"]
					peerType.User = indent["peer_user"]
					peerType.Port, _ = strconv.Atoi(indent["peer_port"])
					peerType.PeerID, _ = strconv.Atoi(indent["peer_id"])
					peerType.Status, _ = strconv.Atoi(indent["peer_status"])
					peerType.RunStatus, _ = strconv.Atoi(indent["peer_run_status"])
					portCheck[peerType.Port] = true
					fireWall[peerType.IP] = portCheck
					peerType.CouchdbPort, _ = strconv.Atoi(indent["couchdb_port"])
					portCheck[peerType.CouchdbPort] = true
					fireWall[peerType.IP] = portCheck
					peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
					portCheck[peerType.ChaincodePort] = true
					fireWall[peerType.IP] = portCheck
					//peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
					peerType.Domain = indent["peer_domain"]
					peerType.CliName = indent["cli_name"]
					peerType.NickName = indent["nick_name"]
					peerType.AccessKey = indent["accesskey"]
					peerArray = append(peerArray, peerType)
					orgType.Peer = peerArray

					//ca
					caType := &objectdefine.CAType{}
					caType.IP = indent["ca_ip"]
					caType.Port, _ = strconv.Atoi(indent["ca_port"])
					portCheck[caType.Port] = true
					fireWall[caType.IP] = portCheck
					caType.Name = indent["ca_name"]
					orgType.CA = caType
					orgMapType[orgName] = orgType
				}
			}
		}
	}
	ret.Org = orgMapType
	ret.FireWall = fireWall
	return ret, nil

}

//IndentDeleteMspToStruct 读取map信息转struct
func IndentDeleteMspToStruct(indentMsp []map[string]string) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	orgMapType := make(map[string]objectdefine.OrgType, 0)

	fireWall := make(map[string]map[int]bool, 0)
	portCheck := make(map[int]bool, 0)
	for i, indent := range indentMsp {
		if i == 0 {
			ret.ID = indent["source_id"]
			ret.ChannelName = indent["channelname"]
			ret.SourceID = indent["source_id"]
			if len(ret.SourceID) == 0 {
				ret.SourceID = "sourceid"
			}
			ret.BaseOutput = filepath.ToSlash(dcache.GetOutputSubPath(ret.ID, ""))
			ret.SourceBaseOutput = filepath.ToSlash(dcache.GetOutputSubPath(ret.SourceID, ""))
			ret.Consensus = indent["consensus"]
			ret.Version = indent["version"]
		}
		if org, ok := orgMapType[indent["org_name"]]; ok {
			orgA := org
			peerA := org.Peer
			peerType := objectdefine.PeerType{}
			peerStatus, _ := strconv.Atoi(indent["peer_status"])
			if peerStatus == 1 || peerStatus == 2 || peerStatus == 3 {
				peerType.IP = indent["peer_ip"]
				peerType.Name = indent["peer_name"]
				peerType.User = indent["peer_user"]
				peerType.Port, _ = strconv.Atoi(indent["peer_port"])
				peerType.PeerID, _ = strconv.Atoi(indent["peer_id"])
				peerType.Status, _ = strconv.Atoi(indent["peer_status"])
				peerType.RunStatus, _ = strconv.Atoi(indent["peer_run_status"])
				portCheck[peerType.Port] = true
				fireWall[peerType.IP] = portCheck
				peerType.CouchdbPort, _ = strconv.Atoi(indent["couchdb_port"])
				portCheck[peerType.CouchdbPort] = true
				fireWall[peerType.IP] = portCheck
				peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
				portCheck[peerType.ChaincodePort] = true
				fireWall[peerType.IP] = portCheck
				//peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
				peerType.Domain = indent["peer_domain"]
				peerType.CliName = indent["cli_name"]
				peerType.NickName = indent["nick_name"]
				peerType.AccessKey = indent["accesskey"]
				peerA = append(peerA, peerType)
				orgA.Peer = peerA
				orgMapType[indent["org_name"]] = orgA
			}
		} else {
			//组织结构
			orgType := objectdefine.OrgType{}
			orgStatus, _ := strconv.Atoi(indent["org_status"])
			if orgStatus == 1 || orgStatus == 2 || orgStatus == 3 {
				orgName := indent["org_name"]
				orgType.Name = indent["org_name"]
				orgType.OrgDomain = indent["org_domain"]
				//节点结构
				peerArray := make([]objectdefine.PeerType, 0)
				peerType := objectdefine.PeerType{}
				peerStatus, _ := strconv.Atoi(indent["peer_status"])
				if peerStatus == 1 || peerStatus == 2 || peerStatus == 3 {
					peerType.IP = indent["peer_ip"]
					peerType.Name = indent["peer_name"]
					peerType.User = indent["peer_user"]
					peerType.Port, _ = strconv.Atoi(indent["peer_port"])
					peerType.PeerID, _ = strconv.Atoi(indent["peer_id"])
					peerType.Status, _ = strconv.Atoi(indent["peer_status"])
					peerType.RunStatus, _ = strconv.Atoi(indent["peer_run_status"])
					portCheck[peerType.Port] = true
					fireWall[peerType.IP] = portCheck
					peerType.CouchdbPort, _ = strconv.Atoi(indent["couchdb_port"])
					portCheck[peerType.CouchdbPort] = true
					fireWall[peerType.IP] = portCheck
					peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
					portCheck[peerType.ChaincodePort] = true
					fireWall[peerType.IP] = portCheck
					//peerType.ChaincodePort, _ = strconv.Atoi(indent["cc_port"])
					peerType.Domain = indent["peer_domain"]
					peerType.CliName = indent["cli_name"]
					peerType.NickName = indent["nick_name"]
					peerType.AccessKey = indent["accesskey"]
					peerArray = append(peerArray, peerType)
					orgType.Peer = peerArray

					//ca
					caType := &objectdefine.CAType{}
					caType.IP = indent["ca_ip"]
					caType.Port, _ = strconv.Atoi(indent["ca_port"])
					portCheck[caType.Port] = true
					fireWall[caType.IP] = portCheck
					caType.Name = indent["ca_name"]
					orgType.CA = caType
					orgMapType[orgName] = orgType
				}
			}
		}
	}
	ret.Org = orgMapType
	ret.FireWall = fireWall
	return ret, nil

}

//IndentOrderToStruct orderer转struct
func IndentOrderToStruct(indent *objectdefine.Indent, indentMsp []map[string]string) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	ret = indent
	ordererArray := make([]objectdefine.OrderType, 0)
	fireWall := ret.FireWall
	portCheck := make(map[int]bool, 0)
	for _, indent := range indentMsp {
		ordererType := objectdefine.OrderType{}
		ordererType.Name = indent["orderer_name"]
		ordererType.OrgDomain = indent["orderer_orgdomain"]
		ordererType.Domain = indent["orderer_domain"]
		ordererType.IP = indent["orderer_ip"]
		ordererType.Port, _ = strconv.Atoi(indent["orderer_port"])
		//portCheck = fireWall[ordererType.IP]
		portCheck[ordererType.Port] = true
		fireWall[ordererType.IP] = portCheck
		ordererArray = append(ordererArray, ordererType)
	}
	ret.Orderer = ordererArray
	return ret, nil
}

// IndentCCToStruct cc转struct
func IndentCCToStruct(indent *objectdefine.Indent, indentMsp []map[string]string) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	ret = indent
	ccStruct := make(map[string][]objectdefine.ChainCodeType, 0)
	//ccArray := make([]objectdefine.ChainCodeType, 0)

	for _, indent := range indentMsp {
		ccName := indent["cc_name"]
		if ccA, ok := ccStruct[ccName]; ok {
			ccAone := ccA
			ccType := objectdefine.ChainCodeType{}
			//ccStatus, _ := strconv.Atoi(indent["status"])
			//if ccStatus == 1 {
			ccType.Version = indent["cc_version"]
			ccType.Policy = indent["cc_policy"]
			if len(indent["detail"]) == 0 {
				ccType.Describe = ""
			} else {
				ccType.Describe = indent["detail"]
			}
			var endorseOrg []string
			json.Unmarshal([]byte(indent["cc_org"]), &endorseOrg)
			ccType.EndorsementOrg = endorseOrg
			ccType.Status, _ = strconv.Atoi(indent["status"])
			ccType.IsInstall, _ = strconv.Atoi(indent["is_install"])
			ccAone = append(ccAone, ccType)
			// ccArray = ccAone
			ccStruct[ccName] = ccAone
			//}
		} else {
			ccType := objectdefine.ChainCodeType{}
			//ccStatus, _ := strconv.Atoi(indent["status"])
			//if ccStatus == 1 {
			ccType.Version = indent["cc_version"]
			ccType.Policy = indent["cc_policy"]
			if len(indent["detail"]) == 0 {
				ccType.Describe = ""
			} else {
				ccType.Describe = indent["detail"]
			}
			var endorseOrg []string
			json.Unmarshal([]byte(indent["cc_org"]), &endorseOrg)
			ccType.EndorsementOrg = endorseOrg
			ccType.Status, _ = strconv.Atoi(indent["status"])
			ccType.IsInstall, _ = strconv.Atoi(indent["is_install"])
			ccArray := make([]objectdefine.ChainCodeType, 0)
			ccArray = append(ccArray, ccType)
			ccStruct[ccName] = ccArray
			//}
		}
	}
	ret.Chaincode = ccStruct
	return ret, nil
}

//AllInfoListToMap 获取所有机器启动的端口
func AllInfoListToMap(indentMsp []map[string]string) (map[string]map[int]bool, error) {
	ret := make(map[string]map[int]bool, len(indentMsp))
	//portCheck := make(map[int]bool, 0)
	for _, indent := range indentMsp {
		if port, ok := ret[indent["ip"]]; ok {
			peerPort, _ := strconv.Atoi(indent["peer_port"])
			if _, ok := port[peerPort]; !ok {
				portE := ret[indent["ip"]]
				portE[peerPort] = true
				ret[indent["ip"]] = portE
			}
			douchdbPort, _ := strconv.Atoi(indent["couchdb_port"])
			if _, ok := port[douchdbPort]; !ok {
				portE := ret[indent["ip"]]
				portE[douchdbPort] = true
				ret[indent["ip"]] = portE
			}
			ccPort, _ := strconv.Atoi(indent["cc_port"])
			if _, ok := port[ccPort]; !ok {
				portE := ret[indent["ip"]]
				portE[ccPort] = true
				ret[indent["ip"]] = portE
			}
			caPort, _ := strconv.Atoi(indent["ca_port"])
			if _, ok := port[caPort]; !ok {
				portE := ret[indent["ip"]]
				portE[caPort] = true
				ret[indent["ip"]] = portE
			}
		} else {
			peerPort, _ := strconv.Atoi(indent["peer_port"])
			portCheckP := make(map[int]bool, 0)
			portCheckP[peerPort] = true
			ret[indent["ip"]] = portCheckP
			douchdbPort, _ := strconv.Atoi(indent["couchdb_port"])
			portCheckC := make(map[int]bool, 0)
			portCheckC[douchdbPort] = true
			ret[indent["ip"]] = portCheckC
			ccPort, _ := strconv.Atoi(indent["cc_port"])
			portCheckCC := make(map[int]bool, 0)
			portCheckCC[ccPort] = true
			ret[indent["ip"]] = portCheckCC
			caPort, _ := strconv.Atoi(indent["ca_port"])
			portCheckCA := make(map[int]bool, 0)
			portCheckCA[caPort] = true
			ret[indent["ip"]] = portCheckCA
		}
	}

	return ret, nil

}

//AllPortListToMap 获取单台机器启动的端口
func AllPortListToMap(indentMsp []map[string]string) ([]int, error) {
	portList := make([]int, 0)
	for _, indent := range indentMsp {
		peerPort, _ := strconv.Atoi(indent["peer_port"])
		portList = append(portList, peerPort)
		douchdbPort, _ := strconv.Atoi(indent["couchdb_port"])
		portList = append(portList, douchdbPort)
		ccPort, _ := strconv.Atoi(indent["cc_port"])
		portList = append(portList, ccPort)
		caPort, _ := strconv.Atoi(indent["ca_port"])
		portList = append(portList, caPort)
	}

	return portList, nil
}

//GetOtherPeerList 查询当期组织其他节点  用来在peer配置文件的 Goosisp_bootstrap 环境变量
func GetOtherPeerList(channelName, orgName string) (string, error) {
	var otherPeer string
	sqld := "select peer_domain,peer_port from ommp_indent where channelname=\"" + channelName + "\" and org_name=\"" + orgName + "\""
	row, err := db.Query(sqld)
	defer row.Close()
	otherPeerMap, err := GetRowsValues(row)
	if err != nil {
		return "", errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}

	for _, opm := range otherPeerMap {
		peerDomain := opm["peer_domain"]
		peerPort, _ := strconv.Atoi(opm["peer_port"])
		str := fmt.Sprintf("%s:%d", peerDomain, peerPort)
		otherPeer += str + " "
	}
	return otherPeer, nil
}
