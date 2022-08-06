package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock "github.com/UndeadDemidov/ya-pr-diploma/internal/app/mocks"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/entity"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalance_Get(t *testing.T) {
	type fields struct {
		bal    entity.Balance
		getter *mock.MockBalanceGetter
	}
	type args struct {
		bal       entity.Balance
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
				bal:       entity.Balance{},
				json:      "",
				reference: "",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "unexpected error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.getter.EXPECT().Get(gomock.Any(), gomock.Any()).Return(f.bal, errDummy),
				)
			},
			args: args{
				bal:       entity.Balance{},
				json:      "",
				reference: "1",
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "empty result",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.getter.EXPECT().Get(gomock.Any(), gomock.Any()).Return(f.bal, nil),
				)
			},
			args: args{
				bal:       entity.Balance{},
				json:      `{"current":0}`,
				reference: "1",
			},
			want: http.StatusOK,
		},
		{
			name: "non empty result",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.getter.EXPECT().Get(gomock.Any(), gomock.Any()).Return(f.bal, nil),
				)
			},
			args: args{
				bal:       entity.Balance{Current: 200, Withdrawn: 50},
				json:      `{"current":2,"withdrawn":0.50}`,
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
			mockBalance := mock.NewMockBalanceGetter(mockCtrl)

			f := fields{
				bal:    tt.args.bal,
				getter: mockBalance,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			reader := strings.NewReader("")
			request := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()

			ord := NewBalance(mockBalance)
			ctx := context.WithValue(request.Context(), middleware.ContextUserIDKey, tt.args.reference)
			ord.Get(w, request.WithContext(ctx))
			result := w.Result()
			require.Equal(t, tt.want, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				b, _ := io.ReadAll(result.Body)
				assert.JSONEq(t, tt.args.json, string(b))
			}
		})
	}
}
