package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

// POST /api/user/order — загрузка пользователем номера заказа для расчёта;
// GET /api/user/order — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;

var (
	ErrProperOrderNumberIsExpected = errors.New("proper order number is expected")
	ErrNoContentToReturn           = errors.New("no content to return")
)

type Order struct {
	processor app.OrderProcessor
}

func NewOrder(proc app.OrderProcessor) *Order {
	if proc == nil {
		panic("missing app.OrderProcessor, parameter must not be nil")
	}
	return &Order{processor: proc}
}

// UploadOrder
// Content-Type: text/plain
// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 409 — номер заказа уже был загружен другим пользователем;
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func (o *Order) UploadOrder(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		utils.ServerError(w, ErrProperOrderNumberIsExpected, http.StatusBadRequest)
		return
	}
	if r.Header.Get(utils.ContentTypeKey) != utils.ContentTypeText {
		utils.ServerError(w, ErrInvalidContentType, http.StatusBadRequest)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	usr := GetUserFromContext(r.Context())
	if usr.ID == "" {
		utils.InternalServerError(w, errors2.ErrSessionUserCanNotBeDefined)
		return
	}
	err = o.processor.Add(r.Context(), usr, string(b))
	switch err {
	case errors2.ErrOrderAlreadyUploaded:
		w.WriteHeader(http.StatusOK)
		return
	case errors2.ErrOrderAlreadyUploadedByAnotherUser:
		utils.ServerError(w, err, http.StatusConflict)
		return
	case errors2.ErrOrderInvalidNumberFormat:
		utils.ServerError(w, err, http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// DownloadOrders
// Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.
// 200 — успешная обработка запроса.
// 204 — нет данных для ответа.
// 500 — внутренняя ошибка сервера.
func (o *Order) DownloadOrders(w http.ResponseWriter, r *http.Request) {
	usr := GetUserFromContext(r.Context())
	if usr.ID == "" {
		utils.InternalServerError(w, errors2.ErrSessionUserCanNotBeDefined)
		return
	}
	list, err := o.processor.List(r.Context(), usr)
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
