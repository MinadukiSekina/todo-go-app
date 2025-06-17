package db

import (
	"fmt"
	"os"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// データベース接続を管理するインターフェース
type SqlHandler interface {
	GetConnection() *gorm.DB
}

// データベース接続の実装を保持する構造体
type sqlHandler struct {
	conn        *gorm.DB
	initialized bool
}

// データベース接続を取得する
// 接続が初期化されていない場合は初期化を行う
func (handler *sqlHandler) GetConnection() *gorm.DB {
	if handler == nil || !handler.initialized {
		Init()
	}
	return handler.conn
}

// シングルトンインスタンスを保持するグローバル変数
var handler *sqlHandler

// データベース接続を初期化する
// 環境変数から接続情報を取得し、GORMを使用して接続を確立する
// 接続後、Todoモデルのマイグレーションを実行する
func Init() {
	if handler != nil && handler.initialized {
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

	handler = &sqlHandler{
		conn:        dB,
		initialized: true,
	}
}

// SqlHandlerのインスタンスを取得する
// インスタンスが存在しない場合は初期化を行う
func GetSqlHandler() *sqlHandler {
	if handler == nil || !handler.initialized {
		Init()
	}
	return handler
}
