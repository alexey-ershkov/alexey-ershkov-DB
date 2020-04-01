const client = require('./connectDB');
const queries = require('./DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getThreadBySlugOrId.values = [
        HttpRequest.params.slug_or_id,
    ];

    client.query(queries.getThreadBySlugOrId)
        .then(response => {
            if (response.rows.length === 0) {
                sendThreadNotFound(HttpRequest, HttpResponse)
            } else {
                createPosts(HttpRequest, HttpResponse, response.rows[0]);
            }
        })
        .catch(e=> {
            console.log('err:');
            console.log(e);
            sendError(HttpResponse);
        })
};

let createPosts = (HttpRequest, HttpResponse, threadInfo) => {

    let createQuery = Object.assign({}, queries.createPosts);
    let globalQueryCount = 0;
    createQuery.values = [];
    if (HttpRequest.body.length !== 0) {
        let users = [];
        let threads = [];
        HttpRequest.body.forEach(elem => {
            users.push(elem.author);
            if (!elem.parent) {
                elem.parent = 0;
            }
            threads.push(elem.parent);
        });

        let GetUsersQuery = doSearchQuery(users, queries.getUsers);
        let GetThreadQuery = doSearchQuery(threads, queries.getParentThread);

        client.query(GetUsersQuery)
            .then(response => {
                if (false) {
                    sendUserNotFound(HttpResponse)
                } else {
                    client.query(GetThreadQuery)
                        .then(response => {
                            if (checkThread(response.rows, threadInfo[0])) {
                                sendAnotherThreadError(HttpResponse);
                            } else {
                                HttpRequest.body.forEach((elem, ind) => {
                                    createQuery.text += '(';
                                    Object.keys(elem).forEach((key, index) => {
                                        globalQueryCount++;
                                        if (Object.keys(elem).length - 1 !== index) {
                                            createQuery.text += `$${globalQueryCount}, `;
                                        } else {
                                            createQuery.text += `$${globalQueryCount++}, $${globalQueryCount}, current_timestamp)`
                                        }
                                        createQuery.values.push(elem[key]);
                                    });
                                    createQuery.values.push(threadInfo[0]);
                                    if (HttpRequest.body.length - 1 !== ind) {
                                        createQuery.text += ', ';
                                    } else {
                                        createQuery.text += 'RETURNING id';
                                    }
                                });
                                client.query(createQuery)
                                    .then(response => {
                                        sendPostInfo(HttpResponse, response.rows);
                                    })
                                    .catch(e => {
                                        console.log(e);
                                        sendError(HttpResponse);
                                    })
                            }
                        })
                }
            });
    } else {
        HttpResponse.status(201).json([]);
    }
};

let checkThread = (threads, threadId) => {
    console.log(threads, threadId);
    let out = false;
    threads.forEach(key => {
        if (key[0] !== threadId) {
            return out = true;
        }
    });
    return out;
};

let sendPostInfo = (HttpResponse, PostIds) => {
      let getPostsQuery = Object.assign({}, queries.getPostsByIds);
      getPostsQuery.text += ' (';
      getPostsQuery.values = [];
      PostIds.forEach((key, index) =>{
          getPostsQuery.values.push(key[0]);
          getPostsQuery.text += ` $${index+1}`;
          if (PostIds.length - 1 !== index){
              getPostsQuery.text += ','
          } else {
              getPostsQuery.text += ')'
          }
      });

      client.query(getPostsQuery)
          .then(DB => {
              let response = [];
              DB.rows.forEach(data => {
                  response.push({
                      author: data[0],
                      created: data[1],
                      forum: data[2],
                      id: data[3],
                      message: data[4],
                      parent: data[5],
                      thread: data[6],
                  })
              });
              HttpResponse.status(201).json(response);
          })
          .catch(e => {
              console.log(e);
              sendError(HttpResponse);
          })
};
let sendThreadNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with slug or id ${HttpRequest.params.slug_or_id}`
    })
};

let sendUserNotFound = (HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};

let sendAnotherThreadError = HttpResponse => {
    HttpResponse.status(409).json({
        message: 'Parent post was created in another thread',
    })
};

let doSearchQuery = function (array, DBQuery) {
    let query = Object.assign({}, DBQuery);
    query.text += ' (';
    query.values = [];

    array.forEach((key, index) =>{
        query.values.push(key);
        query.text += ` $${index+1}`;
        if (array.length - 1 !== index){
            query.text += ','
        } else {
            query.text += ')'
        }
    });

    return query;
};
