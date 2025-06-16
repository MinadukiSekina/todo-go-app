package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
