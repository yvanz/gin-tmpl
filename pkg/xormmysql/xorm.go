/*
@Date: 2021/1/15 下午6:59
@Author: yvan.zhang
@File : xorm
@Desc:
*/

package xormmysql

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"xorm.io/xorm"
	xLog "xorm.io/xorm/log"
)

type MyDB interface {
	Master() *xorm.Engine
	Slave() *xorm.Engine
	Close() error
}

var (
	// 仅用于单例模式下
	_default *DB
)

func My() MyDB {
	return _default
}

type DB struct {
	read    []*xorm.Engine
	write   *xorm.Engine
	isGroup bool
	*xorm.EngineGroup
}

func (c DBConfig) BuildMySQLClient() (*DB, error) {
	var err error
	var master *xorm.Engine
	_default = new(DB)
	if c.WriteDB.Host == "" {
		_default.write = nil
	} else {
		master, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&timeout=5s",
			c.WriteDB.User,
			c.WriteDB.Password,
			c.WriteDB.Host,
			c.WriteDB.Port,
			c.WriteDB.Database),
		)
		if err != nil {
			return nil, err
		}
		_default.write = master
	}

	slaves := make([]*xorm.Engine, 0)
	for _, cf := range c.ReadDB {
		if cf.Host == "" {
			continue
		}
		slave, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&timeout=5s",
			cf.User,
			cf.Password,
			cf.Host,
			cf.Port,
			cf.Database),
		)
		if err != nil {
			return nil, err
		}
		slaves = append(slaves, slave)
	}

	if len(slaves) > 0 {
		if master != nil {
			engines, err := xorm.NewEngineGroup(master, slaves, xorm.LeastConnPolicy())
			if err != nil {
				return nil, err
			}
			_default.isGroup = true
			_default.EngineGroup = engines
			_default.write = engines.Master()
			_default.read = engines.Slaves()
		} else {
			_default.read = slaves
		}
	} else {
		_default.read = []*xorm.Engine{_default.write}
	}

	if _default.write != nil {
		_default.write.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	}

	c.setTablePrefixAndConn()

	logFile, err := c.mysqlLogFile()
	if err != nil {
		return nil, err
	}

	if err := c.clientOtherSettings(logFile); err != nil {
		return nil, err
	}

	return _default, err
}

func (c DBConfig) setTablePrefixAndConn() {
	if len(c.Prefix) > 0 {
		tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, c.Prefix)
		if _default.write != nil {
			_default.write.SetTableMapper(tbMapper)
		}
		for i := 0; i < len(_default.read); i++ {
			_default.read[i].SetTableMapper(tbMapper)
		}
	}

	if c.MaxIdleConns > 0 {
		if _default.write != nil {
			_default.write.SetMaxIdleConns(c.MaxIdleConns)
		}
		for i := 0; i < len(_default.read); i++ {
			_default.read[i].SetMaxIdleConns(c.MaxIdleConns)
		}
	}

	if c.MaxOpenConns > 0 {
		if _default.write != nil {
			_default.write.SetMaxOpenConns(c.MaxOpenConns)
		}
		for i := 0; i < len(_default.read); i++ {
			_default.read[i].SetMaxOpenConns(c.MaxOpenConns)
		}
	}
}

func (c DBConfig) mysqlLogFile() (*os.File, error) {
	var (
		dir  string
		file string
		err  error
	)

	if c.LogDir == "" {
		return os.Stdout, nil
	}

	dir, err = filepath.Abs(c.LogDir)
	if err != nil {
		return nil, err
	}
	f, err := os.Stat(dir)
	if !os.IsNotExist(err) && !f.IsDir() {
		return nil, errors.New("mysql log path is not a directory")
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	file = path.Join(dir, "mysql.log")

	myLog, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return myLog, nil
}

func (c DBConfig) clientOtherSettings(logFile *os.File) (err error) {
	logger := xLog.NewSimpleLogger(logFile)

	if c.Logging {
		logger.ShowSQL(true)
		switch c.LogLevel {
		case `debug`:
			logger.SetLevel(xLog.LOG_DEBUG)
		case `info`:
			logger.SetLevel(xLog.LOG_INFO)
		case `warn`:
			logger.SetLevel(xLog.LOG_WARNING)
		case `error`:
			logger.SetLevel(xLog.LOG_ERR)
		default:
			logger.SetLevel(xLog.LOG_INFO)
		}
	}

	if _default.write != nil {
		_default.write.SetLogger(logger)
	}
	for i := 0; i < len(_default.read); i++ {
		_default.read[i].SetLogger(logger)
	}

	if err := _default.Ping(); err != nil {
		return err
	}

	go func() {
		timer := time.NewTicker(30 * time.Second)
		for {
			<-timer.C
			_ = _default.Ping()
		}
	}()

	return err
}

func (d *DB) Ping() error {
	if d.write != nil {
		if err := d.write.Ping(); err != nil {
			return err
		}
	}
	for i := 0; i < len(d.read); i++ {
		if err := d.read[i].Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) Master() *xorm.Engine {
	defer func() {
		if r := recover(); r != nil {
			d.Logger().Errorf("%v", r)
		}
	}()

	if d.write == nil {
		panic("xorm mysql master not configure")
	}

	return d.write
}

func (d *DB) Slave() *xorm.Engine {
	if d.isGroup {
		return d.EngineGroup.Slave()
	}

	if len(d.Slaves()) == 1 {
		return d.Slaves()[0]
	}

	rand.Seed(time.Now().UnixNano())
	return d.Slaves()[rand.Intn(len(d.Slaves()))]
}

func (d *DB) Slaves() []*xorm.Engine {
	return d.read
}

func (d *DB) Close() error {
	var err error
	if d.isGroup {
		return d.EngineGroup.Close()
	}

	if d.write != nil {
		err = d.write.Close()
		slaves := d.Slaves()
		for i := 0; i < len(slaves); i++ {
			err = slaves[i].Close()
		}
		return err
	}

	return err
}
