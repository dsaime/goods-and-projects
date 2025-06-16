CREATE TABLE good_events
(
    id         SERIAL PRIMARY KEY,                                   -- id записи
    kind       TEXT      NOT NULL NOT NULL DEFAULT '',               -- тип события
    data       TEXT      NOT NULL NOT NULL DEFAULT '',               -- данные события
    ack        BOOLEAN   NOT NULL          DEFAULT FALSE,            -- подтверждение получения, со стороны читателя
    created_at TIMESTAMP NOT NULL          DEFAULT CURRENT_TIMESTAMP -- дата и время
);
/*
   2
   [server] (write)> [pgsql] +entity_update
                            \(write)> [nats] <(read) [event_listener] (write)> [clickhouse]
 */