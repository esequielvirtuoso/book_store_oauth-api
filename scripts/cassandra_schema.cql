CREATE KEYSPACE IF NOT EXISTS oauth WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};
CREATE TABLE IF NOT EXISTS oauth.access_tokens(access_token varchar PRIMARY KEY, user_id bigint, client_id bigint, expires bigint);
