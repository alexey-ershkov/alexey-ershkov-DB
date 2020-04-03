const client = require('../connectDB');
const queries = require('../DbQueries');


module.exports = (HttpRequest, HttpResponse) => {

    queries.insertUser.values = [
        HttpRequest.body.email,
        HttpRequest.body.fullname,
        HttpRequest.params.nickname,
        HttpRequest.body.about
    ];


    client.query(queries.insertUser)
        .then(()=>{
            sendOkResponse(HttpRequest, HttpResponse);
        })
        .catch( () => {

            queries.getUser.values = [
                HttpRequest.params.nickname,
                HttpRequest.body.email,
            ];

            client.query(queries.getUser)
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

