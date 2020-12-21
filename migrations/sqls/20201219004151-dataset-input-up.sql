CREATE TABLE IF NOT EXISTS dataset_db.public.original_data (
  cpf varchar(255) null,
  private varchar(255) null,
  incomplete varchar(255) null,
  last_purchase varchar(255) null,
  avg_ticket varchar(255) null,
  last_ticket varchar(255) null,
  frequent_store varchar(255) null,
  last_store varchar(255) null
);

CREATE TABLE IF NOT EXISTS dataset_db.public.copy_data (
  cpf varchar(255) null,
  private boolean null,
  incomplete boolean null,
  last_purchase date null,
  avg_ticket numeric(5,2) null,
  last_ticket numeric(5,2) null,
  frequent_store varchar(255) null,
  last_store varchar(255) null
);