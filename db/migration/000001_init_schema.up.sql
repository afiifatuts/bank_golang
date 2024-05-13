CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" numeric NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE
);

CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
