const client = require('./connectDB');

let inserUser = {
    text: 'INSERT INTO usr (email, fullname, nickname, about) VALUES ($1, $2, $3, $4)',
    rowMode: 'array',
};

let getUser = {
    text: 'SELECT * FROM usr WHERE nickname = $1 OR email = $2',
    rowMode: 'array',
};



module.exports = (HttpRequest, HttpResponse) => {

    inserUser.values = [
        HttpRequest.body.email,
        HttpRequest.body.fullname,
        HttpRequest.params.nickname,
        HttpRequest.body.about
    ];


    client.query(inserUser)
        .then(()=>{
            sendOkResponse(HttpRequest, HttpResponse);
        })
        .catch( () => {

            getUser.values = [
                HttpRequest.params.nickname,
                HttpRequest.body.email,
            ]

            client.query(getUser)
                .then(response => {
                    sendRepeatedStatus(HttpResponse, response.rows);
                })
                .catch(()=>{
                    sendError(HttpResponse);
                })
        });

};

let sendOkResponse = (HttpRequest, HttpResponse) => {
    HttpResponse.status(201);
    HttpResponse.json({
        about: HttpRequest.body.about,
        email: HttpRequest.body.email,
        fullname: HttpRequest.body.fullname,
        nickname: HttpRequest.params.nickname
    });
};

let sendRepeatedStatus = (HttpResponse, data) => {
    HttpResponse.status(409);
    let response = [];
    data.forEach(elem => {
       response.push({
           about: elem[3],
           email: elem[0],
           fullname: elem[1],
           nickname: elem[2]
       })
    });
    HttpResponse.json(response);
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};

