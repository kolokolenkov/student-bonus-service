CREATE TABLE students (
                          id BIGSERIAL PRIMARY KEY,

                          k_number BIGINT NOT NULL UNIQUE,
                          fio VARCHAR(255) NOT NULL,

                          period DATE,

                          student_contract_number VARCHAR(100),
                          student_contract_date DATE,

                          start_date_of_training DATE NOT NULL,
                          end_date_of_training DATE,

                          employment_contract VARCHAR(100),

                          location_code VARCHAR(50),
                          department_code VARCHAR(50),

                          wage_rate NUMERIC(12, 2),

                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);