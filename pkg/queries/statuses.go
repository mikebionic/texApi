package queries

var GetStatuses = `
    SELECT
        s.id,
        GET_STATUS_TITLE(s.id, 1) AS tk,
        GET_STATUS_TITLE(s.id, 2) AS ru,
        GET_STATUS_TITLE(s.id, 3) AS en
    FROM statuses s
`
