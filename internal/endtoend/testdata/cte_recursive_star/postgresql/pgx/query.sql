-- name: GetDictTree :many
with recursive dictTree(id, code, parent_code, label, value, path, depth) AS (
	select id, code, parent_code, label, value, ARRAY[COALESCE((select id from dict where code=''),'virtual_root'), id], 1 as depth from dict where app_id = '1' and parent_code = '' and is_delete=false
	union
		select d.id, d.code, d.parent_code, d.label, d.value, t.path || ARRAY[d.id], t.depth+1 as depth from dict d join dictTree t on d.parent_code = t.code and not d.id = ANY(t.path) and d.is_delete=false
)
select * from dictTree d order by depth, parent_code;