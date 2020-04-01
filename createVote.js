const client = require('./connectDB');
const queries = require('./DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getUserByNickname.values = [
        HttpRequest.body.nickname
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            if (response.rows.length === 0) {
                sendUserNotFound(HttpRequest, HttpResponse);
            } else {
                userFound(HttpRequest, HttpResponse);
            }
        })
        .catch( e => {
            console.log(e);
            sendError(HttpResponse);
        })
};

let userFound = (HttpRequest, HttpResponse) => {
     queries.getThreadBySlugOrIdWithVotes.values =[
        HttpRequest.params.slug,
     ];

     client.query(queries.getThreadBySlugOrIdWithVotes)
         .then(response => {
             if (response.rows.length === 0) {
                 sendThreadNotFound(HttpRequest,HttpResponse);
             } else {
                 userAndThreadFound(HttpRequest, HttpResponse, response.rows[0]);
             }
         })
         .catch(e => {
             console.log(e);
             sendError(HttpResponse);
         })
};

let userAndThreadFound = (HttpRequest, HttpResponse, thread) => {
    queries.createVote.values = [
        HttpRequest.body.voice,
        HttpRequest.body.nickname,
        thread[0],
    ];

    queries.updateVote.values = [
        HttpRequest.body.voice,
        HttpRequest.body.nickname,
        thread[0],
    ];

    client.query(queries.updateVote)
        .then(response => {
            if (response.rows.length === 0) {
                client.query(queries.createVote)
                    .then(() => {
                        sendVoteCreated(HttpRequest ,HttpResponse, thread);
                    })
                    .catch(e=> {
                        console.log(e);
                        sendError(HttpResponse)
                    })
            } else {
                sendVoteUpdated(HttpRequest, HttpResponse);
            }

        })
        .catch(e=> {
           console.log(e);
           sendError(HttpResponse)
        });
};

let sendVoteUpdated = (HttpRequest, HttpResponse) => {
    queries.getThreadBySlugOrIdWithVotes.values =[
        HttpRequest.params.slug,
    ];

    client.query(queries.getThreadBySlugOrIdWithVotes)
        .then(response => {
            let data = response.rows[0];
            HttpResponse.status(200);
            HttpResponse.json({
                author: data[5],
                created: data[3],
                id: data[0],
                forum: data[6],
                message: data[2],
                slug: data[4],
                title: data[1],
                votes: data[7],
            });
        })
};


let sendVoteCreated = (HttpRequest, HttpResponse, data) => {
    let now;
    HttpRequest.body.voice === 1 ? now = 1 : now = -1;
    HttpResponse.status(200);
    HttpResponse.json({
        author: data[5],
        created: data[3],
        id: data[0],
        forum: data[6],
        message: data[2],
        slug: data[4],
        title: data[1],
        votes: data[7] + now,
    })
};

let sendUserNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${HttpRequest.body.nickname}`
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
