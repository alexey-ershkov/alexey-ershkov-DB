const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getThreadBySlugOrIdWithVotes.values = [
       HttpRequest.params.slug,
    ];

    client.query(queries.getThreadBySlugOrIdWithVotes)
        .then( response => {
            if (response.rows.length !== 0) {
                sendThreadInfo(HttpResponse, response.rows[0]);
            } else {
                sendThreadNotFound(HttpRequest, HttpResponse);
            }
        })
        .catch( e => {
            console.log(e);
            sendError(HttpResponse);
        })
};

let sendThreadInfo = (HttpResponse, data) => {
    HttpResponse.status(200);
    HttpResponse.json({
        author: data[5],
        created: data[3],
        id: data[0],
        forum: data[6],
        message: data[2],
        slug: data[4],
        title: data[1],
        votes: data[7]
    })
};

let sendThreadNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with slug or id ${HttpRequest.params.slug}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
