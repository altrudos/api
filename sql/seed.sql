INSERT INTO drives (source_url, uri, amount)
VALUES ('https://www.reddit.com/r/wholesomememes', 'PrettyPinkMoon', 0);

INSERT INTO charities (name, logo_url, description, summary, jg_charity_id)
VALUES ('The Demo Charity', 'https://images.staging.justgiving.com/image/fd300863-43d6-4da7-b5ac-724e008f483d.png"', '29c50192-e194-4fd8-9ae5-333d54e9c357', '', 2050);

INSERT INTO donations (drive_id, charity_id, amount, currency_code, reference_code, status, created, next_check) VALUES
(1, 1, 10.00, 'USD', 'ch-1234567890', 'pending', NOW(), NOW()),
(1, 1, 13.00, 'USD', 'ch-1234567891', 'pending', '2019-01-01', '2019-02-20'),
(1, 1, 14.00, 'USD', 'ch-1234567892', 'approved', NOW(), '2019-01-01')
;