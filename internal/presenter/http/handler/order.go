package handler

// POST /api/user/register — регистрация пользователя;
// POST /api/user/login — аутентификация пользователя;
// POST /api/user/order — загрузка пользователем номера заказа для расчёта;
// GET /api/user/order — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
// GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.