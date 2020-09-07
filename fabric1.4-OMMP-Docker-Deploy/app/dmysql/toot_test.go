package dmysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Test_qweqew(t *testing.T) {
	channelName := "1234"
	orgName := "qqq"
	sql := "select count(*) from  indent where channelname = \"" + channelName + "\" and org_name = \"" + orgName + "\""
	fmt.Println("sql", sql)
}

func Test_structInfo(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	channel := "score-channel"
	sqlq := "select * from  indent where channelname = \"" + channel + "\" "
	row, err := DB.Query(sqlq)
	if err != nil {
		fmt.Println("err", err)
	}
	// Get column names
	columns, err := row.Columns()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("columns", columns)
	values := make([]sql.RawBytes, len(columns))
	fmt.Println("values", values)
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
		fmt.Println("map", indentM)
	}
	fmt.Println("map", indentM)
}

func Test_distinct(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	channel := "score-channel"
	sqlq := "select * from  chaincode where channelname = \"" + channel + "\" "
	var count int
	row, err := DB.Query(sqlq)
	if err != nil {
		fmt.Println("Count", err)
	}

	for row.Next() {
		count += 1
	}
	fmt.Println("Count", count)
}

func Test_insert(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	mmm := make(map[string][]string, 10)
	aaa := []string{"123", "456"}
	bbb := []string{"123", "456"}
	mmm["mychannel"] = aaa
	mmm["testchannel"] = bbb

	mjson, _ := json.Marshal(mmm)
	mString := string(mjson)
	fmt.Printf("print mString:%s", mString)
	//channel := "score-channel"
	ip := "127.0.0.1"
	port := 7000
	sqlq := "insert into deploy (ip,port,channel_org) values (?,?,?)"
	fmt.Println("sqlq", sqlq)
	//var count int
	stem, err := DB.Exec(sqlq, ip, port, mString)
	if err != nil {
		fmt.Println("stem", err)
	}

	lastId, err := stem.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := stem.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID=%d, affected=%d\n", lastId, rowCnt)

}

func Test_Qinsert(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	mmm := make(map[string][]string, 10)
	aaa := []string{"123", "456"}
	bbb := []string{"123", "456"}
	mmm["mychannel"] = aaa
	mmm["testchannel"] = bbb

	mjson, _ := json.Marshal(mmm)
	mString := string(mjson)
	fmt.Printf("print mString:%s", mString)
	//channel := "score-channel"
	ip := "127.0.1.1"
	sqlq := "select * from  deploy where ip = \"" + ip + "\""
	fmt.Println("sqlq", sqlq)
	//var count int
	row, err := DB.Query(sqlq)
	if err != nil {
		fmt.Println("Count", err)
	}
	columns, err := row.Columns()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("columns", columns)
	values := make([]sql.RawBytes, len(columns))
	fmt.Println("values", values)
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
		fmt.Println("map", indentM)
	}
	kkk := make(map[string][]string, 0)
	for _, qqq := range indentM {
		mss := qqq["channel_org"]
		json.Unmarshal([]byte(mss), &kkk)
	}
	fmt.Println("kkk", kkk)

}

func Test_InitialsToUpper(t *testing.T) {
	s := "Org3"

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
		fmt.Println("zifu", string(by))
	}
	fmt.Println("zifu", s)
}

func Test_slice(t *testing.T) {
	s := "bysms192.168.8.13.tar"
	str := strings.Split(s, "192")
	fmt.Println("zifu", str)
}

func Test_updateM(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}

	ccPolicy := "OR('TonrenMSP.member')"
	endorseString := "[\"Tongren\"]"
	ccChannel := "score-channel"
	ccDesc := ""
	ccName := "test4"
	ccVersion := "1.0"
	sqlU := "update chaincode set cc_policy=?, cc_org=?,channelname=?,detail=?,is_install=? where cc_name=? and cc_version=?"
	result, err := DB.Exec(sqlU, ccPolicy, endorseString, ccChannel, ccDesc, 1, ccName, ccVersion)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := result.RowsAffected()
	if rows <= 0 {
		fmt.Println(err)
	}
	// stmt, err := DB.Prepare(`update chaincode set cc_policy=?, cc_org=?,channelname=?,detail=? where cc_name=? and cc_version=?`)
	// res, err := stmt.Exec(ccPolicy, endorseString, ccChannel, ccDesc, ccName, ccVersion)
	// rows, err := res.RowsAffected()
	// if rows <= 0 {
	// 	fmt.Println(err)
	// }
}

func Test_port(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	peeip := "192.168.8.13"
	sql := "select distinct peer_port,couchdb_port from  indent where peer_ip = \"" + peeip + "\""
	row, err := DB.Query(sql)
	defer row.Close()
	if err != nil {
		//return 0, errors.WithMessage(err, "mysql query indent fail")
	}
	portArrayMap, err := GetRowsValues(row)
	if err != nil {
		//return 0, errors.WithMessage(err, "mysql query indent result rows to map fail")
	}
	allPort := make(map[string]bool, 0)
	for _, port := range portArrayMap {
		peerPort := port["peer_port"]
		// couchdbPort := port["couchdb_port"]
		// caPort := port["ca_port"]
		allPort[peerPort] = true
	}
	peerPort := 7051
	couchdbPort := 5084

LOOP:
	for {
		peerPort += 100
		couchdbPort += 100
		pS := strconv.Itoa(peerPort)
		if _, ok := allPort[pS]; ok {
			goto LOOP
		}
		goto END
	}
END:
	fmt.Println(peerPort)

	caip := "192.168.8.13"
	sqlC := "select distinct ca_port from  indent where peer_ip = \"" + caip + "\""
	rowC, err := DB.Query(sqlC)
	defer row.Close()
	if err != nil {
		//return 0, errors.WithMessage(err, "mysql query indent fail")
	}
	caArrayMap, err := GetRowsValues(rowC)
	if err != nil {
		//return 0, errors.WithMessage(err, "mysql query indent result rows to map fail")
	}
	caPortMap := make(map[string]bool, 0)
	for _, ca := range caArrayMap {
		caPort := ca["ca_port"]
		// couchdbPort := port["couchdb_port"]
		// caPort := port["ca_port"]
		caPortMap[caPort] = true
	}
	fmt.Println(len(caPortMap))
	couchebPort := 7054
CLOOP:
	for {
		couchebPort += 1000
		cS := strconv.Itoa(couchebPort)
		if _, ok := caPortMap[cS]; ok {
			goto CLOOP
		}
		goto CEND
	}
CEND:
	fmt.Println(couchebPort)
}

func Test_stringEQ(t *testing.T) {
	ccVerison := "2.0"
	str := []string{"1.0", "1.1", "2.0", "3.0"}
	for _, ss := range str {
		if ss >= ccVerison {
			fmt.Printf("12345")
		}
	}
}

func Test_NumQE(t *testing.T) {
	var errBuf string
	errString := ""
	errInfoA := strings.Split(errString, "Error")
	if len(errInfoA) >= 2 {
		errBuf = errInfoA[1]
	} else {
		errBuf = errInfoA[0]
	}
	fmt.Println(errBuf)
}

func Test_quyu(t *testing.T) {
	total := 11
	pagesize := 2
	var pageTotal int
	if total <= pagesize {
		pageTotal = 1
		fmt.Println(pageTotal)
	}
	if total%pagesize == 0 {
		pageT := total / pagesize
		fmt.Println(pageT)
	} else {
		pageT := (total / pagesize) + 1
		fmt.Println(pageT)
	}
}

func Test_sqlRows(t *testing.T) {
	userName = "root"
	passWord = "123456"
	ip = "localhost"
	port = "3306"
	dbName = "fabric"
	path := strings.Join([]string{userName, ":", passWord, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("database path:", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Println(err)
	}
	ccName := "ee"
	ccVersion := "1.0"
	sqlS := "select count(*) from ommp_chaincode where cc_name=\"" + ccName + "\" and cc_version=\"" + ccVersion + "\" and is_install=0 and status=0"
	fmt.Println("sqls info ", ccName, ccVersion)
	var count int
	err := DB.QueryRow(sqlS).Scan(&count)
	if err != nil {

	}
	fmt.Println(count)
}
