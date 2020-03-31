const client = require('./connectDB');
const queries = require('./DbQueries');

module.exports = (HttpRequest, HttpResponse) => {


    queries.getUserByNickname.values = [
        HttpRequest.body.user
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            if (response.rows.length !== 0) {
                createForum(HttpRequest, HttpResponse);
            } else {
                sendNotFound(HttpRequest, HttpResponse);
            }
        })
        .catch(()=>{
            sendError(HttpResponse);
        })
};

let createForum = (HttpRequest, HttpResponse) => {
    queries.createForum.values = [
        HttpRequest.body.slug,
        HttpRequest.body.title,
        HttpRequest.body.user
    ];

    client.query(queries.createForum)
        .then(() => {
            HttpResponse.status(201);
        })
        .catch(() => {
            HttpResponse.status(409);
        });

    sendResponse(HttpRequest,HttpResponse);
};

let sendResponse = (HttpRequest, HttpResponse) => {
    queries.getForumBySlugSimple.values = [
        HttpRequest.body.slug
    ];

    client.query(queries.getForumBySlugSimple)
        .then(response => {
            const info = response.rows[0];
            HttpResponse.json({
                slug: info[0],
                title: info[1],
                user: info[2],
            })
        })
};

let sendNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${HttpRequest.body.user}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
