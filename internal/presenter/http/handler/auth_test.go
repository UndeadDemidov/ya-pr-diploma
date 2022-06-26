package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler/mocks"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errDummy = errors.New("dummy error")

func TestAuth_RegisterUser(t *testing.T) {
	type fields struct {
		auth *mock_handler.MockAuthenticator
	}
	type args struct {
		request     string
		contentType string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    int
	}{
		{
			name: "status 200",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auth.EXPECT().SignIn(context.Background(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusOK,
		},
		{
			name: "status 400 invalid json",
			prepare: func(f *fields) {
			},
			args: args{
				request:     `{"login": "test","password": "test"`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "status 400 invalid content type",
			prepare: func(f *fields) {
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: "xxx",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "status 409",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auth.EXPECT().SignIn(context.Background(), gomock.Any(), gomock.Any()).Return(errors2.ErrLoginIsInUseAlready),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusConflict,
		},
		{
			name: "status 500",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auth.EXPECT().SignIn(context.Background(), gomock.Any(), gomock.Any()).Return(errDummy),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockAuth := mock_handler.NewMockAuthenticator(mockCtrl)

			f := fields{
				auth: mockAuth,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.args.request)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			request.Header.Set(utils.ContentTypeKey, tt.args.contentType)
			w := httptest.NewRecorder()

			auth := NewAuth(mockAuth)
			auth.RegisterUser(w, request)
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}
