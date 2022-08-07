-- DATABASE_URI=user=postgres password=postgres dbname=ya_pract sslmode=disable
DROP TRIGGER IF EXISTS set_timestamp ON orders;
DROP FUNCTION IF EXISTS trigger_set_timestamp;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS order_status;
