const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    client.query(queries.clearDB);
    HttpResponse.status(200).send('Database cleared');
};
