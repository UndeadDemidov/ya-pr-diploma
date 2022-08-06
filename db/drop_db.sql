-- DATABASE_URI=user=postgres password=postgres dbname=ya_pract sslmode=disable
DROP TYPE IF EXISTS order_status;
DROP TABLE IF EXISTS orders;
DROP FUNCTION IF EXISTS trigger_set_timestamp;

DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;