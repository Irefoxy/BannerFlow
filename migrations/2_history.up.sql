CREATE TABLE IF NOT EXISTS banner_history
(
    id        SERIAL PRIMARY KEY,
    bannerId  INT REFERENCES banners (id) ON DELETE CASCADE,
    content   JSONB NOT NULL,
    version   INT,
    tagIds    INT[],
    featureId INT
);

CREATE OR REPLACE PROCEDURE save_banner_to_history (records banners)
LANGUAGE plpgsql
AS $$
DECLARE
    vs INT;
BEGIN
    vs := (SELECT COALESCE(max(version), 0) FROM banner_history WHERE bannerId = records.id) + 1;
    INSERT INTO banner_history (bannerId, content, version, tagIds, featureId)
    VALUES (records.id, records.content, vs, records.tagids, records.featureid);
    DELETE FROM banner_history bh WHERE bh.bannerId = records.id AND bh.version
    NOT IN (SELECT bh2.version FROM banner_history bh2 WHERE bh2.bannerId = records.id
            ORDER BY bh2.version DESC LIMIT 3);
END;
$$;

CREATE OR REPLACE FUNCTION update_banner_func()
    RETURNS TRIGGER AS
$$
BEGIN
    IF old.tagIds IS DISTINCT FROM new.tagIds THEN
        DELETE FROM feature_tag WHERE bannerid = old.id;
        CALL insert_into_feature_tag(NEW);
    END IF;
    CALL save_banner_to_history(OLD);
    NEW.updated = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_banner
    BEFORE UPDATE ON banners
    FOR EACH ROW
EXECUTE FUNCTION update_banner_func();
