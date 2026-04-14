CREATE TABLE IF NOT EXISTS "status_task" (
  "id" serial PRIMARY KEY,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE IF NOT EXISTS "type_task" (
  "id" serial PRIMARY KEY,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE IF NOT EXISTS "produto" (
  "id" serial PRIMARY KEY,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE IF NOT EXISTS "escopo" (
  "id" serial PRIMARY KEY,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE IF NOT EXISTS "cliente" (
  "id" serial PRIMARY KEY,
  "nome" varchar,
  "ativo" bool
);

CREATE TABLE IF NOT EXISTS "ticket" (
  "id" serial PRIMARY KEY,
  "descricao" varchar,
  "aberto_em" timestamp,
  "id_status" integer,
  "id_epic" integer,
  "departamento" varchar,
  "categoria" varchar,
  "responsavel" varchar,
  FOREIGN KEY ("id_status") REFERENCES "status_task" ("id")
);

CREATE TABLE IF NOT EXISTS "task" (
  "id" serial PRIMARY KEY,
  "task" integer,
  "id_type_task" integer,
  "descricao" varchar,
  "id_status" integer,
  "id_produto" integer,
  "id_escopo" integer,
  "id_cliente" integer,
  FOREIGN KEY ("id_type_task") REFERENCES "type_task" ("id"),
  FOREIGN KEY ("id_status") REFERENCES "status_task" ("id"),
  FOREIGN KEY ("id_produto") REFERENCES "produto" ("id"),
  FOREIGN KEY ("id_escopo") REFERENCES "escopo" ("id"),
  FOREIGN KEY ("id_cliente") REFERENCES "cliente" ("id")
);
