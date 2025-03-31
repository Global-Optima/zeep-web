-- Re-add the unique constraint on employee_id
ALTER TABLE employee_tokens ADD CONSTRAINT employee_tokens_employee_id_key UNIQUE (employee_id);