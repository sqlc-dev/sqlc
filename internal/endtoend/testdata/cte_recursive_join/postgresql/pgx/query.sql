-- name: GetReferralChain :many
WITH RECURSIVE referral_chain AS (
    SELECT user_id,
           referrer_id,
           1 AS level
    FROM user_referrals r
    WHERE r.user_id = $1

    UNION ALL

    SELECT r.user_id,
           r.referrer_id,
           rc.level + 1
    FROM user_referrals r
             JOIN referral_chain rc ON r.user_id = rc.referrer_id
    WHERE rc.level < @level
)
SELECT *
FROM referral_chain
;