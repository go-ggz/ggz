package helper

import "testing"

func TestIsFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "file match",
			args: args{
				filePath: "file.go",
			},
			want: true,
		},
		{
			name: "file not found",
			args: args{
				filePath: "file1.go",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.args.filePath); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
