package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_app "github.com/UndeadDemidov/ya-pr-diploma/internal/app/mocks"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	midware "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var errDummy = errors.New("dummy error")

func TestAuth_RegisterUser(t *testing.T) {
	type fields struct {
		auth *mock_app.MockAuthenticator
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
			name: "status 400 empty body",
			prepare: func(f *fields) {
			},
			args: args{
				request:     ``,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
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
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockAuth := mock_app.NewMockAuthenticator(mockCtrl)

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

			auth := NewAuth(mockAuth, midware.NewDefaultSessions())
			auth.RegisterUser(w, request)
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func TestAuth_LoginUser(t *testing.T) {
	type fields struct {
		auth *mock_app.MockAuthenticator
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
					f.auth.EXPECT().Login(context.Background(), gomock.Any(), gomock.Any()).Return(user.User{ID: "1"}, nil),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusOK,
		},
		{
			name: "status 400 empty body",
			prepare: func(f *fields) {
			},
			args: args{
				request:     ``,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
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
			name: "status 401",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auth.EXPECT().Login(context.Background(), gomock.Any(), gomock.Any()).Return(user.User{}, errors2.ErrPairLoginPwordIsNotExist),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "status 500",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auth.EXPECT().Login(context.Background(), gomock.Any(), gomock.Any()).Return(user.User{}, errDummy),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusInternalServerError,
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockAuth := mock_app.NewMockAuthenticator(mockCtrl)

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

			auth := NewAuth(mockAuth, midware.NewDefaultSessions())
			auth.LoginUser(w, request)
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}
