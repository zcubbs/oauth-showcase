package internal

import (
	"strings"
	"testing"
)

func TestGetJSONString(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetJSONString",
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
			},
			want: `{
				"sub": "test",
				"name": "John Doe",
				"admin": true
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := strings.ReplaceAll(tt.want, "\t", "")
			want = strings.ReplaceAll(want, " ", "")
			got := GetJSONString(tt.args.token)
			got = strings.ReplaceAll(got, "\t", "")
			got = strings.ReplaceAll(got, " ", "")

			if got != want {
				t.Errorf("GetJSONString() = %v,\n want:\n %v", got, want)
			}
		})
	}
}
