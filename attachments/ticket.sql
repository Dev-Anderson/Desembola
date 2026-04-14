CREATE TABLE "ticket" (
  "id" serial,
  "descricao" varchar,
  "aberto_em" timestamp,
  "id_status" integer,
  "id_epic" integer,
  "departamento" varchar,
  "categoria" varchar,
  "responsavel" varchar
);

CREATE TABLE "status_task" (
  "id" serial,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE "task" (
  "id" serial,
  "task" integer,
  "id_type_task" integer,
  "descricao" varchar,
  "id_status" integer,
  "id_produto" integer,
  "id_escopo" integer,
  "id_cliente" integer
);

CREATE TABLE "type_task" (
  "id" serial,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE "produto" (
  "id" serial,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE "escopo" (
  "id" serial,
  "descricao" varchar,
  "ativo" bool
);

CREATE TABLE "cliente" (
  "id" serial,
  "nome" varchar,
  "ativo" bool
);

ALTER TABLE "ticket" ADD FOREIGN KEY ("id_status") REFERENCES "status_task" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "task" ADD FOREIGN KEY ("id_type_task") REFERENCES "type_task" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "task" ADD FOREIGN KEY ("id_status") REFERENCES "status_task" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "task" ADD FOREIGN KEY ("id_produto") REFERENCES "produto" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "task" ADD FOREIGN KEY ("id_escopo") REFERENCES "escopo" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "task" ADD FOREIGN KEY ("id_cliente") REFERENCES "cliente" ("id") DEFERRABLE INITIALLY IMMEDIATE;
