const express = require('express');
const mainRouter = express.Router();

//User Handlers
const createUser = require('./UserHandlers/createUser');
const updateUser = require('./UserHandlers/updateUser');
const getUser = require('./UserHandlers/getUser');

//Forum Handlers
const createForum = require('./ForumHandlers/createForum');
const getForum = require('./ForumHandlers/getForum');
const getForumThreads = require('./ForumHandlers/getForumThreads');
const getForumUsers = require('./getForumUsers');

//Thread Handlers
const createThread = require('./ThreadHandlers/createThread');
const getThreadInfo = require('./ThreadHandlers/getThreadInfo');
const createVote = require('./ThreadHandlers/createVote');
const updateThread = require('./ThreadHandlers/updateThread');
const getThreadPosts = require('./getThreadPosts');

//Post Handlers
const createPost = require('./createPost');
const updatePost = require('./updatePost');
const getPost = require('./getPost');

//Info Handlers
const getInfo = require('./ServiceHandlers/getInfo');

//Clear Handlers
const clearDB = require('./ServiceHandlers/clearDB');

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
mainRouter.get('/forum/:forum/threads', getForumThreads);
mainRouter.get('/forum/:slug/users', getForumUsers);

//Thread URL selection
mainRouter.post('/forum/:forum/create', createThread);
mainRouter.get('/thread/:slug/details', getThreadInfo);
mainRouter.post('/thread/:slug/details', updateThread);
mainRouter.post('/thread/:slug/vote', createVote);
mainRouter.get('/thread/:slug_or_id/posts', getThreadPosts);

//Post URL section
mainRouter.post('/thread/:slug_or_id/create', createPost);
mainRouter.get('/post/:id/details', getPost);
mainRouter.post('/post/:id/details', updatePost);

//Info URL section
mainRouter.get('/service/status', getInfo);

//Clear URL section
mainRouter.post('/service/clear', clearDB);

module.exports = mainRouter;
