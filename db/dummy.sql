CREATE TABLE dummy
(
    id          BIGINT         NOT NULL PRIMARY KEY,
    product     VARCHAR(255)   NOT NULL,
    description VARCHAR(255),
    price       DECIMAL(10, 2) NOT NULL,
    qty         BIGINT         NOT NULL,
    date        VARCHAR(255)   NOT NULL
)