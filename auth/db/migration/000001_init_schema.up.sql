CREATE TABLE "users"(
	"username" varchar PRIMARY KEY,
	"name" varchar NOT NULL,
	"password" varchar NOT NULL,
	"about" text NOT NULL,
	"photo" varchar NOT NULL,
	"role" varchar NOT NULL,
	"created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "roles"(
	"role" varchar PRIMARY KEY
);

ALTER TABLE "users" ADD FOREIGN KEY ("role") REFERENCES "roles" ("role");
