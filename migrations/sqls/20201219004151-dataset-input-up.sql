CREATE TABLE IF NOT EXISTS dataset_db.public.copy_data (
  cpf varchar(255) null,
  private boolean null,
  incomplete boolean null,
  last_purchase date null,
  avg_ticket float(53) null,
  last_ticket float(53) null,
  frequent_store varchar(255) null,
  last_store varchar(255) null,
  cpf_invalid boolean not null default false,
  cnpj_invalid boolean not null default false
);