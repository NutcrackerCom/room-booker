insert into users (id, email, role)
values
    ('00000000-0000-0000-0000-000000000001', 'admin@example.com', 'admin'),
    ('00000000-0000-0000-0000-000000000002', 'user@example.com', 'user')
on conflict (id) do nothing;
