package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	type args struct {
		t   Type
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Internal check",
			args: args{Internal, EInternal("Internal error", nil)},
			want: true,
		},
		{
			name: "NotFound check",
			args: args{NotFound, ENotFound("NotFound error", nil)},
			want: true,
		},
		{
			name: "BadRequest check",
			args: args{BadRequest, EBadRequest("BadRequest error", nil)},
			want: true,
		},
		{
			name: "Validation check",
			args: args{Validation, EValidation("Validation error", nil)},
			want: true,
		},
		{
			name: "EAlreadyExists check",
			args: args{AlreadyExists, EAlreadyExists("EAlreadyExists error", nil)},
			want: true,
		},
		{
			name: "EUnauthorized check",
			args: args{Unauthorized, EUnauthorized("EUnauthorized error", nil)},
			want: true,
		},
		{
			name: "new error check",
			args: args{Internal, New(Internal, "error", nil)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.t, tt.args.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorType(t *testing.T) {
	t.Run("EValidation", func(t *testing.T) {
		err := EValidation("not validation", errors.New("test"))
		assert.Equal(t, "not validation", err.Error())
		assert.Equal(t, "test", err.(*Error).Cause())
		assert.Equal(t, http.StatusForbidden, err.(*Error).Type.Code())
		assert.Equal(t, "Validation error", err.(*Error).Type.String())
	})

	t.Run("EBadRequest", func(t *testing.T) {
		err := EBadRequest("not EBadRequest", nil)
		assert.Equal(t, "not EBadRequest", err.Error())
		assert.Equal(t, "", err.(*Error).Cause())
		assert.Equal(t, http.StatusBadRequest, err.(*Error).Type.Code())
		assert.Equal(t, "BadRequest error", err.(*Error).Type.String())
	})

	t.Run("EAlreadyExists", func(t *testing.T) {
		err := EAlreadyExists("not EAlreadyExists", errors.New("test"))
		assert.Equal(t, "not EAlreadyExists", err.Error())
		assert.Equal(t, "test", err.(*Error).Cause())
		assert.Equal(t, http.StatusBadRequest, err.(*Error).Type.Code())
		assert.Equal(t, "Item already exists", err.(*Error).Type.String())
	})

	t.Run("EInternal", func(t *testing.T) {
		err := EInternal("not EInternal", errors.New("test"))
		assert.Equal(t, "not EInternal", err.Error())
		assert.Equal(t, "test", err.(*Error).Cause())
		assert.Equal(t, http.StatusInternalServerError, err.(*Error).Type.Code())
		assert.Equal(t, "Internal Error", err.(*Error).Type.String())
	})

	t.Run("ENotFound", func(t *testing.T) {
		err := ENotFound("not validation", errors.New("test"))
		assert.Equal(t, "not validation", err.Error())
		assert.Equal(t, "test", err.(*Error).Cause())
		assert.Equal(t, http.StatusNotFound, err.(*Error).Type.Code())
		assert.Equal(t, "Item not found", err.(*Error).Type.String())
	})

	t.Run("EUnauthorized", func(t *testing.T) {
		err := EUnauthorized("not validation", errors.New("test"))
		assert.Equal(t, "not validation", err.Error())
		assert.Equal(t, "test", err.(*Error).Cause())
		assert.Equal(t, http.StatusUnauthorized, err.(*Error).Type.Code())
		assert.Equal(t, "Unauthorized error", err.(*Error).Type.String())
	})

	t.Run("Unknown Error", func(t *testing.T) {
		unknown := Type("Unknown")
		assert.Equal(t, http.StatusInternalServerError, unknown.Code())
		assert.Equal(t, "Unknown error", unknown.String())
	})
}

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "Internal",
			t:    Internal,
			want: "Internal Error",
		},
		{
			name: "NotFound",
			t:    NotFound,
			want: "Item not found",
		},
		{
			name: "AlreadyExists",
			t:    AlreadyExists,
			want: "Item already exists",
		},
		{
			name: "Validation",
			t:    Validation,
			want: "Validation error",
		},
		{
			name: "BadRequest",
			t:    BadRequest,
			want: "BadRequest error",
		},
		{
			name: "Unauthorized",
			t:    Unauthorized,
			want: "Unauthorized error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
