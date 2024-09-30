CREATE TABLE dummy
(
    id          BIGINT        NOT NULL PRIMARY KEY,
    product     CHAR(100)     NOT NULL,
    description VARCHAR(255),
    price       DOUBLE(10, 2) NOT NULL,
    qty         BIGINT        NOT NULL,
    date        DATETIME      NOT NULL
)