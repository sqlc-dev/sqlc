CREATE TABLE wishlist_item (
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    date_from  DATE         NOT NULL,
    updated_at TIMESTAMP    NULL DEFAULT NULL
);
