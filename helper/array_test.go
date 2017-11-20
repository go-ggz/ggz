package helper

import (
	"reflect"
	"testing"
)

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
			name: "test in array",
			args: args{
				needle:   "a",
				haystack: []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "test not in array",
			args: args{
				needle:   "d",
				haystack: []string{"a", "b", "c"},
			},
			want: false,
		},
		{
			name: "test empty target array",
			args: args{
				needle:   "d",
				haystack: []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if _, got := InArray(tt.args.needle, tt.args.haystack); got != tt.want {
			t.Errorf("%q. InArray() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDiffArray(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test not in array",
			args: args{
				a: []string{"d"},
				b: []string{"a", "b", "c"},
			},
			want: []string{"d", "a", "b", "c"},
		},
		{
			name: "test partial not in array",
			args: args{
				a: []string{"a", "c"},
				b: []string{"a", "b", "c"},
			},
			want: []string{"b"},
		},
		{
			name: "test all match in array",
			args: args{
				a: []string{"a", "c", "b"},
				b: []string{"a", "b", "c"},
			},
			want: []string{},
		},
		{
			name: "test empty source in array",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: []string{},
		},
		{
			name: "test source len > target len",
			args: args{
				a: []string{"a", "b", "c", "d", "e"},
				b: []string{"a", "c"},
			},
			want: []string{"b", "d", "e"},
		},
	}
	for _, tt := range tests {
		if got := DiffArray(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DiffArray() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// from array_test.go
func BenchmarkDiffArray2(b *testing.B) {
	// run the DiffArray function b.N times
	for n := 0; n < b.N; n++ {
		DiffArray([]string{"a", "c"}, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
	}
}

func BenchmarkDiffArray5(b *testing.B) {
	// run the DiffArray function b.N times
	for n := 0; n < b.N; n++ {
		DiffArray([]string{"a", "c", "e", "g", "i"}, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
	}
}
