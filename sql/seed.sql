INSERT INTO drives (id, source_url, uri, amount)
VALUES ('3656cf1d-8826-404c-8f85-77f3e1f50464', 'https://www.reddit.com/r/wholesomememes', 'PrettyPinkMoon', 0);

INSERT INTO charities (id, name, logo_url, description, summary, jg_charity_id)
VALUES ('9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 'The Demo Charity', 'https://images.staging.justgiving.com/image/fd300863-43d6-4da7-b5ac-724e008f483d.png"', '29c50192-e194-4fd8-9ae5-333d54e9c357', '', 2050);

INSERT INTO donations (id, drive_id, charity_id, amount, currency_code, reference_code, status, created, next_check) VALUES
('b17f0a2d-8de6-4009-8efa-9aca898338c3', '3656cf1d-8826-404c-8f85-77f3e1f50464', '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 10.00, 'USD', 'ch-1234567890', 'Pending', NOW(), NOW()),
('3cd425a8-c376-4352-b5c6-8ebec5e1d8fe', '3656cf1d-8826-404c-8f85-77f3e1f50464', '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 13.00, 'USD', 'ch-1234567891', 'Pending', '2019-01-01', '2019-02-20'),
('899304f6-ea9d-4bc2-b7c1-df04aaa3b798', '3656cf1d-8826-404c-8f85-77f3e1f50464', '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 14.00, 'USD', 'ch-1234567892', 'Accepted', NOW(), '2019-01-01')
;
