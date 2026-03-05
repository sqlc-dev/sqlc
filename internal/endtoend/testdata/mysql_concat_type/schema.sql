CREATE TABLE ticket
(
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
    ticket_status       TINYINT      NOT NULL,
    title               VARCHAR(255) NOT NULL,
    created_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
