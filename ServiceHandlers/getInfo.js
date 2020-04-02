const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    client.query(queries.getInfo)
        .then(resp => {
            const data = resp.rows;
            HttpResponse.status(200).json({
                forum: Number(data[0][0]),
                post: Number(data[1][0]),
                thread: Number(data[2][0]),
                user: Number(data[3][0]),
            })
        })
};
