package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

type Withdrawal struct {
	processor app.WithdrawalProcessor
}

func NewWithdrawal(processor app.WithdrawalProcessor) *Withdrawal {
	if processor == nil {
		panic("missing app.WithdrawalProcessor, parameter must not be nil")
	}
	return &Withdrawal{processor: processor}
}

// CashOut
// 200 — успешная обработка запроса;
// 402 — на счету недостаточно средств;
// 422 — неверный номер заказа;
// 500 — внутренняя ошибка сервера.
func (wd Withdrawal) CashOut(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		utils.ServerError(w, ErrProperJSONIsExpected, http.StatusBadRequest)
		return
	}
	if r.Header.Get(utils.ContentTypeKey) != utils.ContentTypeJSON {
		utils.ServerError(w, ErrInvalidContentType, http.StatusBadRequest)
		return
	}

	var req wtdrwlRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ServerError(w, ErrProperJSONIsExpected, http.StatusBadRequest)
		return
	}

	usr := GetUserFromContext(r.Context())
	if usr.ID == "" {
		utils.InternalServerError(w, errors2.ErrSessionUserCanNotBeDefined)
		return
	}

	err = wd.processor.Add(r.Context(), usr, req.Order, req.Sum)
	switch {
	case errors.Is(err, errors2.ErrWithdrawalNotEnoughFund):
		utils.ServerError(w, err, http.StatusPaymentRequired)
		return
	case errors.Is(err, errors2.ErrOrderInvalidNumberFormat):
		utils.ServerError(w, err, http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// History
// 200 — успешная обработка запроса.
// 204 — нет данных для ответа.
// 500 — внутренняя ошибка сервера.
func (wd Withdrawal) History(w http.ResponseWriter, r *http.Request) {
	usr := GetUserFromContext(r.Context())
	if usr.ID == "" {
		utils.InternalServerError(w, errors2.ErrSessionUserCanNotBeDefined)
		return
	}
	list, err := wd.processor.List(r.Context(), usr)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if len(list) == 0 {
		utils.ServerError(w, ErrNoContentToReturn, http.StatusNoContent)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// ToDo было бы не плохо вставить адаптер из list в response
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
}

// {
// "order": "2377225624",
// "sum": 751
// }
type wtdrwlRequest struct {
	Order string          `json:"order"`
	Sum   primit.Currency `json:"sum"`
}
