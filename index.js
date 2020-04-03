const mainRouter = require('./mainRouter');
const express = require('express');
const bodyParser = require('body-parser');
const app = express();
const port = 5000;

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(bodyParser.raw());

app.use('/api/', mainRouter);

app.listen(port, (err) => {
    if (err) {
        return console.log(`[ERROR] Can't start server on PORT:${port}`, err)
    }
    console.log(`[INFO] server is listening on http://localhost:${port}`)
});

