DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users"
(
    "id"         varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
    "first_name" varchar(64) COLLATE "pg_catalog"."default",
    "last_name"  varchar(64) COLLATE "pg_catalog"."default",
    "nickname"   varchar(64) COLLATE "pg_catalog"."default",
    "password"   varchar(64) COLLATE "pg_catalog"."default",
    "email"      varchar(64) COLLATE "pg_catalog"."default",
    "country"    varchar(64) COLLATE "pg_catalog"."default",
    "created_at" timestamp(6),
    "updated_at" timestamp(6)
);

ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");