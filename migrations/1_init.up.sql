CREATE TABLE IF NOT EXISTS banners
(
    id        SERIAL PRIMARY KEY,
    content   JSONB NOT NULL,
    tagIds     INT[],
    featureId INT,
    created   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated   TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS deactivated (
    bannerId  INT PRIMARY KEY REFERENCES banners (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS feature_tag
(
    bannerId  INT REFERENCES banners (id) ON DELETE CASCADE,
    tagId     INT,
    featureId INT,
    PRIMARY KEY (tagId, featureId)
);

CREATE OR REPLACE PROCEDURE insert_into_feature_tag(records banners)
LANGUAGE plpgsql
AS $$
DECLARE
    tag INT;
BEGIN
    FOREACH tag IN ARRAY records.tagIds
        LOOP
            INSERT INTO feature_tag (bannerId, tagId, featureId)
            VALUES (records.id, tag, records.featureId);
        END LOOP;
END;
$$;

CREATE OR REPLACE FUNCTION insert_trigger_func()
    RETURNS TRIGGER AS $$
BEGIN
    CALL insert_into_feature_tag(NEW);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_banner
    AFTER INSERT ON banners
    FOR EACH ROW EXECUTE FUNCTION insert_trigger_func();