package db

import (
	"fmt"
	"os"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlHandler struct {
	conn        *gorm.DB
	initialized bool
}

func (handler *SqlHandler) GetConnection() *gorm.DB {
	if sqlHandler == nil || !sqlHandler.initialized {
		Init()
	}
	return sqlHandler.conn
}

var sqlHandler *SqlHandler

func Init() {
	if sqlHandler != nil && sqlHandler.initialized {
		return
	}

	// db.envに定義したDB関係の環境変数を取得
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	// tcp（）の中にdocker-composeで定義したDB用コンテナのサービス名を入れれば、
	// 自動的にホストとポートを読み取ってくれる
	dsn := fmt.Sprintf(
		"%s:%s@tcp(db)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser,
		dbPassword,
		dbName,
	)
	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database.")
	}

	// マイグレーションを行う
	dB.AutoMigrate(&models.Todo{})

	sqlHandler = &SqlHandler{
		conn:        dB,
		initialized: true,
	}
}

func GetSqlHandler() *SqlHandler {
	if sqlHandler == nil || !sqlHandler.initialized {
		Init()
	}
	return sqlHandler
}
