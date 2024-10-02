package queries

var GetContentTypes = `
	SELECT
	id,
	uuid,
	name,
	title,
	title_ru,
	description,
	parent_id,
	parent_name FROM tbl_content_type
`

var GetContentTypesWithContent = `
SELECT
    ct.id,
    ct.uuid,
    ct.name,
    ct.title,
    ct.title_ru,
    ct.description,
    ct.parent_id,
    ct.parent_name,
    json_agg(
        json_build_object(
            'id', c.id,
            'uuid', c.uuid,
            'lang_id', c.lang_id,
            'content_type_id', c.content_type_id,
            'title', c.title,
            'slogan', c.slogan,
            'subtitle', c.subtitle,
            'description', c.description,
            'image_url', c.image_url,
            'video_url', c.video_url,
            'step', c.step,
            'created_at', c.created_at::VARCHAR,
            'updated_at', c.updated_at::VARCHAR,
            'active', c.active,
            'deleted', c.deleted
        )
    ) AS contents
FROM
    tbl_content_type ct
LEFT JOIN
    tbl_content c ON c.content_type_id = ct.id AND c.deleted = 0
GROUP BY
    ct.id
ORDER BY
    ct.id;
`

//
//var GetNewOrders = `
//    SELECT
//        o.id,
//        o.order_number,
//        ARRAY(
//            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
//        ) AS file_paths,
//        o.users_id AS user_id,
//        o.workers_id AS worker_id,
//        o.address,
//        o.date::VARCHAR,
//        o.time::VARCHAR,
//        o.time_duration,
//        (
//            SELECT JSON_BUILD_OBJECT(
//                'id', s.id,
//                'tk', GET_STATUS_TITLE(s.id, 1),
//                'ru', GET_STATUS_TITLE(s.id, 2),
//                'en', GET_STATUS_TITLE(s.id, 3)
//            )
//            FROM statuses s
//            WHERE s.id = o.statuses_id
//        ) AS status,
//        o.description,
//        ARRAY(
//            SELECT
//                JSON_BUILD_OBJECT(
//                    'id', os.services_id,
//                    'image', (
//                        SELECT s.image FROM services s
//                        WHERE s.id = os.services_id
//                    ),
//                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
//                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
//                    'en', GET_SERVICE_TITLE(os.services_id, 3)
//                )
//            FROM orders_services os
//            WHERE os.orders_id = o.id
//        ) AS services,
//        o.read_by_admin
//    FROM orders o
//    WHERE o.read_by_admin = false
//    ORDER BY o.id DESC
//`
