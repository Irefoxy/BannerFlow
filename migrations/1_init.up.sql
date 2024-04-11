CREATE TABLE IF NOT EXISTS banners
(
    id        SERIAL PRIMARY KEY,
    content   JSONB NOT NULL ,
    is_active BOOL NOT NULL,
    created   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated   TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS feature_tag
(
    bannerId INT REFERENCES banners(id),
    tagId   INT,
    featureId INT,
    PRIMARY KEY (tagId, featureId)
);