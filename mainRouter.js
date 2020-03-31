const express = require('express');
const mainRouter = express.Router();
const createUser = require('./createUser');
const updateUser = require('./updateUser');
const getUser = require('./getUser');

mainRouter.use((request, response, next) => {
    response.set('Content-Type', 'application/json');
    console.log(`[DEBUG] Request URL is http://localhost:3000${request.path}`);
    next();
});

mainRouter.get('/user/:nickname/profile', getUser);
mainRouter.post('/user/:nickname/profile', updateUser);
mainRouter.post('/user/:nickname/create',createUser);

module.exports = mainRouter;
