-- create "tasks" table
CREATE TABLE "public"."tasks" (
  "id" text NOT NULL,
  "title" text NULL,
  "description" text NULL,
  "status" text NULL,
  "user_id" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
