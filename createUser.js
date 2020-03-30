const client = require('./connectDB');

let query = {
    text: 'INSERT INTO usr (email, fullname, nickname, about) VALUES ($1, $2, $3, $4)',
    rowMode: 'array',
};


module.exports = (HttpRequest, HttpResponse) => {


    query.values = [
        HttpRequest.body.email,
        HttpRequest.body.fullname,
        HttpRequest.params.nickname,
        HttpRequest.body.about
    ];

    HttpResponse.status(201);

    client.query(query, (DBError, DBResponse) => {
        CallbackDB(
            HttpRequest,
            HttpResponse,
            DBError,
            DBResponse
        );
    })
};

let CallbackDB = (
        HttpRequest,
        HttpResponse,
        DbError,
        DbResponse) => {
    if (DbError) {
        console.log(DbError.stack);
        HttpResponse.status(500).send('Ops');
    } else {
        HttpResponse.json({
            about: HttpRequest.body.about,
            email: HttpRequest.body.email,
            fullname: HttpRequest.body.fullname,
            nickname: HttpRequest.params.nickname
        });
    }};

