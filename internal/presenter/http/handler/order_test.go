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
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/order"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder_UploadOrder(t *testing.T) {
	type fields struct {
		processor *mock.MockOrderProcessor
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
			name:    "no content uploaded",
			prepare: nil,
			args: args{
				request:     "",
				reference:   "",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "invalid content type",
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
				request:     "1",
				reference:   "",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "order already uploaded",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors2.ErrOrderAlreadyUploaded),
				)
			},
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusOK,
		},
		{
			name: "order already uploaded by another user",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors2.ErrOrderAlreadyUploadedByAnotherUser),
				)
			},
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusConflict,
		},
		{
			name: "invalid number format",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors2.ErrOrderInvalidNumberFormat),
				)
			},
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "unexpected error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(errDummy),
				)
			},
			args: args{
				request:     "1",
				reference:   "1",
				contentType: utils.ContentTypeText,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "everything is good",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				contentType: utils.ContentTypeText,
				request:     "1",
				reference:   "1",
			},
			want: http.StatusAccepted,
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockOrder := mock.NewMockOrderProcessor(mockCtrl)

			f := fields{
				processor: mockOrder,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader(tt.args.request)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			request.Header.Set(utils.ContentTypeKey, tt.args.contentType)
			w := httptest.NewRecorder()

			ord := NewOrder(mockOrder)
			ctx := context.WithValue(request.Context(), middleware.ContextUserIDKey, tt.args.reference)
			ord.UploadOrder(w, request.WithContext(ctx))
			result := w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func TestOrder_DownloadOrders(t *testing.T) {
	type fields struct {
		orders    []order.Order
		processor *mock.MockOrderProcessor
	}
	type args struct {
		orders    []order.Order
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
				orders:    nil,
				json:      "",
				reference: "",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "unexpected error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.orders, errDummy),
				)
			},
			args: args{
				orders:    nil,
				json:      "",
				reference: "1",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "empty result",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.orders, nil),
				)
			},
			args: args{
				orders:    make([]order.Order, 0),
				json:      "",
				reference: "1",
			},
			want: http.StatusNoContent,
		},
		{
			name: "set of valid results",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.processor.EXPECT().List(gomock.Any(), gomock.Any()).Return(f.orders, nil),
				)
			},
			args: args{
				orders: []order.Order{
					{
						ID:        "1",
						User:      user.User{ID: "1"},
						Number:    9278923470,
						Status:    order.Processed,
						Accrual:   50000,
						Unloaded:  utils.TimeRFC3339ParseHelper("2020-12-10T15:15:45+03:00"),
						Processed: time.Now(),
					},
					{
						ID:        "1",
						User:      user.User{ID: "1"},
						Number:    12345678903,
						Status:    order.Processing,
						Accrual:   0,
						Unloaded:  utils.TimeRFC3339ParseHelper("2020-12-10T15:12:01+03:00"),
						Processed: time.Now(),
					},
					{
						ID:        "1",
						User:      user.User{ID: "1"},
						Number:    346436439,
						Status:    order.Invalid,
						Accrual:   0,
						Unloaded:  utils.TimeRFC3339ParseHelper("2020-12-09T16:09:53+03:00"),
						Processed: time.Now(),
					},
				},
				json: `
[
	{
        "number": "9278923470",
        "status": "PROCESSED",
        "accrual": 500,
        "uploaded_at": "2020-12-10T15:15:45+03:00"
    },
    {
        "number": "12345678903",
        "status": "PROCESSING",
        "uploaded_at": "2020-12-10T15:12:01+03:00"
    },
    {
        "number": "346436439",
        "status": "INVALID",
        "uploaded_at": "2020-12-09T16:09:53+03:00"
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
			mockOrder := mock.NewMockOrderProcessor(mockCtrl)

			f := fields{
				orders:    tt.args.orders,
				processor: mockOrder,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader("")
			request := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()

			ord := NewOrder(mockOrder)
			ctx := context.WithValue(request.Context(), middleware.ContextUserIDKey, tt.args.reference)
			ord.DownloadOrders(w, request.WithContext(ctx))
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
