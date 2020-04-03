const client = require('./connectDB');
const queries = require('./DbQueries');

module.exports = async (HttpRequest, HttpResponse) => {
    queries.getThreadBySlugOrId.values = [
        HttpRequest.params.slug_or_id,
    ];

    let getThread = await client.query(queries.getThreadBySlugOrId);

    if (getThread.rows.length === 0) {
        sendThreadNotFound(HttpRequest, HttpResponse);
    } else {
        let getQuery = undefined;

        if (HttpRequest.query.sort === 'tree' ) {
            getQuery = Object.assign({}, queries.getThreadPostsTree);
            getQuery.values = [
                HttpRequest.params.slug_or_id,
            ];
            if (HttpRequest.query.since) {
                if (HttpRequest.query.desc === 'true') {
                    getQuery.text += ` WHERE r.path::int[] < (SELECT path FROM recursetree WHERE id = ${HttpRequest.query.since})::int[] `
                } else {
                    getQuery.text += ` WHERE r.path::int[] > (SELECT path FROM recursetree WHERE id = ${HttpRequest.query.since})::int[] `
                }
            }

            getQuery.text += ' ORDER BY r.path ';

            if (HttpRequest.query.desc === 'true') {
                getQuery.text += ' DESC '
            }

            if (HttpRequest.query.limit) {
                getQuery.text += ` LIMIT ${HttpRequest.query.limit} `;
            }

        }
        else if (HttpRequest.query.sort === 'parent_tree') {
            getQuery = Object.assign({}, queries.getThreadPostsTree);
            getQuery.values = [
                HttpRequest.params.slug_or_id,
            ];

            if (HttpRequest.query.since) {
                if (HttpRequest.query.desc === 'true') {
                    getQuery.text += ` WHERE r.path[1]::int < (SELECT path[1] FROM recursetree WHERE id = ${HttpRequest.query.since})::int `;
                } else {
                    getQuery.text += ` WHERE r.path[1]::int > (SELECT path[1] FROM recursetree WHERE id = ${HttpRequest.query.since})::int `;
                }
            }


            if (HttpRequest.query.desc === 'true') {
                if (!HttpRequest.query.since) {
                    getQuery.text += ' ORDER BY r.path[1] DESC, r.path ';
                } else {
                    getQuery.text += ' ORDER BY r.path[1] , r.path ';
                }
            } else {
                getQuery.text += ' ORDER BY r.path ';
            }


        }
        else {
            getQuery = Object.assign({}, queries.getThreadPostsFlat);
            getQuery.values = [
                HttpRequest.params.slug_or_id,
            ];

            if (HttpRequest.query.since) {
                if (HttpRequest.query.desc === 'true'){
                    getQuery.text += ` AND p.id < ${Number(HttpRequest.query.since)} `
                } else {
                    getQuery.text += ` AND p.id > ${Number(HttpRequest.query.since)} `
                }
            }

            if (HttpRequest.query.sort === 'flat') {
                getQuery.text += ' ORDER BY  p.created ';
                if (HttpRequest.query.desc === 'true') {
                    getQuery.text += ' DESC '
                }
                getQuery.text += ' ,p.id ';

            } else {
                getQuery.text += ' ORDER BY p.id';
            }

            if (HttpRequest.query.desc === 'true') {
                getQuery.text += ' DESC '
            }



            if (HttpRequest.query.limit) {
                getQuery.text += ` LIMIT ${HttpRequest.query.limit} `;
            }
        }


        console.log(getQuery);

        let posts = await client.query(getQuery);
        let response = [];

        if (HttpRequest.query.limit && HttpRequest.query.sort === 'parent_tree') {
            let prev = undefined;
            // if (HttpRequest.query.since && HttpRequest.query.limit && HttpRequest.query.desc === 'true') {
            //     HttpRequest.query.limit = 1;
            // }
            for (const data of posts.rows) {
                if (!prev)
                    prev = data[8];
                if (data[8] !== prev) {
                    HttpRequest.query.limit--;
                    prev = data[8];
                }
                if (HttpRequest.query.limit === 0){
                    break;
                }

                response.push({
                    author: data[0],
                    created: data[1],
                    forum: data[2],
                    id: data[3],
                    isEdited: Boolean(data[4]),
                    message: data[5],
                    parent: data[6],
                    thread: data[7],
                })

            }
        } else {
            for (const data of posts.rows) {
                response.push({
                    author: data[0],
                    created: data[1],
                    forum: data[2],
                    id: data[3],
                    isEdited: Boolean(data[4]),
                    message: data[5],
                    parent: data[6],
                    thread: data[7],
                })
            }
        }


        HttpResponse.status(200).json(response);

    }
};

let sendThreadNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with slug or id ${HttpRequest.params.slug_or_id}`
    })
};
