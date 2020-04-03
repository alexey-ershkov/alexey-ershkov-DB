const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getUserByNickname.values = [
        HttpRequest.params.nickname
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            if (response.rows.length !== 0) {
                updateUser(HttpRequest, HttpResponse);
            } else {
                sendNotFound(HttpRequest, HttpResponse);
            }
        })
        .catch(()=>{
            sendError(HttpResponse);
        })
};

let updateUser = (HttpRequest, HttpResponse) => {

    let insertText = 'UPDATE usr SET ';
    queries.updateUser.values = [];
    Object.keys(HttpRequest.body).forEach( (key, index) => {
        if (Object.keys(HttpRequest.body).length - 1 !== index) {
            insertText += `${key} = $${index+1}, `;
        } else {
            insertText += `${key} = $${index+1} WHERE nickname = $${index+2}`;
        }
        queries.updateUser.values.push(HttpRequest.body[key]);
    });

    queries.updateUser.values.push(HttpRequest.params.nickname);
    queries.updateUser.text = insertText;


    if (Object.keys(HttpRequest.body).length !== 0) {
        client.query(queries.updateUser)
            .then(() => {
                sendOkResponse(HttpRequest, HttpResponse);
            })
            .catch(error => {
                sendCannotUpdate(HttpRequest, HttpResponse);
            })
    } else {
        sendOkResponse(HttpRequest,HttpResponse);
    }
};

let sendOkResponse = (HttpRequest, HttpResponse) => {
    queries.getUserByNickname.values = [
        HttpRequest.params.nickname
    ];

    client.query(queries.getUserByNickname)
        .then(response => {
            sendUserInfo(HttpResponse, response.rows[0]);
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

let sendCannotUpdate = (HttpRequest, HttpResponse) => {
    HttpResponse.status(409);
    HttpResponse.json({
        message: `Can't update user with nickname ${HttpRequest.params.nickname}`
    })
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
