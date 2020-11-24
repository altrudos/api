INSERT INTO drives (id, source_url, source_type, source_key, source_meta, uri, amount, created_at)
VALUES ('3656cf1d-8826-404c-8f85-77f3e1f50464', 'https://www.reddit.com/r/PublicFreakout/comments/fnnujk/ny_not_handling_this_shit_well/flbfd1j/', 'reddit_comment', 'flbfd1j', '{"subreddit": "PublicFreakout", "author": "BRuX-", "title": "NY not handling this shit well"}', 'PrettyPinkMoon', 0, NOW()- INTERVAL '10 DAYS'),
('edfa5906-303f-4a5c-941c-0006d604cf28', 'https://www.reddit.com/r/PublicFreakout/comments/fnnujk/ny_not_handling_this_shit_well/', 'reddit_post', 'fnnujk', '{"subreddit": "PublicFreakout", "author": "BRuX-", "title": "NY not handling this shit well"}', 'PrettyRedMoon', 0, NOW()- INTERVAL '3 DAYS'),
('1181e5c7-c82c-424e-b1ef-b13f61d616dc', 'https://gist.github.com/ThaumRystra/ffb264dea8c32e15de95f775596194a4', 'url', 'https://gist.github.com/ThaumRystra/ffb264dea8c32e15de95f775596194a4', '{}', 'LargeBearFire', 0, NOW()- INTERVAL '1 DAY');

INSERT INTO charities (id, name, logo_url, description, summary, jg_charity_id, feature_score)
VALUES ('9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 'The Demo Charity',
        'https://images.staging.justgiving.com/image/fd300863-43d6-4da7-b5ac-724e008f483d.png"',
        'The Demo charity takes demos and makes them charitable.', '', 2050, 1),
       ('627e0410-c75d-48c8-b41f-6318d04f1e65', 'Fake Charity',
        'https://images.staging.justgiving.com/image/fd300863-43d6-4da7-b5ac-724e008f483d.png"',
        'This charity is a fake one used for testing. It provides free water to other fake test data.', '', 2051, 2);

INSERT INTO donations (id, drive_id, charity_id, donor_amount, donor_currency, final_amount, final_currency, usd_amount, reference_code,
                       donor_name, status, created_at, next_check)
VALUES ('b17f0a2d-8de6-4009-8efa-9aca898338c3', '3656cf1d-8826-404c-8f85-77f3e1f50464',
        '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 1000, 'USD', 98, 'USD', 98, 'ch-1586573040421095000', '', 'Pending', NOW(), NOW()),
       ('3cd425a8-c376-4352-b5c6-8ebec5e1d8fe', '3656cf1d-8826-404c-8f85-77f3e1f50464',
        '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 1300, 'GBP', 1243, 'USD', 1243, 'ch-1234567891', '', 'Pending', '2019-01-01',
        '2019-02-20'),
       ('899304f6-ea9d-4bc2-b7c1-df04aaa3b798', '3656cf1d-8826-404c-8f85-77f3e1f50464',
        '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 800, 'USD', 780, 'USD', 780, 'ch-1234567892', '', 'Accepted',
        NOW() - INTERVAL '8 DAYS', '2019-01-01'),
       ('71800a64-fb75-4ced-bb0e-76eaaa4cf440', 'edfa5906-303f-4a5c-941c-0006d604cf28',
        '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 32000, 'USD', 31001, 'USD', 31001, 'ch-1234567893', 'Big Spender', 'Accepted',
        NOW() - INTERVAL '1 MIN', '2019-01-01'),
       ('4f6241f4-903f-459b-bf02-fd14c825d8d8', '3656cf1d-8826-404c-8f85-77f3e1f50464',
        '9d0b23cd-657b-4cc4-8258-a8cabb1f6847', 1400, 'CAD', 1332, 'USD', 1332, 'ch-1234567894', 'FordonGreeman', 'Accepted', NOW(),
        '2019-01-01'),
    ('66b0765d-fac6-4e8b-b1ee-b381ef0df228', '1181e5c7-c82c-424e-b1ef-b13f61d616dc',
    '627e0410-c75d-48c8-b41f-6318d04f1e65', 4325, 'CAD', 4101, 'USD', 4101, 'ch-5984302543', 'Alex Blake Joel Miles Miles Nick Wyatt', 'Accepted', NOW(), NOW());
