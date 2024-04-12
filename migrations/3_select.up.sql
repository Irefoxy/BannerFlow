CREATE OR REPLACE PROCEDURE choose_banner_from_history (bid INT, vn INT)
LANGUAGE plpgsql
AS $$
DECLARE
    history banner_history;
BEGIN
    SELECT * INTO history FROM banner_history WHERE bannerid = bid AND version = vn;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'No history for ID % and version %.', bid, vn;
    END IF;
    DELETE FROM banner_history WHERE bannerid = bid AND version = vn;
    UPDATE banners SET featureid = history.featureid, tagids = history.tagids, updated = current_timestamp, content = history.content
    WHERE id = bid;
END;
$$;