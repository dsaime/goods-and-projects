# goods-and-projects
Выполнение очередного тестового задания

sequenceDiagram
```mermaid
---
config:
  look: classic
  theme: redux-dark-color
---
sequenceDiagram
    participant server as server
    participant pgsql_tx as pgsql_tx
    participant P1 as nats
    participant P2 as event_listener
    participant P3 as clickhouse
    P2 ->> P1: subscribe
    server ->> pgsql_tx: modify_entity
    pgsql_tx ->> pgsql_tx: insert or update row
    pgsql_tx ->> P1: publish_event
    P1 -->> pgsql_tx: ok
    pgsql_tx -->> server: ok
    P1 -->> P2: new_msg
    P2 ->> P3: add_event
    P3 -->> P2: ok

```