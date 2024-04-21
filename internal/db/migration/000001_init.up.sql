CREATE TABLE "articles" (
  "id" bigserial PRIMARY KEY,
  "author_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "content" varchar NOT NULL,
  "rating" int NOT NULL DEFAULT -1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "authors" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "description" varchar NOT NULL,
  "rating" int NOT NULL DEFAULT -1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "articles" ("author_id");

CREATE INDEX ON "articles" ("rating");

CREATE INDEX ON "authors" ("name");

CREATE INDEX ON "authors" ("rating");

COMMENT ON COLUMN "articles"."rating" IS '-1 Stands for not-checked yet, 0-9 is actual rating';

COMMENT ON COLUMN "authors"."rating" IS '-1 Stands for not-checked yet, 0-3 is actual rating';

ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "authors" ("id");
