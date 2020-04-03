const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getForumBySlug.values = [
        HttpRequest.params.slug
    ];

    client.query(queries.getForumBySlug)
        .then(response => {
            if (response.rows.length !== 0) {
                sendInfo(HttpResponse, response.rows[0]);
            } else {
                sendNotFound(HttpRequest, HttpResponse);
            }
        })
        .catch( error =>{
            console.log(error);
            sendError(HttpResponse);
        })
};

let sendInfo = (HttpResponse, data) => {
    HttpResponse.status(200);
    HttpResponse.json({
        posts: Number(data[0]),
        slug: data[1],
        threads: Number(data[2]),
        title: data[3],
        user: data[4],
    });
};

let sendNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${HttpRequest.params.slug}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
