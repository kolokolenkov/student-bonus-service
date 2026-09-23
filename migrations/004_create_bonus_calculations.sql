CREATE TABLE bonus_calculations (
                                    id BIGSERIAL PRIMARY KEY,

                                    k_number BIGINT NOT NULL,

                                    period DATE NOT NULL,

                                    bonus_type VARCHAR(50) NOT NULL,

                                    hours NUMERIC(10, 2) NOT NULL,

                                    rate NUMERIC(12, 2) NOT NULL,

                                    amount NUMERIC(12, 2) NOT NULL,

                                    status VARCHAR(50) NOT NULL DEFAULT 'CALCULATED',

                                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

                                    CONSTRAINT fk_bonus_calculations_student
                                        FOREIGN KEY (k_number)
                                            REFERENCES students(k_number),

                                    CONSTRAINT positive_calculation_hours
                                        CHECK (hours >= 0),

                                    CONSTRAINT positive_calculation_rate
                                        CHECK (rate >= 0),

                                    CONSTRAINT positive_calculation_amount
                                        CHECK (amount >= 0),

                                    CONSTRAINT valid_bonus_status
                                        CHECK (
                                            status IN (
                                                       'CALCULATED',
                                                       'SENT',
                                                       'SEND_ERROR'
                                                )
                                            ),

                                    CONSTRAINT unique_student_bonus_period
                                        UNIQUE (
                                                k_number,
                                                period,
                                                bonus_type
                                            )
);