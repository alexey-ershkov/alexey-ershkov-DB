const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = async (HttpRequest, HttpResponse) => {
    queries.getForumBySlug.values = [
        HttpRequest.params.slug,
    ];

    let forumResponse = await client.query(queries.getForumBySlug);
    if (forumResponse.rows.length === 0) {
        sendForumNotFound(HttpRequest, HttpResponse);
    } else {
        let getQuery = undefined;
        if (!HttpRequest.query.since) {
            getQuery = Object.assign({}, queries.getForumUsersBySlug);
        } else {
            if (HttpRequest.query.desc === 'true') {
                getQuery = Object.assign({}, queries.getForumUsersBySlugSinceDESC);
            } else {
                getQuery = Object.assign({}, queries.getForumUsersBySlugSinceASC)
            }
        }
        getQuery.values = [
            HttpRequest.params.slug,
        ];

        if (HttpRequest.query.since) {
            getQuery.values.push(HttpRequest.query.since);
        }

        if (HttpRequest.query.desc === 'true') {
            getQuery.text += ' DESC '
        }

        if (HttpRequest.query.limit) {
            if (!HttpRequest.query.since) {
                getQuery.text += ' LIMIT $2 ';
            } else {
                getQuery.text += ' LIMIT $3 '
            }
            getQuery.values.push(HttpRequest.query.limit);
        }

        let getForumUsers = await client.query(getQuery);
        let response = [];

        for (const elem of getForumUsers.rows) {
            queries.getUserByNickname.values = [
                elem[0],
            ];

            let getUserResponse = await client.query(queries.getUserByNickname);
            let userData = getUserResponse.rows[0];
            response.push({
                about: userData[3],
                email: userData[0],
                fullname: userData[1],
                nickname: userData[2],
            })
        }

        HttpResponse.status(200).json(response);

    }

};

let sendForumNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find forum with slug  ${HttpRequest.params.slug}`
    })
};
