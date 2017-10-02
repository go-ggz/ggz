package helper

import "testing"

func TestInArray(t *testing.T) {
	type args struct {
		needle   string
		haystack []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"in array",
			args{
				needle:   "apple",
				haystack: []string{"apple", "banana"},
			},
			true,
		},
		{
			"not in array",
			args{
				needle:   "apple",
				haystack: []string{"appleboy", "banana"},
			},
			false,
		},
	}
	for _, tt := range tests {
		if got := InArray(tt.args.needle, tt.args.haystack); got != tt.want {
			t.Errorf("%q. InArray() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
