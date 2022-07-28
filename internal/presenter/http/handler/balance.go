package handler

import (
	"encoding/json"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
// GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.

type Balance struct {
	getter app.BalanceGetter
}

func NewBalance(getter app.BalanceGetter) *Balance {
	if getter == nil {
		panic("missing app.BalanceGetter, parameter must not be nil")
	}
	return &Balance{getter: getter}
}

// Get
// 200 — успешная обработка запроса.
// 500 — внутренняя ошибка сервера.
func (b Balance) Get(w http.ResponseWriter, r *http.Request) {
	usr := GetUserFromContext(r.Context())
	if usr.ID == "" {
		utils.InternalServerError(w, errors2.ErrSessionUserCanNotBeDefined)
		return
	}
	bal, err := b.getter.Get(r.Context(), usr)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// ToDo было бы не плохо вставить адаптер из balance в response
	err = json.NewEncoder(w).Encode(bal)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
}
