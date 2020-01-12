create table recurrings
(
    id           uuid primary key default uuid_generate_v4(),
    workspace_id uuid not null references workspaces (id),
    title        text not null,
    interval     int,
    created_at   ts_default
);

create table recurring_records
(
    id                uuid primary key default uuid_generate_v4(),
    description       text not null,
    actor_id          uuid not null references users (id),
    interval_snapshot int,
    created_at        ts_default,
    recurring_id      uuid not null references recurrings (id)
);