CREATE TABLE IF NOT EXISTS dataset_db.public.txtdata (
  cpf varchar(255) null,
  private boolean null,
  incomplete boolean null,
  last_purchase date null,
  avg_ticket numeric(5,2) null,
  last_ticket numeric(5,2) null,
  frequent_store varchar(255) null,
  last_store varchar(255) null
);