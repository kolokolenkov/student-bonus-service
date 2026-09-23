CREATE TABLE bonus_rates (
                             id BIGSERIAL PRIMARY KEY,

                             bonus_type VARCHAR(50) NOT NULL,

                             rate NUMERIC(12, 2) NOT NULL,

                             valid_from DATE NOT NULL,
                             valid_to DATE,

                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                             CONSTRAINT positive_bonus_rate
                                 CHECK (rate >= 0),

                             CONSTRAINT valid_rate_period
                                 CHECK (
                                     valid_to IS NULL
                                         OR valid_to >= valid_from
                                     )
);