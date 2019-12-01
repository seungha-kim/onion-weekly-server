create table users (
    id uuid primary key default uuid_generate_v4(),
    created_at ts_default
);

create table user_profiles (
    user_id uuid not null references users (id),
    full_name text not null
);

create table email_credentials (
    user_id uuid not null references users (id),
    email citext not null unique,
    hashed_password text not null,
    created_at ts_default
);
