package smallinterface

import (
	"context"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

//go:generate mockgen -destination=./mocks/mock_userstore.go -package=mocks github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface UserStore

// UserStore はユーザー関連の操作のみを定義した小さなインターフェースです
// これは凝集度が高いアプローチの例です
type UserStore interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}
