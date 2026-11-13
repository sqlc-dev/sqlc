-- name: Methods :one
SELECT doc.value('(/a/@n)[1]', 'int') AS v, doc.value('(/a)[1]', 'nvarchar(20)') AS s,
    doc.exist('/a') AS e, geo.STDistance(geo) AS dist,
    geo.STAsText() AS txt, shape.ToString() AS shape_text, node.ToString() AS path,
    node.GetLevel() AS lvl, node.GetAncestor(@up) AS anc,
    hierarchyid::GetRoot() AS root, geography::Point(@lat, @long, 4326) AS pt
FROM m WHERE id = @id;

-- name: Nearby :many
SELECT id FROM m WHERE geo.STDistance(geography::Point(@lat, @long, 4326)) < @meters AND node.IsDescendantOf(@ancestor) = 1;
