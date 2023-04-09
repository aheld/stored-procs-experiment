\set banner_id '''f4bd6cdc-eb4b-4f74-8565-c243d3fdf20a'''

insert into list_items(banner_id, user_text) 
values (:banner_id, 'first item')

call list_items_insert(:banner_id, 1, 'testing proc insert');

call list_items_update(:banner_id, 1, 1, 'testing proc update');