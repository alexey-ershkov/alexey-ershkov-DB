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
            client.query(queries.getTimestamp)
                .then(async response => {
                    for (const elem of HttpRequest.body) {
                        queries.createSinglePost.values = [
                            elem.author,
                            elem.message,
                            elem.parent,
                            threadInfo[0],
                            response.rows[0].current_timestamp,
                        ];
                        let resp = await client.query(queries.createSinglePost);
                        ids.push(resp.rows[0]);
                    }
                    sendPostInfo(HttpRequest, HttpResponse,ids);
                });
        }
    } else {
        HttpResponse.status(201).json([]);
    }
};



let sendPostInfo = async (HttpRequest, HttpResponse, PostIds) => {
      let getPostsQuery = Object.assign({}, queries.getPostsByIds);
      getPostsQuery.text += ' (';
      getPostsQuery.values = [];
      PostIds.forEach((key, index) =>{
          getPostsQuery.values.push(key[0]);
          getPostsQuery.text += ` $${index+1}`;
          if (PostIds.length - 1 !== index){
              getPostsQuery.text += ','
          } else {
              getPostsQuery.text += `) AND p.message = $${index+2}`
          }
      });
    let response = [];
        for (const elem of HttpRequest.body) {
            getPostsQuery.values.push(elem.message);
            let resp = await client.query(getPostsQuery);
            let data = resp.rows[0];
            response.push({
                 author: data[0],
                 created: data[1],
                 forum: data[2],
                 id: data[3],
                 isEdited: Boolean(data[4]),
                 message: data[5],
                 parent: data[6],
                thread: data[7],
            });
            getPostsQuery.values.pop();
        }

    HttpResponse.status(201).json(response);
};

let sendThreadNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with slug or id ${HttpRequest.params.slug_or_id}`
    })
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

