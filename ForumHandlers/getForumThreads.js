const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {

    queries.getForumBySlug.values = [
        HttpRequest.params.forum,
    ];

    client.query(queries.getForumBySlug)
        .then(resp => {
            if (resp.rows.length === 0) {
                sendForumNotFound(HttpRequest, HttpResponse);
            } else {
                getThreads(HttpRequest, HttpResponse);
            }
        })
        .catch(e => {
           console.log(e);
           sendError(HttpResponse);
        });
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};

let sendForumNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with forum ${HttpRequest.params.forum}`
    })
};

let getThreads = (HttpRequest, HttpResponse) => {

    let customQuery;
    if (HttpRequest.query.desc === 'true') {
        customQuery = Object.assign({}, queries.getThreadsInForumBySlagDESC);
    } else {
        customQuery = Object.assign({}, queries.getThreadsInForumBySlagASC);
    }


    customQuery.values = [
        HttpRequest.params.forum,
    ];

    if (HttpRequest.query.since) {
        customQuery.values.push(HttpRequest.query.since)
    } else {
        if (HttpRequest.query.desc === 'true') {
            customQuery.values.push('infinity')
        } else {
            customQuery.values.push('-infinity')
        }
    }



    if (HttpRequest.query.limit) {
        customQuery.text += ` LIMIT ${HttpRequest.query.limit}`;
    }

    client.query(customQuery)
        .then(response => {
            sendThreads(HttpResponse, response);
        })
        .catch(e => {
            console.log(e);
            sendError(HttpResponse);
        });
};

let sendThreads = (HttpResponse, DB) => {
    HttpResponse.status(200);
    let response = [];
    DB.rows.forEach(data => {
       response.push({
           author: data[5],
           created: data[3],
           id: data[0],
           forum: data[6],
           message: data[2],
           slug: data[4],
           title: data[1],
       })
    });
    HttpResponse.json(response);
};
