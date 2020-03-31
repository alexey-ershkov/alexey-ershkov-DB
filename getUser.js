let client = require('./connectDB');
let queries = require('./DbQueries');


module.exports = (HttpRequest, HttpResponse) => {
    queries.getUserByNickname.values = [
        HttpRequest.params.nickname,
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            if (response.rows.length !== 0)
                sendUserInfo(HttpResponse, response.rows[0]);
            else
                sendNotFound(HttpRequest, HttpResponse);
        })
        .catch( () => {
            sendError(HttpResponse);
        })
};

let sendUserInfo = (HttpResponse, data) => {
    HttpResponse.status(200);
    HttpResponse.json({
        about: data[3],
        email: data[0],
        fullname: data[1],
        nickname: data[2]
    });
};

let sendNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${HttpRequest.params.nickname}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
