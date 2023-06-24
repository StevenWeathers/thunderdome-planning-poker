-- Set Storyboard Owner --
CREATE PROCEDURE set_storyboard_owner(storyboardId UUID, ownerId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard SET updated_date = NOW(), owner_id = ownerId WHERE id = storyboardId;
END;
$$;