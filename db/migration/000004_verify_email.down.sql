DROP TABLE IF EXISTS "verify_emails";

ALTER TABLE IF EXISTS "users" DROP COLUMN IF EXISTS "is_verified";
