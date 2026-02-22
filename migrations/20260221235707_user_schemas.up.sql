-- create "users" table
CREATE TABLE "public"."users" (
  "id" text NOT NULL,
  "username" text NULL,
  "email" text NULL,
  "password" text NULL,
  "role" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
