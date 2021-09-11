package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var (
	MYSQLDB  = "root:mysql@tcp(127.0.0.1:3306)/fund"
	SQLITEDB = "fund.db"
)

var (
	DBconn    *gorm.DB
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 10, // 慢 SQL 阈值
			LogLevel:      logger.Error,     // Log level
			Colorful:      false,            // 禁用彩色打印
		},
	)
	config = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	}
)

func MysqlInit() {
	//连接数据库
	db, err := gorm.Open(mysql.Open(MYSQLDB), config)
	if err != nil {
		panic(err)
	}
	mySQL, err := db.DB()
	if err != nil {
		panic(err)
	}
	//设置最大空闲连接
	mySQL.SetMaxIdleConns(10)
	//设置最大连接数
	mySQL.SetMaxOpenConns(100)
	//设置连接超时时间:1分钟
	mySQL.SetConnMaxLifetime(time.Minute)
	// 设置了连接可复用的最大时间
	mySQL.SetConnMaxLifetime(time.Hour)
	DBconn = db
}

func SqliteInit() {
	db, err := gorm.Open(sqlite.Open(SQLITEDB), config)
	if err != nil {
		panic("failed to connect database")
	}
	sqliteDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	//设置最大空闲连接
	sqliteDB.SetMaxIdleConns(10)
	//设置最大连接数
	sqliteDB.SetMaxOpenConns(100)
	//设置连接超时时间:1分钟
	sqliteDB.SetConnMaxLifetime(time.Minute)
	// 设置了连接可复用的最大时间
	sqliteDB.SetConnMaxLifetime(time.Hour)
	DBconn = db
}
