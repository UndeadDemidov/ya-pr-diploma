package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAuth_RegisterUser(t *testing.T) {
	type fields struct {
		ctrl        *gomock.Controller
		mockManager *MockCredentialManager
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
					f.mockManager.EXPECT().New(gomock.Any(), gomock.Any()).Return(NewMockCredentialValidator(f.ctrl), true),
					f.mockManager.EXPECT().Add(gomock.Any()).Return(nil),
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
					f.mockManager.EXPECT().New(gomock.Any(), gomock.Any()).Return(NewMockCredentialValidator(f.ctrl), false),
				)
			},
			args: args{
				request:     `{"login": "test","password": "test"}`,
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockManager := NewMockCredentialManager(mockCtrl)

			f := fields{
				ctrl:        mockCtrl,
				mockManager: mockManager,
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.args.request)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			request.Header.Set(utils.ContentTypeKey, tt.args.contentType)
			w := httptest.NewRecorder()

			auth := NewAuth(mockManager)
			auth.RegisterUser(w, request)
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}
