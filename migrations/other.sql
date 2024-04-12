INSERT INTO banners (content, tagIds, featureId)
VALUES ('"asdv"'::jsonb, '{1,2,3}', 3)
RETURNING id;


UPDATE banners b SET tagids = '{33}', featureid = 1 WHERE id = 3;

SELECT ARRAY_AGG(DISTINCT bannerid) ids FROM feature_tag WHERE tagid = 1
GROUP BY tagid;

CALL choose_banner_from_history(1,3);

DELETE FROM banners WHERE id = ANY ('{3}');