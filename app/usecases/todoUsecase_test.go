package usecases

import (
	"errors"
	"testing"

	"github.com/MinadukiSekina/todo-go-app/app/domain/models"
	mock_repository "github.com/MinadukiSekina/todo-go-app/app/mock/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSearchByID(t *testing.T) {

	type args struct {
		ID uint
	}

	testTodo := models.Todo{Title: "test", Status: models.NotStarted}
	testTodo.ID = 1

	cases := map[string]struct {
		args      args
		want      *models.Todo
		expectErr bool
		err       error
	}{
		"正常ケース:データあり": {
			args:      args{ID: testTodo.ID},
			want:      &testTodo,
			expectErr: false,
			err:       nil,
		},
		"異常ケース:データなし": {
			args:      args{ID: 2},
			want:      nil,
			expectErr: true,
			err:       errors.New("Record Not found"),
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_repository.NewMockTodoRepository(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			// 違う引数で呼び出すとエラーになるらしい
			mock.EXPECT().FindById(tt.args.ID).Return(tt.want, tt.err)

			// mockを利用してテストする
			Usecase := NewTodoUsecase(mock)
			result, err := Usecase.SearchByID(tt.args.ID)

			// 結果を確認
			assert.Equal(t, tt.want, result)

			// エラーの確認
			if tt.expectErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShow(t *testing.T) {

	todo1 := models.Todo{Title: "test1", Status: models.NotStarted}
	todo1.ID = 1

	todo2 := models.Todo{Title: "test2", Status: models.NotStarted}
	todo2.ID = 1

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
		"正常ケース:1つのデータあり": {
			want:      &onlyOneTodos,
			expectErr: false,
			err:       nil,
		},
		"正常ケース:2つのデータあり": {
			want:      &manyHasTodos,
			expectErr: false,
			err:       nil,
		},
		"異常ケース:エラーあり": {
			want:      nil,
			expectErr: true,
			err:       errors.New("Some error is occured"),
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_repository.NewMockTodoRepository(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			// 違う引数で呼び出すとエラーになるらしい
			mock.EXPECT().FindAll().Return(tt.want, tt.err)

			// mockを利用してテストする
			Usecase := NewTodoUsecase(mock)
			result, err := Usecase.Show()

			// 結果を確認
			assert.Equal(t, tt.want, result)

			// エラーの確認
			if tt.expectErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestAdd(t *testing.T) {

	type args struct {
		todo *models.Todo
	}

	testTodo := models.Todo{Title: "test", Status: models.NotStarted}

	cases := map[string]struct {
		args      args
		expectErr bool
		err       error
	}{
		"正常ケース:登録完了": {
			args:      args{&testTodo},
			expectErr: false,
			err:       nil,
		},
		"異常ケース:登録失敗": {
			args:      args{},
			expectErr: true,
			err:       errors.New("Save todo is failed"),
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_repository.NewMockTodoRepository(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			// 違う引数で呼び出すとエラーになるらしい
			mock.EXPECT().Create(tt.args.todo).Return(tt.err)

			// mockを利用してテストする
			Usecase := NewTodoUsecase(mock)
			err := Usecase.Add(tt.args.todo)

			// 結果を確認
			if tt.expectErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEdit(t *testing.T) {

	type args struct {
		todo *models.Todo
	}

	testTodo := models.Todo{Title: "test", Status: models.NotStarted}

	cases := map[string]struct {
		args      args
		expectErr bool
		err       error
	}{
		"正常ケース:更新完了": {
			args:      args{&testTodo},
			expectErr: false,
			err:       nil,
		},
		"異常ケース:更新失敗": {
			args:      args{&testTodo},
			expectErr: true,
			err:       errors.New("Save todo is failed"),
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_repository.NewMockTodoRepository(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			// 違う引数で呼び出すとエラーになるらしい
			mock.EXPECT().Update(tt.args.todo).Return(tt.err)

			// mockを利用してテストする
			Usecase := NewTodoUsecase(mock)
			err := Usecase.Edit(tt.args.todo)

			// 結果を確認
			if tt.expectErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	type args struct {
		ID uint
	}

	cases := map[string]struct {
		args      args
		expectErr bool
		err       error
	}{
		"正常ケース:更新完了": {
			args:      args{ID: 1},
			expectErr: false,
			err:       nil,
		},
		"異常ケース:更新失敗": {
			args:      args{ID: 1},
			expectErr: true,
			err:       errors.New("Save todo is failed"),
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			// モックの呼び出しを管理するControllerを生成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// モックの生成
			mock := mock_repository.NewMockTodoRepository(mockCtrl)
			// テスト中に呼ばれるべき関数と帰り値を指定
			// 違う引数で呼び出すとエラーになるらしい
			mock.EXPECT().Delete(tt.args.ID).Return(tt.err)

			// mockを利用してテストする
			Usecase := NewTodoUsecase(mock)
			err := Usecase.Delete(tt.args.ID)

			// 結果を確認
			if tt.expectErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
