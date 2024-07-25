package db

import (
	"XArr-Rss/app/model/dbmodel"
	"github.com/glebarez/sqlite"

	//"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	//"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

var MainDb *gorm.DB

func InitDb(confDir string) *gorm.DB {
	// 连接sqlite
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)
	//_pragma=journal_mode(WAL)
	drive := sqlite.Open(confDir + "/xarr-rss.db?cache=shared&mode=rwc&_journal_mode=WAL&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&loc=local")

	MainDb, err = gorm.Open(drive, &gorm.Config{
		//SkipDefaultTransaction:                   false,
		//NamingStrategy:                           nil,
		//FullSaveAssociations:                     false,
		Logger: newLogger,
		//NowFunc:                                  nil,
		//DryRun:                                   false,
		//PrepareStmt:                              false,
		//DisableAutomaticPing:                     false,
		//DisableForeignKeyConstraintWhenMigrating: false,
		//DisableNestedTransaction:                 false,
		//AllowGlobalUpdate:                        false,
		//QueryFields:                              false,
		//CreateBatchSize:                          0,
		//ClauseBuilders:                           nil,
		//ConnPool:                                 nil,
		//Dialector:                                nil,
		//Plugins:                                  nil,
	})

	if err != nil {
		panic("数据库连接失败,请检查 conf/xarr-rss.db 权限")
	}

	// 去掉复数

	// 迁移 schema
	//MainDb.SingularTable
	MainDb.AutoMigrate(&dbmodel.Group{})
	MainDb.AutoMigrate(&dbmodel.GroupMedia{})
	MainDb.AutoMigrate(&dbmodel.Sources{})
	MainDb.AutoMigrate(&dbmodel.SourceItem{})
	MainDb.AutoMigrate(&dbmodel.Media{})
	MainDb.AutoMigrate(&dbmodel.GroupTemplate{})
	if !MainDb.Migrator().HasColumn(&dbmodel.Media{}, "tags") {
		MainDb.Migrator().AddColumn(&dbmodel.Media{}, "tags")
	}

	MainDb.Exec("DROP TABLE IF EXISTS media_episodes")

	MainDb.AutoMigrate(&dbmodel.MediaEpisodeList{})
	MainDb.AutoMigrate(&dbmodel.Options{})
	MainDb.AutoMigrate(&dbmodel.DownloadList{})

	db, err := MainDb.DB()
	if err != nil {
		panic("数据库获取失败:" + err.Error())
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	return MainDb
}
