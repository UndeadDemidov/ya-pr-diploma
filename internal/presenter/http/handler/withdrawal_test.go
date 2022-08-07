package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mock "github.com/UndeadDemidov/ya-pr-diploma/internal/app/mocks"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/entity"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawal_CashOut(t *testing.T) {
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
			wtdrwl.CashOut(w, request.WithContext(ctx))
			result := w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func TestWithdrawal_History(t *testing.T) {
	type fields struct {
		wtdrwls   []entity.Withdrawal
		processor *mock.MockWithdrawalProcessor
	}
	type args struct {
		wtdrwls   []entity.Withdrawal
		json      string
		reference string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    int
	}{
		{
			name:    "session error",
			prepare: nil,
			args: args{
				wtdrwls:   nil,
				json:      "",
				reference: "",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "unexpected error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.wtdrwls, errDummy),
				)
			},
			args: args{
				wtdrwls:   nil,
				json:      "",
				reference: "1",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "empty result",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.wtdrwls, nil),
				)
			},
			args: args{
				wtdrwls:   make([]entity.Withdrawal, 0),
				json:      "",
				reference: "1",
			},
			want: http.StatusNoContent,
		},
		{
			name: "set of valid results",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.wtdrwls, nil),
				)
			},
			args: args{
				wtdrwls: []entity.Withdrawal{
					{
						ID:   "1",
						User: user.User{ID: "1"},
						Order: entity.Order{
							ID:        "1",
							User:      user.User{ID: "1"},
							Number:    2377225624,
							Status:    entity.Processed,
							Accrual:   200,
							Unloaded:  time.Time{},
							Processed: time.Time{},
						},
						Sum:       50000,
						Processed: utils.TimeRFC3339ParseHelper("2020-12-09T16:09:57+03:00"),
					},
				},
				json: `
[
	{
        "order": "2377225624",
        "sum": 500,
        "processed_at": "2020-12-09T16:09:57+03:00"
    }
]`,
				reference: "1",
			},
			want: http.StatusOK,
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockWtdrwls := mock.NewMockWithdrawalProcessor(mockCtrl)

			f := fields{
				wtdrwls:   tt.args.wtdrwls,
				processor: mockWtdrwls,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader("")
			request := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()

			ord := NewWithdrawal(mockWtdrwls)
			ctx := context.WithValue(request.Context(), middleware.ContextUserIDKey, tt.args.reference)
			ord.History(w, request.WithContext(ctx))
			result := w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				b, _ := io.ReadAll(result.Body)
				assert.JSONEq(t, tt.args.json, string(b))
			}
		})
	}
}
