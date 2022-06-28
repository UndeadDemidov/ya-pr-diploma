package auth

import (
	"context"
	"errors"
	"testing"

	mock_auth "github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth/mocks"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	mock_user "github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user/mocks"
	"github.com/golang/mock/gomock"
)

var errDummy = errors.New("dummy error")

func TestService_SignIn(t *testing.T) {
	type fields struct {
		userSvc *mock_user.MockRegisterer
		credMan *mock_auth.MockCredentialManager
	}
	type args struct {
		login string
		pword string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "login is in use already",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.credMan.EXPECT().GetUser(context.Background(), gomock.Any()).Return(user.User{}, nil),
				)
			},
			args: args{
				login: "test",
				pword: "test",
			},
			wantErr: true,
		},
		{
			name: "can't register new user",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.credMan.EXPECT().GetUser(context.Background(), gomock.Any()).Return(user.User{}, errDummy),
					f.userSvc.EXPECT().RegisterNewUser(context.Background(), gomock.Any()).Return(errDummy),
				)
			},
			args: args{
				login: "test",
				pword: "test",
			},
			wantErr: true,
		},
		{
			name: "can't add credentials for new user",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.credMan.EXPECT().GetUser(context.Background(), gomock.Any()).Return(user.User{}, errDummy),
					f.userSvc.EXPECT().RegisterNewUser(context.Background(), gomock.Any()).Return(nil),
					f.credMan.EXPECT().AddNewUser(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errDummy),
				)
			},
			args: args{
				login: "test",
				pword: "test",
			},
			wantErr: true,
		},
		{
			name: "added new user",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.credMan.EXPECT().GetUser(context.Background(), gomock.Any()).Return(user.User{}, errDummy),
					f.userSvc.EXPECT().RegisterNewUser(context.Background(), gomock.Any()).Return(nil),
					f.credMan.EXPECT().AddNewUser(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				login: "test",
				pword: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			reg := mock_user.NewMockRegisterer(mockCtrl)
			man := mock_auth.NewMockCredentialManager(mockCtrl)

			f := fields{
				userSvc: reg,
				credMan: man,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := NewService(reg, man)
			if err := s.SignIn(context.Background(), tt.args.login, tt.args.pword); (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
