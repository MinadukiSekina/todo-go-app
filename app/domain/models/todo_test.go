package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type args struct {
	target     string
	correspond map[string]Status
}

var correspond = map[string]Status{
	"NotStarted": NotStarted,
	"Done":       Done,
}

func TestStrToStatus(t *testing.T) {
	cases := map[string]struct {
		args       args
		want       Status
		expectErr  bool
		errMessage string
	}{
		"正常ケース：NotStarted": {
			args:       args{target: "NotStarted", correspond: correspond},
			want:       NotStarted,
			expectErr:  false,
			errMessage: "",
		},
		"正常ケース：Done": {
			args:       args{target: "Done", correspond: correspond},
			want:       Done,
			expectErr:  false,
			errMessage: "",
		},
		"異常ケース：対象の文字列が空文字列": {
			args:       args{target: "", correspond: correspond},
			want:       invalid,
			expectErr:  true,
			errMessage: "target string is empty",
		},
		"異常ケース：対象の文字列が変換マップに無い": {
			args:       args{target: "test", correspond: correspond},
			want:       invalid,
			expectErr:  true,
			errMessage: "invalid status value: test",
		},
		"異常ケース：変換マップが空": {
			args:       args{target: "NotStarted", correspond: map[string]Status{}},
			want:       invalid,
			expectErr:  true,
			errMessage: "correspond map is empty",
		},
		"異常ケース：変換後の値が列挙定数の範囲を超えている(負の数)": {
			args:       args{target: "NotStarted", correspond: map[string]Status{"NotStarted": -1}},
			want:       invalid,
			expectErr:  true,
			errMessage: "invalid status value: NotStarted",
		},
		"異常ケース：変換後の値が列挙定数の範囲を超えている(正の数)": {
			args:       args{target: "NotStarted", correspond: map[string]Status{"NotStarted": 99999999999999}},
			want:       invalid,
			expectErr:  true,
			errMessage: "invalid status value: NotStarted",
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := StrToStatus(tt.args.target, tt.args.correspond)

			// check error
			if tt.expectErr {
				if err == nil {
					t.Error("No error is occured.")
				} else {
					assert.Equal(t, tt.errMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// check result
			assert.Equal(t, tt.want, result)
		})
	}
}
