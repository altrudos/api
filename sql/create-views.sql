DROP VIEW IF EXISTS drives_view;
CREATE VIEW drives_view AS
WITH sums AS (
    SELECT drive_id,
           SUM(final_amount) AS final_amount_total,
           SUM(donor_amount) AS donor_amount_total,
           MAX(final_amount) AS final_amount_max,
           MAX(donor_amount) AS donor_amount_max,
           most_recent_final_amount,
           most_recent_donor_amount,
           most_recent_time
    FROM (
             SELECT *,
                    LAST_VALUE(donor_amount)
                    OVER (PARTITION BY drive_id) AS most_recent_donor_amount,
                    LAST_VALUE(final_amount)
                    OVER (PARTITION BY drive_id) AS most_recent_final_amount,
                    LAST_VALUE(created)
                    OVER (PARTITION BY drive_id) AS most_recent_time
             FROM donations
             WHERE status = 'Accepted'
             ORDER BY created ASC
         ) T
    GROUP BY drive_id, most_recent_donor_amount, most_recent_final_amount, most_recent_time
)
SELECT drives.*,
       COALESCE(sums.final_amount_total, 0)       AS final_amount_total,
       COALESCE(sums.final_amount_max, 0)         AS final_amount_max,
       COALESCE(sums.donor_amount_total, 0)       AS donor_amount_total,
       COALESCE(sums.donor_amount_max, 0)         AS donor_amount_max,
       COALESCE(sums.most_recent_donor_amount, 0) AS most_recent_donor_amount,
       COALESCE(sums.most_recent_final_amount, 0) AS most_recent_final_amount,
       sums.most_recent_time
FROM drives
         LEFT JOIN sums ON sums.drive_id = drives.id;

DROP VIEW IF EXISTS charities_view;
CREATE VIEW charities_view AS
WITH sums AS (
    SELECT charity_id,
           SUM(final_amount) AS final_amount_total,
           SUM(donor_amount) AS donor_amount_total,
           MAX(final_amount) AS final_amount_max,
           MAX(donor_amount) AS donor_amount_max,
           most_recent_final_amount,
           most_recent_donor_amount,
           most_recent_time
    FROM (
             SELECT *,
                    LAST_VALUE(donor_amount)
                    OVER (PARTITION BY charity_id) AS most_recent_donor_amount,
                    LAST_VALUE(final_amount)
                    OVER (PARTITION BY charity_id) AS most_recent_final_amount,
                    LAST_VALUE(created)
                    OVER (PARTITION BY charity_id) AS most_recent_time
             FROM donations
             WHERE status = 'Accepted'
             ORDER BY created ASC
         ) T
    GROUP BY charity_id, most_recent_donor_amount, most_recent_final_amount, most_recent_time
)
SELECT charities.*,
       COALESCE(sums.final_amount_total, 0)       AS final_amount_total,
       COALESCE(sums.final_amount_max, 0)         AS final_amount_max,
       COALESCE(sums.donor_amount_total, 0)       AS donor_amount_total,
       COALESCE(sums.donor_amount_max, 0)         AS donor_amount_max,
       COALESCE(sums.most_recent_donor_amount, 0) AS most_recent_donor_amount,
       COALESCE(sums.most_recent_final_amount, 0) AS most_recent_final_amount,
       sums.most_recent_time
FROM charities
         LEFT JOIN sums ON sums.charity_id = charities.id;
