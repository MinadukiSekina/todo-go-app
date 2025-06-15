package db

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// テスト用のDBハンドラー
type testHandler struct {
	conn *gorm.DB
}

// SqlHandlerインターフェイスの実装
func (th *testHandler) GetConnection() *gorm.DB {
	return th.conn
}

// テストスイートの構造体
type todoRepositoryTestSuite struct {
	suite.Suite
	dbConn *gorm.DB
}

// テストスイートを実行する
func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(todoRepositoryTestSuite))
}

// テストスイートの実行前に処理される
func (s *todoRepositoryTestSuite) SetupSuite() {
	// db.envに定義したDB関係の環境変数を取得
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("TEST_DATABASE")
	// tcp（）の中にdocker-composeで定義したDB用コンテナのサービス名を入れれば、
	// 自動的にホストとポートを読み取ってくれる
	dsn := fmt.Sprintf(
		"%s:%s@tcp(db)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser,
		dbPassword,
		dbName,
	)

	// txdbに登録する
	txdb.Register("txdb", "mysql", dsn)

	// マイグレーションは一度だけ実行
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		s.Failf("failed to connect to database", "%v", err)
	}
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		s.Failf("failed to migrate database", "%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		s.Failf("failed to get underlying sql.DB", "%v", err)
	}

	// マイグレーション用の接続を閉じる
	defer sqlDB.Close()
}

// 各テスト終了時にDB接続を閉じる
func (s *todoRepositoryTestSuite) Close() {
	// テスト用に開いた接続を閉じる
	sqlDB, err := s.dbConn.DB()
	if err != nil {
		s.Failf("failed to get underlying sql.DB", "%v", err)
	}
	defer sqlDB.Close()
}

func (s *todoRepositoryTestSuite) TestFindAll() {

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}

	todo2 := models.Todo{Title: "test2", Status: models.NotStarted}

	nothingTodos := []models.Todo{}
	onlyOneTodos := []models.Todo{todo1}
	manyHasTodos := []models.Todo{todo1, todo2}

	cases := map[string]struct {
		want      *[]models.Todo
		expectErr bool
		err       error
	}{
		"正常ケース:データなし": {
			want:      &nothingTodos,
			expectErr: false,
			err:       nil,
		},
		"正常ケース:1件データあり": {
			want:      &onlyOneTodos,
			expectErr: false,
			err:       nil,
		},
		"正常ケース:2件データあり": {
			want:      &manyHasTodos,
			expectErr: false,
			err:       nil,
		},
	}
	for name, tt := range cases {
		s.T().Run(name, func(t *testing.T) {
			// テスト用DBに接続する
			db, err := gorm.Open(mysql.New(mysql.Config{DSN: uuid.NewString(), DriverName: "txdb"}))
			if err != nil {
				s.Failf("database connection is not established", "%v", err)
			}
			// コネクションを格納する
			s.dbConn = db

			defer s.Close()

			// 初期処理
			sqlHandler := testHandler{conn: s.dbConn}
			todoRepository := NewTodoRepository(&sqlHandler)

			// 期待する結果が空以外の場合
			if len(*tt.want) > 0 {
				// データが無いはずなので登録する
				for i := range *tt.want {
					result := s.dbConn.Create(&(*tt.want)[i])
					if result.Error != nil {
						s.T().Errorf("Creation is failed. error: %v", result.Error)
					}
				}
			}

			todos, err := todoRepository.FindAll()

			// 結果を確認
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
				assert.Nil(t, todos)
			} else {
				if assert.NoError(t, err) {
					assert.ElementsMatch(t, *tt.want, *todos)
				}
			}
		})
	}
}

func (s *todoRepositoryTestSuite) TestFindById() {

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo2 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo2.ID = 1

	cases := map[string]struct {
		want      *models.Todo
		expectErr bool
		err       error
		setup     func(*gorm.DB)
	}{
		"正常ケース:指定したIDのデータあり": {
			want:      &todo1,
			expectErr: false,
			err:       nil,
			setup:     func(d *gorm.DB) { _ = d.Create(&todo1) },
		},
		"異常ケース:指定したIDのデータが無い": {
			want:      &todo2,
			expectErr: true,
			err:       errors.New("record not found"),
			setup:     func(d *gorm.DB) {},
		},
	}
	for name, tt := range cases {
		s.T().Run(name, func(t *testing.T) {
			// テスト用DBに接続する
			db, err := gorm.Open(mysql.New(mysql.Config{DSN: uuid.NewString(), DriverName: "txdb"}))
			if err != nil {
				s.Failf("database connection is not established", "%v", err)
			}
			// コネクションを格納する
			s.dbConn = db

			defer s.Close()

			// セットアップ関数を実行する
			tt.setup(s.dbConn)

			// 初期処理
			sqlHandler := testHandler{conn: s.dbConn}
			todoRepository := NewTodoRepository(&sqlHandler)

			todo, err := todoRepository.FindById(tt.want.ID)

			// 結果を確認
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
				assert.Nil(t, todo)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, *tt.want, *todo)
				}
			}
		})
	}
}

func (s *todoRepositoryTestSuite) TestCreate() {

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}

	cases := map[string]struct {
		want      *models.Todo
		expectErr bool
		err       error
	}{
		"正常ケース:作成に成功": {
			want:      &todo1,
			expectErr: false,
			err:       nil,
		},
	}
	for name, tt := range cases {
		s.T().Run(name, func(t *testing.T) {
			// テスト用DBに接続する
			db, err := gorm.Open(mysql.New(mysql.Config{DSN: uuid.NewString(), DriverName: "txdb"}))
			if err != nil {
				s.Failf("database connection is not established", "%v", err)
			}
			// コネクションを格納する
			s.dbConn = db

			defer s.Close()

			// 初期処理
			sqlHandler := testHandler{conn: s.dbConn}
			todoRepository := NewTodoRepository(&sqlHandler)

			err = todoRepository.Create(tt.want)

			// 結果を確認
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				if assert.NoError(t, err) {
					todo, err := todoRepository.FindById(tt.want.ID)
					if err != nil {
						s.Failf("can't get todo. error: %v", err.Error())
					}
					assert.Equal(t, *tt.want, *todo)
				}
			}
		})
	}
}

func (s *todoRepositoryTestSuite) TestUpdate() {
	// 更新前のデータ
	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	// 存在しないデータを更新するケース用
	notExistTodo := models.Todo{Title: "not_exist_record", Status: models.NotStarted}

	cases := map[string]struct {
		before    *models.Todo
		update    func(*models.Todo) // 更新するフィールドを指定する関数
		expectErr bool
		err       error
		setup     func(*gorm.DB)
	}{
		"正常ケース:更新成功": {
			before: &todo1,
			update: func(t *models.Todo) {
				t.Title = "test1_updated"
				t.Status = models.Done
			},
			expectErr: false,
			err:       nil,
			setup:     func(d *gorm.DB) { _ = d.Create(&todo1) },
		},
		"異常ケース:存在しないIDの更新": {
			before: &notExistTodo,
			update: func(t *models.Todo) {
				t.Title = "not_exist_updated"
				t.Status = models.Done
			},
			expectErr: true,
			err:       errors.New("record not found"),
			setup:     func(d *gorm.DB) {}, // 何もしない
		},
	}

	for name, tt := range cases {
		s.T().Run(name, func(t *testing.T) {
			// テスト用DBに接続する
			db, err := gorm.Open(mysql.New(mysql.Config{DSN: uuid.NewString(), DriverName: "txdb"}))
			if err != nil {
				s.Failf("database connection is not established", "%v", err)
			}
			// コネクションを格納する
			s.dbConn = db

			defer s.Close()

			// セットアップ関数を実行する
			tt.setup(s.dbConn)

			// 初期処理
			sqlHandler := testHandler{conn: s.dbConn}
			todoRepository := NewTodoRepository(&sqlHandler)

			// 更新対象のレコードを取得
			todo, err := todoRepository.FindById(tt.before.ID)

			// 存在しないIDの場合は、直接Updateを実行する
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
				// 存在しないIDのレコードを直接更新
				tt.update(tt.before)
				err = todoRepository.Update(tt.before)
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
				return
			}

			if assert.NoError(t, err) {
				// 更新するフィールドを設定
				tt.update(todo)

				// 更新処理を実行
				err = todoRepository.Update(todo)

				// 結果を確認
				if assert.NoError(t, err) {
					// 更新後のデータを取得して確認
					updated, err := todoRepository.FindById(todo.ID)
					if assert.NoError(t, err) {
						assert.Equal(t, todo.Title, updated.Title)
						assert.Equal(t, todo.Status, updated.Status)
					}
				}
			}
		})
	}
}

func (s *todoRepositoryTestSuite) TestDelete() {
	// 更新前のデータ
	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	// 存在しないデータを更新するケース用
	notExistTodo := models.Todo{Title: "not_exist_record", Status: models.NotStarted}

	cases := map[string]struct {
		todo      *models.Todo
		expectErr bool
		err       error
		setup     func(*gorm.DB)
	}{
		"正常ケース:削除成功": {
			todo:      &todo1,
			expectErr: false,
			err:       errors.New("record not found"),
			setup:     func(d *gorm.DB) { _ = d.Create(&todo1) },
		},
		"異常ケース:存在しないIDの削除": {
			todo:      &notExistTodo,
			expectErr: true,
			err:       errors.New("record not found"),
			setup:     func(d *gorm.DB) {}, // 何もしない
		},
	}

	for name, tt := range cases {
		s.T().Run(name, func(t *testing.T) {
			// テスト用DBに接続する
			db, err := gorm.Open(mysql.New(mysql.Config{DSN: uuid.NewString(), DriverName: "txdb"}))
			if err != nil {
				s.Failf("database connection is not established", "%v", err)
			}
			// コネクションを格納する
			s.dbConn = db

			defer s.Close()

			// セットアップ関数を実行する
			tt.setup(s.dbConn)

			// 初期処理
			sqlHandler := testHandler{conn: s.dbConn}
			todoRepository := NewTodoRepository(&sqlHandler)

			// 削除処理を実行する
			err = todoRepository.Delete(tt.todo.ID)

			// 結果を確認する（異常系）
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Equal(t, tt.err, err)
				}
				return
			}

			// 結果を確認する（正常系）
			if assert.NoError(t, err) {
				todo, err := todoRepository.FindById(tt.todo.ID)
				assert.Nil(t, todo)
				if assert.Error(t, err) {
					assert.Equal(t, errors.New("record not found"), err)
				}
			}
		})
	}
}
