package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock "github.com/UndeadDemidov/ya-pr-diploma/internal/app/mocks"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestWithdrawal_Register(t *testing.T) {
	type fields struct {
		processor *mock.MockWithdrawalProcessor
	}
	type args struct {
		contentType string
		request     string
		reference   string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    int
	}{
		{
			name:    "no content passed",
			prepare: nil,
			args: args{
				request:     "",
				reference:   "",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "invalid content type",
			prepare: nil,
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "invalid content",
			prepare: nil,
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "session error",
			prepare: nil,
			args: args{
				request: `
{
	"order": "2377225624",
    "sum": 751
}`,
				reference:   "",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "not enough fund",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors2.ErrWithdrawalNotEnoughFund),
				)
			},
			args: args{
				request: `
{
	"order": "2377225624",
    "sum": 751
}`,
				reference:   "1",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusPaymentRequired,
		},
		{
			name: "not enough fund",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors2.ErrOrderInvalidNumberFormat),
				)
			},
			args: args{
				request: `
{
	"order": "2377225624",
    "sum": 751
}`,
				reference:   "1",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "not enough fund",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				request: `
{
	"order": "2377225624",
    "sum": 751
}`,
				reference:   "1",
				contentType: utils.ContentTypeJSON,
			},
			want: http.StatusOK,
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockWtdrwl := mock.NewMockWithdrawalProcessor(mockCtrl)

			f := fields{
				processor: mockWtdrwl,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.args.request)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			request.Header.Set(utils.ContentTypeKey, tt.args.contentType)
			w := httptest.NewRecorder()

			wtdrwl := NewWithdrawal(mockWtdrwl)
			ctx := context.WithValue(request.Context(), middleware.ContextUserIDKey, tt.args.reference)
			wtdrwl.Register(w, request.WithContext(ctx))
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}
