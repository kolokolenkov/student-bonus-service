CREATE TABLE student_hours (
                               id BIGSERIAL PRIMARY KEY,

                               k_number BIGINT NOT NULL,

                               hours NUMERIC(10, 2) NOT NULL DEFAULT 0,
                               initial_training NUMERIC(10, 2) NOT NULL DEFAULT 0,

                               period_date DATE NOT NULL,

                               student_contract_number VARCHAR(100),

                               created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                               CONSTRAINT fk_student_hours_student
                                   FOREIGN KEY (k_number)
                                       REFERENCES students(k_number),

                               CONSTRAINT positive_hours
                                   CHECK (
                                       hours >= 0
                                           AND initial_training >= 0
                                       )
);