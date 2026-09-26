-- name: RepeatedName :many
SELECT id FROM events
WHERE name = {name:String} OR tag = {name:String} OR id = {id:UInt64};

-- name: InsertEvent :exec
INSERT INTO events (id, name, tag, amount, uid, ip, big)
VALUES ({id:UInt64}, {name:String}, {tag:Nullable(String)}, {amount:Decimal(10, 2)}, {uid:UUID}, {ip:IPv4}, {big:Int128});

-- name: PositionalEvent :exec
INSERT INTO events (id, name, tag, amount, uid, ip, big)
VALUES (?, ?, ?, ?, ?, ?, ?);
