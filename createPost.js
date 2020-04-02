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
            console.log(e);
            sendError(HttpResponse);
        })
};

let createPosts = async (HttpRequest, HttpResponse, threadInfo) => {


    if (HttpRequest.body.length !== 0) {
        let isSend = false;
        for (const elem of HttpRequest.body) {
            queries.getUserByNickname.values = [
                elem.author,
            ];

            let resp = await client.query(queries.getUserByNickname);
            if (resp.rows.length === 0) {
                isSend = true;
                sendUserNotFound(HttpResponse, elem.author);
            }

            if (!elem.parent) {
                elem.parent = 0;
            } else {
                queries.getParentThread.values = [
                    elem.parent,
                ];

                resp = await client.query(queries.getParentThread);
                if (resp.rows.length===0 || resp.rows[0][0] !== threadInfo[0]) {
                    isSend = true;
                    sendAnotherThreadError(HttpResponse)
                }
            }


        }

        if (!isSend) {
            let ids = [];
            client.query('BEGIN')
                .then(async () => {
                    for (const elem of HttpRequest.body) {
                        queries.createSinglePost.values = [
                            elem.author,
                            elem.message,
                            elem.parent,
                            threadInfo[0],
                        ];
                        let resp = await client.query(queries.createSinglePost);
                        ids.push(resp.rows[0]);
                    }
                    await client.query('COMMIT');
                    sendPostInfo(HttpResponse,ids);
                });
        }
    } else {
        HttpResponse.status(201).json([]);
    }
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

let doInputQuery = (HttpRequest, threadId) => {
    let createQuery = Object.assign({}, queries.createPosts);
    let globalQueryCount = 0;
    createQuery.values = [];
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
        createQuery.values.push(threadId);
        if (HttpRequest.body.length - 1 !== ind) {
            createQuery.text += ', ';
        } else {
            createQuery.text += 'RETURNING id';
        }
    });

    return createQuery;
};

let sendUserNotFound = (HttpResponse, nickanme) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find user with nickname ${nickanme}`,
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

