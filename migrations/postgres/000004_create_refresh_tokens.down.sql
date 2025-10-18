-- Drop refresh_tokens table and related objects
DROP FUNCTION IF EXISTS cleanup_expired_tokens();
DROP TABLE IF EXISTS refresh_tokens;

