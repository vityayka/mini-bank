CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "code" varchar NOT NULL,
  "is_used" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL
);


ALTER TABLE "verify_emails" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD "is_verified" BOOLEAN DEFAULT false;