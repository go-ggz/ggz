package helper

import (
	"reflect"
	"testing"
)

func TestGetCacheKey(t *testing.T) {
	type args struct {
		module string
		id     interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "for string id",
			args: args{
				module: "user",
				id:     "100",
			},
			want: "user:100",
		},
		{
			name: "for int64 id",
			args: args{
				module: "user",
				id:     int64(100),
			},
			want: "user:100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCacheKey(tt.args.module, tt.args.id); got != tt.want {
				t.Errorf("GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCacheID(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "user module",
			args: args{
				key: "user:100",
			},
			want:    int64(100),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCacheID(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCacheID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCacheID() = %v, want %v", got, tt.want)
			}
		})
	}
}
