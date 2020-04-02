const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getUserByNickname.values = [
        HttpRequest.body.author
    ];

    queries.getForumBySlugSimple.values = [
      HttpRequest.params.forum
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            if (response.rows.length !== 0) {

                client.query(queries.getForumBySlugSimple)
                    .then( ans => {
                            if (ans.rows.length !== 0) {
                                createThread(HttpRequest, HttpResponse);
                            } else {
                                sendForumNotFound(HttpRequest, HttpResponse);
                            }
                        }
                    )
                    .catch( e =>{
                        console.log(e);
                        sendError(HttpResponse);
                    } )
            } else {
                sendUserNotFound(HttpRequest, HttpResponse);
            }
        })
        .catch(()=>{
            sendError(HttpResponse);
        })
};

let createThread = (HttpRequest, HttpResponse) => {
    queries.getThread.values = [
        HttpRequest.body.author,
        HttpRequest.params.forum,
        HttpRequest.body.message,
        HttpRequest.body.title,
    ];

    client.query(queries.getThread)
        .then(response => {
            if (response.rows.length !== 0) {
                HttpResponse.status(409);
                sendThreadInfo(HttpResponse, response.rows[0]);
            } else {
                queries.createThread.values = [
                    HttpRequest.body.author,
                    HttpRequest.body.created,
                    HttpRequest.params.forum,
                    HttpRequest.body.message,
                    HttpRequest.body.title,
                    HttpRequest.body.slug,
                ];
                client.query(queries.createThread)
                    .then( () => {
                        client.query(queries.getThread)
                            .then( resp => {
                                HttpResponse.status(201);
                                sendThreadInfo(HttpResponse,resp.rows[0]);
                            })
                            .catch( e => {
                                console.log(e);
                                sendError(HttpResponse);
                            })
                    })
                    .catch( () => {

                        queries.getThreadBySlug.values = [
                            HttpRequest.body.slug,
                        ];

                        client.query(queries.getThreadBySlug)
                            .then(resp => {
                                HttpResponse.status(409);
                                sendThreadInfo(HttpResponse, resp.rows[0]);
                            });
                    });
            }
        })
        .catch(e => {
            console.log(e);
            sendError(HttpResponse);
        })

};

let sendThreadInfo = (HttpResponse, data) => {
    HttpResponse.json({
        author: data[5],
        created: data[3],
        id: data[0],
        forum: data[6],
        message: data[2],
        slug: data[4],
        title: data[1],
    })
};

let sendForumNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find forum with slug ${HttpRequest.params.forum}`
    })
};

let sendUserNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${HttpRequest.body.author}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
