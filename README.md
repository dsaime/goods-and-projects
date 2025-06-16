# goods-and-projects
Выполнение очередного тестового задания


[server] (write)> [pgsql] +entity_update
               \(write)> [nats] <(read) [event_listener] (write)> [clickhouse]

flowchart
```mermaid
---
config:
  theme: redux-dark
  layout: fixed
  look: classic
---
flowchart TD
    A(["server"]) -- write --> B@{ label: "<div style=\"color:\" data-darkreader-inline-color=\"\"><pre style=\"font-family:'JetBrains\"><span style=\"color:\" data-darkreader-inline-color=\"\">pgsql</span></pre></div>" }
    B -- write --> C@{ label: "<div style=\"color:\" data-darkreader-inline-color=\"\"><pre style=\"font-family:'JetBrains\"><span style=\"color:\" data-darkreader-inline-color=\"\">nats</span></pre></div>" }
    B -- save --> D@{ label: "<div style=\"color:\" data-darkreader-inline-color=\"\"><pre style=\"font-family:'JetBrains\">entity_update</pre></div>" }
    n1(["event_listener"]) -- read --> C
    n1 -- write --> n2(["clickhouse"])
    B@{ shape: stadium}
    C@{ shape: stadium}
    D@{ shape: rect}

```

sequenceDiagram
```mermaid
---
config:
  look: classic
  theme: redux-dark-color
---
sequenceDiagram
  participant Alice as server
  participant Bob as pgsql_tx
  participant P1 as nats
  participant P2 as event_listener
  participant P3 as clickhouse
  P2 ->> P1: subscribe
  Alice ->> Bob: modify_entity
  Bob ->> Bob: insert or update row
  Bob ->> P1: publish_event
  P1 ->> P2: new_msg
  P2 ->> P3: add_event

```