CREATE TABLE IF NOT EXISTS list_items(
   id serial PRIMARY KEY,
   user_id int NOT NULL,
   banner_id UUID, 
   user_text TEXT,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS user_id_idx ON list_items (id, user_id, banner_id);


CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE TRIGGER update_list_items_updated_at
    BEFORE UPDATE
    ON
        list_items
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE OR REPLACE FUNCTION list_items_insert("tenantId" UUID, "userId" integer, "item" character varying) returns int
LANGUAGE SQL
AS $$
    insert into list_items(banner_id, user_id, user_text) values ("tenantId", "userId", "item")
    RETURNING id;
$$;

CREATE OR REPLACE PROCEDURE list_items_update("tenantId" UUID, "userId" integer, "listItemId" integer, "item" character varying)
LANGUAGE plpgsql
AS $$
BEGIN
    update list_items set user_text="item" where banner_id="tenantId" and user_id="userId" and id="listItemId";
    IF NOT FOUND THEN -- UPDATE didn't touch anything
        RAISE EXCEPTION 'list item % not found', "listItemId";
    END IF;
END
$$;
