DROP TRIGGER IF EXISTS trg_update_finance_balance ON finance_transactions;
DROP FUNCTION IF EXISTS update_jurisdiction_balance();
DROP TRIGGER IF EXISTS trg_protect_finance_transactions ON finance_transactions;
DROP FUNCTION IF EXISTS protect_transactions();
DROP TABLE IF EXISTS finance_balances;
DROP TABLE IF EXISTS finance_transactions;
DROP TABLE IF EXISTS finance_categories;
DROP TYPE IF EXISTS transaction_type;
