SELECT event_object_table AS table_name ,trigger_name         
FROM information_schema.triggers
WHERE event_object_table not like 'pg_%'
GROUP BY table_name , trigger_name 
ORDER BY table_name ,trigger_name ;

SELECT
    routine_name
FROM 
    information_schema.routines
WHERE 
    routine_type = 'FUNCTION'
AND
    routine_schema = 'public';