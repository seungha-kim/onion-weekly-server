create table workspaces
(
    id         uuid primary key default uuid_generate_v4(),
    name       text not null,
    created_by uuid not null references users (id),
    created_at ts_default
);

create table workspace_members
(
    user_id      uuid not null references users (id),
    workspace_id uuid not null references workspaces (id),
    created_at   ts_default
)