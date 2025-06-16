package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	mock_usecases "github.com/MinadukiSekina/todo-go-app/app/mock/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIndex(t *testing.T) {

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo1.ID = 1

	todo2 := models.Todo{Title: "test2", Status: models.NotStarted}
	todo2.ID = 1

	nothingTodos := []models.Todo{}
	onlyOneTodos := []models.Todo{todo1}
	manyHasTodos := []models.Todo{todo1, todo2}

	cases := map[string]struct {
		prepareMockFn func(m *mock_usecases.MockTodoUsecase)
		want          int
		expectErr     bool
		err           error
	}{
		"正常ケース:データなし": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) { m.EXPECT().Show().Return(&nothingTodos, nil) },
			want:          http.StatusOK,
			expectErr:     false,
			err:           nil,
		},
		"正常ケース:1件データあり": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) { m.EXPECT().Show().Return(&onlyOneTodos, nil) },
			want:          http.StatusOK,
			expectErr:     false,
			err:           nil,
		},
		"正常ケース:2件データあり": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) { m.EXPECT().Show().Return(&manyHasTodos, nil) },
			want:          http.StatusOK,
			expectErr:     false,
			err:           nil,
		},
		"異常ケース:エラーあり": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().Show().Return(nil, errors.New("something is wrong"))
			},
			want:      http.StatusInternalServerError,
			expectErr: true,
			err:       nil,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_usecases.NewMockTodoUsecase(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			tt.prepareMockFn(mock)

			// gin contextの生成
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// テンプレートの読み込み
			// route.goと同じ指定だとエラーになったため、appからのパスで指定する
			r.LoadHTMLGlob("/app/templates/*/*.html")

			// リクエストを設定
			req, _ := http.NewRequest("GET", "/todo", nil)
			c.Request = req

			// mockを利用してテストする
			handler := NewTodoHandler(mock)
			handler.Index(c)

			// 結果を確認
			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestShowByID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// テスト用の引数を格納する
	type args struct {
		id any
	}

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo1.ID = 1

	cases := map[string]struct {
		prepareMockFn func(m *mock_usecases.MockTodoUsecase)
		args          args
		want          int
		expectErr     bool
		err           error
	}{
		"正常ケース:データあり": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) { m.EXPECT().SearchByID(uint(1)).Return(&todo1, nil) },
			args:          args{id: 1},
			want:          http.StatusOK,
			expectErr:     false,
			err:           nil,
		},
		"異常ケース:IDを数値に変換できない": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				// IDが数値に変換できない場合はSearchByIDは呼ばれない
			},
			args:      args{id: "string"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
		"異常ケース:データなし": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().SearchByID(uint(1)).Return(nil, errors.New("record not found"))
			},
			args:      args{id: 1},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_usecases.NewMockTodoUsecase(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			tt.prepareMockFn(mock)

			// gin contextの生成
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// テンプレートの読み込み
			// route.goと同じ指定だとエラーになったため、appからのパスで指定する
			r.LoadHTMLGlob("/app/templates/*/*.html")

			// リクエストを設定
			req, _ := http.NewRequest("GET", fmt.Sprintf("/todo/%v", tt.args.id), nil)
			c.Request = req

			// パラメータを設定
			c.Params = []gin.Param{{Key: "id", Value: fmt.Sprint(tt.args.id)}}

			// mockを利用してテストする
			handler := NewTodoHandler(mock)
			handler.ShowById(c)

			// 結果を確認
			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestCreate(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// テスト用の引数を格納する
	type args struct {
		title string
	}

	cases := map[string]struct {
		prepareMockFn func(m *mock_usecases.MockTodoUsecase)
		args          args
		want          int
		expectErr     bool
		err           error
	}{
		"正常ケース:作成に成功": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) { m.EXPECT().Add(gomock.Any()).Return(nil) },
			args:          args{title: "test1"},
			want:          http.StatusFound,
			expectErr:     false,
			err:           nil,
		},
		"異常ケース:作成に失敗": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().Add(gomock.Any()).Return(errors.New("create todo is failed"))
			},
			args:      args{title: "failed"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_usecases.NewMockTodoUsecase(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			tt.prepareMockFn(mock)

			// gin contextの生成
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)
			// テンプレートの読み込み
			// route.goと同じ指定だとエラーになったため、appからのパスで指定する
			r.LoadHTMLGlob("/app/templates/*/*.html")

			// フォームデータの組み立て
			formData := url.Values{}
			formData.Add("title", tt.args.title)

			// リクエストを設定
			req, _ := http.NewRequest("POST", "/todo", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request = req

			// mockを利用してテストする
			handler := NewTodoHandler(mock)
			handler.Create(c)

			// GETの場合と異なり、POSTの場合はリダイレクトのステータスコードが書き込まれないらしい
			// 回避策として、c.Redirectの後に明示的に書き込むようにする（Create内部でc.Redirectを呼んでいる）
			// ref: https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for
			c.Writer.WriteHeaderNow()

			// 結果を確認
			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestUpdate(t *testing.T) {

	gin.SetMode(gin.TestMode)

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo1.ID = 1

	// テスト用の引数を格納する
	type args struct {
		id     any
		title  string
		status string
	}

	cases := map[string]struct {
		prepareMockFn func(m *mock_usecases.MockTodoUsecase)
		args          args
		want          int
		expectErr     bool
		err           error
	}{
		"正常ケース:更新に成功": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().SearchByID(uint(1)).Return(&todo1, nil)
				m.EXPECT().Edit(gomock.Any()).Return(nil)
			},
			args:      args{id: 1, title: "test1", status: "completed"},
			want:      http.StatusFound,
			expectErr: false,
			err:       nil,
		},
		"異常ケース:IDが数値に変換できない": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				// 変換できない場合はusecaseの処理が走る前にReturnするので何もしない
			},
			args:      args{id: "string", title: "failed", status: "completed"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
		"異常ケース:ステータスの値が変換できない": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				// 変換できない場合はusecaseの処理が走る前にReturnするので何もしない
			},
			args:      args{id: 1, title: "failed", status: "cannotConverted"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
		"異常ケース:対象のタスクが存在しない": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().SearchByID(uint(1)).Return(nil, errors.New("record not found"))
			},
			args:      args{id: 1, title: "failed", status: "completed"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
		"異常ケース:更新に失敗": {
			prepareMockFn: func(m *mock_usecases.MockTodoUsecase) {
				m.EXPECT().SearchByID(uint(1)).Return(&todo1, nil)
				m.EXPECT().Edit(gomock.Any()).Return(errors.New("something is wrong"))
			},
			args:      args{id: 1, title: "failed", status: "completed"},
			want:      http.StatusSeeOther,
			expectErr: true,
			err:       nil,
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_usecases.NewMockTodoUsecase(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			tt.prepareMockFn(mock)

			// gin contextの生成
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)
			// テンプレートの読み込み
			// route.goと同じ指定だとエラーになったため、appからのパスで指定する
			r.LoadHTMLGlob("/app/templates/*/*.html")

			// フォームデータの組み立て
			formData := url.Values{}
			formData.Add("title", tt.args.title)
			formData.Add("status", tt.args.status)

			// リクエストを設定
			req, _ := http.NewRequest("POST", fmt.Sprintf("/todo/%v", tt.args.id), strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request = req

			// パラメータを設定
			c.Params = []gin.Param{{Key: "id", Value: fmt.Sprint(tt.args.id)}}

			// mockを利用してテストする
			handler := NewTodoHandler(mock)
			handler.Update(c)

			// GETの場合と異なり、POSTの場合はリダイレクトのステータスコードが書き込まれないらしい
			// 回避策として、c.Redirectの後に明示的に書き込むようにする（Create内部でc.Redirectを呼んでいる）
			// ref: https://stackoverflow.com/questions/76319196/unit-testing-of-gins-context-redirect-works-for-get-response-code-but-fails-for
			c.Writer.WriteHeaderNow()

			// 結果を確認
			assert.Equal(t, tt.want, w.Code)
		})
	}
}
