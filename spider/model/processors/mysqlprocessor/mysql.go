package mysqlprocessor

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sync"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/l-dandelion/spider-go/lib/library/pool"
	"fmt"
)

var (
	DB *sql.DB
	tableMap map[string]bool
	tableMapLock sync.Mutex
	syncPool pool.Pool
)

func InitMysql(mysql string) {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysql)
	orm.RegisterModel(&C{})
	orm.RunSyncdb("default", false, true)
	syncPool.Init(20)
	go Start()
}

type C struct {
	Id int
}

func CreateTable(m *DBModel) error{
	tableName := m.Name
	sql := m.TableSql()
	tableMapLock.Lock()
	defer tableMapLock.Unlock()
	if tableMap == nil {
		tableMap = make(map[string]bool)
	}
	if _, ok := tableMap[tableName]; ok {
		return nil
	}
	o := orm.NewOrm()
	_,err := o.Raw(sql).Exec()
	if err != nil{
		return err
	}
	tableMap[tableName] = true
	beego.Info("创建表 ",m.Name," 成功 【完成】")
	return nil
}

func Add(m *DBModel) error{
	o := orm.NewOrm()
	_,err := o.Raw(m.InsertSql(),m.InsertArgs()...).Exec()
	if err != nil{
		return err
	}
	beego.Info("插入数据成功")
	return nil
}

var (
	sqlMap map[string][]*DBModel
	sqlMapLock sync.Mutex
	CanPrepare bool = true
)

func AddPrepare(m *DBModel) {
	for ;!CanPrepare; {
		time.Sleep(100*time.Millisecond)
	}
	sqlMapLock.Lock()
	defer sqlMapLock.Unlock()
	if sqlMap == nil {
		sqlMap = map[string][]*DBModel{}
	}
	tmp := sqlMap[m.InsertSql()]
	if tmp == nil {
		tmp = []*DBModel{}
	}
	tmp = append(tmp, m)
	sqlMap[m.InsertSql()] = tmp
}

func Exec() {
	CanPrepare = false
	defer func(){
		CanPrepare = true
	}()
	sqlMapLock.Lock()
	defer sqlMapLock.Unlock()
	if sqlMap == nil {
		return
	}
	for key, models := range sqlMap {
		if len(models) != 0 {
			syncPool.Add()
			go func(models []*DBModel) {
				defer syncPool.Done()
				o := orm.NewOrm()
				err := CreateTable(models[0])
				if err != nil {
					log.Info("Mysql Create Table ERROR:", err)
				}
				sql := GenInsertModelsSql(models)
				args := GenInsertModelsArgs(models)
				_, err = o.Raw(sql, args...).Exec()
				if err != nil {
					log.Info("Mysql Exec ERROR:", err)
				}
			}(models)
			sqlMap[key] = []*DBModel{}
		}
	}
}

func Start() {
	for {
		Exec()
		time.Sleep(5*time.Second)
	}
}

func GenInsertModelsSql(models []*DBModel) string {
	model := models[0]
	sql := fmt.Sprintf("INSERT INTO `%s`", model.Name)
	modelsLen := len(models)
	filedsLen := len(model.Fields)
	filedNamesStr := "("
	tmp := "("
	for j := 1; j < filedsLen; j ++ {
		if j !=1 {
			tmp = tmp + ",?"
			filedNamesStr = filedNamesStr + fmt.Sprintf(",`%s`", model.Fields[j].Name)
		} else {
			tmp = tmp + "?"
			filedNamesStr = filedNamesStr + fmt.Sprintf("`%s`", model.Fields[j].Name)
		}
	}
	tmp += ")"
	filedNamesStr += ")"
	sql += filedNamesStr
	for i := 0; i < modelsLen; i ++ {
		if i == 0 {
			sql += "VALUES" + tmp
		} else {
			sql += "," + tmp
		}
	}
	return sql
}

func GenInsertModelsArgs(models []*DBModel) []interface{} {
	args := []interface{}{}
	for _, model := range models {
		args = append(args, model.InsertArgs()...)
	}
	return args
}