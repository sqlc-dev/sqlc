-- name: InsertOrders :exec
insert into Orders (id,name)
select id , CASE WHEN @name_do_update::BOOLEAN THEN @name ELSE s.name END 
from Orders s ;
