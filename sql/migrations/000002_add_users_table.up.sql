CREATE TABLE "usuarios" (
  "id" serial PRIMARY KEY,
  "nome" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "senha_hash" varchar NOT NULL
);
