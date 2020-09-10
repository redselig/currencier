DROP DATABASE IF EXISTS currencier;
CREATE DATABASE currencier;
CREATE USER igor WITH encrypted password 'igor';
GRANT ALL PRIVILEGES ON DATABASE currencier to igor;
