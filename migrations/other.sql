INSERT INTO banners (content, is_active, tagIds, featureId)
VALUES ('"asdv"'::jsonb, true, '{4,5,6}', 3)
RETURNING id;


UPDATE banners b SET tagids = '{1,2}', featureid = 2 WHERE id = 6;


SELECT b.id,
       b.content,
       b.is_active,
       b.created,
       b.updated,
       b.featureId,
       b.tagIds
FROM feature_tag ft
         JOIN banners b on b.id = ft.bannerId
WHERE EXISTS (SELECT 1
              FROM feature_tag ft2
              WHERE ft2.bannerId = b.id
                AND ft2.tagId = 2)
GROUP BY b.id
ORDER BY b.id
LIMIT 3 OFFSET 0;