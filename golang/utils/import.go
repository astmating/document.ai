package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type DBTool struct {
	SqlPath                                    string
	Username, Password, Server, Port, Database string
}

func (this *DBTool) connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", this.Username, this.Password, this.Server, this.Port, this.Database)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("数据库连接失败:", err)
		//panic("数据库连接失败!")
		return nil, err
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(59 * time.Second)
	return db, nil
}

//执行查询语句
func (this *DBTool) QuerySql(sqlStr string) ([]map[string]interface{}, error) {
	db, err := this.connect()
	var result []map[string]interface{}
	if err != nil {
		return result, err
	}
	sqlStr = strings.TrimSpace(sqlStr)
	if sqlStr == "" {
		return result, errors.New("sql语句为空")
	}
	rows, err := db.DB().Query(sqlStr)
	if err != nil {
		return result, err
	}

	//获取列名
	columns, _ := rows.Columns()

	//定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
	values := make([]sql.RawBytes, len(columns))
	//定义一个切片,元素类型是interface{} 接口
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		//把sql.RawBytes类型的地址存进去了
		scanArgs[i] = &values[i]
	}
	//获取字段值
	for rows.Next() {
		res := make(map[string]interface{})
		rows.Scan(scanArgs...)
		for i, col := range values {
			res[columns[i]] = string(col)
		}
		result = append(result, res)
	}
	return result, nil
}

//执行增删改
func (this *DBTool) ExecuteSql(sql string) error {
	db, err := this.connect()
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return errors.New("sql语句为空")
	}
	err = db.Exec(sql).Error
	if err != nil {
		log.Println("执行失败:" + err.Error())
		return err
	} else {
		log.Println(sql, "\t success!")
	}
	return nil
}

//执行sql文件
func (this *DBTool) ImportSql() error {
	_, err := os.Stat(this.SqlPath)
	if os.IsNotExist(err) {
		log.Println("数据库SQL文件不存在:", err)
		return err
	}

	db, err := this.connect()

	sqls, _ := ioutil.ReadFile(this.SqlPath)
	sqlArr := strings.Split(string(sqls), ";")
	for _, sql := range sqlArr {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		err := db.Exec(sql).Error
		if err != nil {
			log.Println("数据库导入失败:" + err.Error())
			return err
		} else {
			log.Println(sql, "\t success!")
		}
	}
	return nil
}
