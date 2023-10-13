-- name: UpdateAttribute :one
with updated_attribute as (UPDATE attribute_value
    SET
        val = CASE WHEN @filter_value::bool THEN @value ELSE val END
    WHERE attribute_value.id = @id
    RETURNING id,attribute,val)
select updated_attribute.id, val, name
from updated_attribute
         left join attribute on updated_attribute.attribute = attribute.id;
