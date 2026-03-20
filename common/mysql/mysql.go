package mysql

import (
	"GoNexus/config"
	"GoNexus/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMysql 初始化mysql
func InitMysql() error {
	port := config.GetConfig().MysqlPort
	host := config.GetConfig().MysqlHost
	user := config.GetConfig().MysqlUser
	password := config.GetConfig().MysqlPassword
	databaseName := config.GetConfig().MysqlDatabaseName
	charset := config.GetConfig().MysqlCharset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		user, password, host, port, databaseName, charset)

	var log logger.Interface
	if gin.Mode() == "debug" {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default
	}

	// 初始化并建立一个到Mysql数据库的GORM连接实例
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{Logger: log})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	return migration()
}

// migration 自动根据model创建对应的mysql数据库表
func migration() error {
	return DB.AutoMigrate(
		&model.User{},
	)
}

// GetUserByUsername 根据用户名查询用户信息
func GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := DB.Where("username = ?", username).First(user).Error
	return user, err
}

// GetUserByEmail 根据邮箱查询用户信息
func GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := DB.Where("email = ?", email).First(user).Error
	return user, err
}

// InsertUser 写入用户信息
func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	return user, err
}
