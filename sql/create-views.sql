DROP VIEW IF EXISTS drives_view;
CREATE VIEW drives_view AS
WITH sums AS (
    SELECT drive_id,
           SUM(usd_amount) AS usd_amount_total,
           SUM(donor_amount) AS donor_amount_total,
           COUNT(*) as num_donations,
           most_recent_usd_amount,
           most_recent_donor_amount,
           most_recent_time
    FROM (
             SELECT *,
                    LAST_VALUE(donor_amount)
                    OVER (PARTITION BY drive_id) AS most_recent_donor_amount,
                    LAST_VALUE(usd_amount)
                    OVER (PARTITION BY drive_id) AS most_recent_usd_amount,
                    LAST_VALUE(created)
                    OVER (PARTITION BY drive_id) AS most_recent_time
             FROM donations
             WHERE status = 'Accepted'
             ORDER BY created ASC
         ) T
    GROUP BY drive_id, most_recent_donor_amount, most_recent_usd_amount, most_recent_time
)
SELECT drives.*,
       COALESCE(sums.usd_amount_total, 0)       AS usd_amount_total,
       COALESCE(sums.donor_amount_total, 0)       AS donor_amount_total,
       COALESCE(sums.most_recent_donor_amount, 0) AS most_recent_donor_amount,
       COALESCE(sums.most_recent_usd_amount, 0) AS most_recent_usd_amount,
       COALESCE(sums.num_donations, 0) AS num_donations,
       sums.most_recent_time
FROM drives
         LEFT JOIN sums ON sums.drive_id = drives.id;

DROP VIEW IF EXISTS charities_view;
CREATE VIEW charities_view AS
WITH sums AS (
    SELECT charity_id,
           SUM(usd_amount) AS usd_amount_total,
           SUM(donor_amount) AS donor_amount_total,
           most_recent_usd_amount,
           most_recent_donor_amount,
           most_recent_time
    FROM (
             SELECT *,
                    LAST_VALUE(donor_amount)
                    OVER (PARTITION BY charity_id) AS most_recent_donor_amount,
                    LAST_VALUE(usd_amount)
                    OVER (PARTITION BY charity_id) AS most_recent_usd_amount,
                    LAST_VALUE(created)
                    OVER (PARTITION BY charity_id) AS most_recent_time
             FROM donations
             WHERE status = 'Accepted'
             ORDER BY created ASC
         ) T
    GROUP BY charity_id, most_recent_donor_amount, most_recent_usd_amount, most_recent_time
)
SELECT charities.*,
       COALESCE(sums.usd_amount_total, 0)       AS usd_amount_total,
       COALESCE(sums.donor_amount_total, 0)       AS donor_amount_total,
       COALESCE(sums.most_recent_donor_amount, 0) AS most_recent_donor_amount,
       COALESCE(sums.most_recent_usd_amount, 0) AS most_recent_usd_amount,
       sums.most_recent_time
FROM charities
         LEFT JOIN sums ON sums.charity_id = charities.id;

DROP VIEW IF EXISTS featured_charities_view;
CREATE VIEW featured_charities_view AS
SELECT * FROM charities_view WHERE feature_score > 0;