CREATE TABLE changes (
    ranked INT NOT NULL
);

CREATE TABLE changes_ranked (
    LIKE changes INCLUDING ALL,
    rank_by_effect_size INT NOT NULL,
    rank_by_abs_percent_change INT NOT NULL
);
