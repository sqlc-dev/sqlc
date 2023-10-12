-- name: Percentile :exec
UPDATE group_calc_totals gct 
SET npn = nem.npn 
FROM producer_group_attribute ga 
JOIN npn_external_map nem ON ga.npn_external_map_id = nem.id 
WHERE gct.group_id = ga.group_id;
