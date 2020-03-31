const express = require('express');
const mainRouter = express.Router();

//User Handlers
const createUser = require('./UserHandlers/createUser');
const updateUser = require('./UserHandlers/updateUser');
const getUser = require('./UserHandlers/getUser');

//Forum Handlers
const createForum = require('./createForum');
const getForum = require('./getForum');

//Thread Handlers
const createThread = require('./createThread');

mainRouter.use((request, response, next) => {
    response.set('Content-Type', 'application/json');
    console.log(`[DEBUG] ${request.method}: Request URL is http://localhost:3000${request.path}`);
    next();
});

//User URL section
mainRouter.get('/user/:nickname/profile', getUser);
mainRouter.post('/user/:nickname/profile', updateUser);
mainRouter.post('/user/:nickname/create',createUser);

//Forum URL section
mainRouter.post('/forum/create', createForum);
mainRouter.get('/forum/:slug/details/', getForum);

//Thread URL selection
mainRouter.post('/forum/:forum/create', createThread);


module.exports = mainRouter;
