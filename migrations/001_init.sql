create extension if not exists "pgcrypto";

create table if not exists users (
    id uuid primary key,
    email text unique not null,
    password_hash text,
    role text not null check (role in ('admin', 'user')),
    created_at timestamptz not null default now()
);

create table if not exists rooms (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    description text,
    capacity integer,
    created_at timestamptz not null default now()
);

create table if not exists schedules (
    id uuid primary key default gen_random_uuid(),
    room_id uuid not null unique references rooms(id) on delete cascade,
    days_of_week integer[] not null,
    start_time time not null,
    end_time time not null,
    created_at timestamptz not null default now(),
    check (array_length(days_of_week, 1) >= 1),
    check (start_time < end_time)
);

create table if not exists slots (
    id uuid primary key default gen_random_uuid(),
    room_id uuid not null references rooms(id) on delete cascade,
    start_at timestamptz not null,
    end_at timestamptz not null,
    created_at timestamptz not null default now(),
    check (end_at > start_at),
    unique (room_id, start_at)
);

create index if not exists idx_slots_room_start on slots(room_id, start_at);

create table if not exists bookings (
    id uuid primary key default gen_random_uuid(),
    slot_id uuid not null references slots(id) on delete cascade,
    user_id uuid not null references users(id),
    status text not null check (status in ('active', 'cancelled')),
    conference_link text,
    created_at timestamptz not null default now()
);

create unique index if not exists ux_bookings_active_slot
on bookings(slot_id)
where status = 'active';
