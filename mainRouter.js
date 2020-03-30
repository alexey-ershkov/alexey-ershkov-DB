const express = require('express');
const mainRouter = express.Router();
const createUser = require('./createUser');

mainRouter.use((request, response, next) => {
    response.set('Content-Type', 'application/json');
    console.log(`[DEBUG] Request URL is http://localhost:3000${request.path}`);
    next();
});

mainRouter.post('/user/:nickname/create',createUser);

module.exports = mainRouter;
