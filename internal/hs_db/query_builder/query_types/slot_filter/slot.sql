SELECT
    s.id,
    s.uploader,
    s.name,
    s.description,
    s.location_x,
    s.location_y,
    s.background,
    s.root_level,
    s.icon,
    s.initially_locked,
    s.sub_level,
    s.lbp1only,
    s.shareable,
    s.min_players,
    s.max_players,
    s.game,

    COUNT(DISTINCT h.owner) AS heart_count,
    COUNT(DISTINCT tu.owner) AS thumbs_up_count,
    COUNT(DISTINCT td.owner) AS thumbs_down_count,
    COUNT(DISTINCT p.owner) AS play_count
FROM
    slots AS s
        LEFT JOIN
    hearts AS h ON s.id = h.slot_id
        LEFT JOIN
    thumbs_up AS tu ON s.id = tu.slot_id
        LEFT JOIN
    thumbs_down AS td ON s.id = td.slot_id
        LEFT JOIN
    plays AS p ON s.id = p.slot_id
WHERE %s
GROUP BY
    s.id