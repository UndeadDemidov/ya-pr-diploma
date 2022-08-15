-- DATABASE_URI postgres://postgres:postgres@localhost:5432/ya_pract?sslmode=disable
DROP TRIGGER IF EXISTS trigger_before_insert_withdrawals ON withdrawals;
DROP FUNCTION IF EXISTS update_withdrawn;

DROP TRIGGER IF EXISTS trigger_after_update_orders ON orders;
DROP FUNCTION IF EXISTS update_accrual;
DROP TRIGGER IF EXISTS set_timestamp ON orders;
DROP FUNCTION IF EXISTS trigger_set_timestamp;

DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS order_status;