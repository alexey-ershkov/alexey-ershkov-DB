const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = async (HttpRequest, HttpResponse) => {
    queries.getPostById.values = [
        HttpRequest.params.id,
    ];

    let response = await client.query(queries.getPostById);

    if (response.rows.length === 0) {
        sendPostNotFound(HttpRequest, HttpResponse);
    } else {
        if (HttpRequest.body.message && HttpRequest.body.message !== response.rows[0][5]) {
            queries.updatePost.values = [
                HttpRequest.body.message,
                HttpRequest.params.id,
            ];

            client.query(queries.updatePost)
                .then(() => {
                    let data = response.rows[0];
                    HttpResponse.status(200).json({
                        author: data[0],
                        created: data[1],
                        forum: data[2],
                        id: data[3],
                        isEdited: true,
                        message: HttpRequest.body.message,
                        parent: data[6],
                        thread: data[7],
                    })
                })
                .catch(e => {
                    console.log(e);
                    sendError(HttpResponse);
                })
        } else {
            let data = response.rows[0];
            HttpResponse.status(200).json({
                author: data[0],
                created: data[1],
                forum: data[2],
                id: data[3],
                isEdited: data[4],
                message: data[5],
                parent: data[6],
                thread: data[7],
            })
        }
    }
};

let sendPostNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find post with id ${HttpRequest.params.id}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
