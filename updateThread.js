const client = require('./connectDB');
const queries = require('./DbQueries');

module.exports = (HttpRequest, HttpResponse) => {
    queries.getThreadBySlugOrId.values =[
        HttpRequest.params.slug,
    ];

    client.query(queries.getThreadBySlugOrId)
        .then(response => {
            if (response.rows.length === 0) {
                sendThreadNotFound(HttpRequest, HttpResponse, response.rows[0]);
            } else {
                updateThread(HttpRequest, HttpResponse, response.rows[0]);
            }
        })
        .catch(e=> {
            console.log(e);
            sendError(HttpResponse);
        })
};

let updateThread = (HttpRequest, HttpResponse, thread) => {
      if (Object.keys(HttpRequest.body).length === 0) {
          sendThread(HttpResponse, thread);
      } else {
          let updateQuery = {
              rowMode: 'array',
              text: 'UPDATE thread SET ',
              values: []
          };
          Object.keys(HttpRequest.body).forEach((key, index) => {
              if (Object.keys(HttpRequest.body).length - 1 !== index) {
                  updateQuery.text += `${key} = $${index+1}, `;
              } else {
                  updateQuery.text += `${key} = $${index+1} WHERE id = $${index+2} RETURNING id`;
              }
              updateQuery.values.push(HttpRequest.body[key]);
          });
          updateQuery.values.push(thread[0]);

          client.query(updateQuery)
              .then(response => {
                  queries.getThreadBySlugOrId.values = [
                        response.rows[0][0],
                  ];

                  client.query(queries.getThreadBySlugOrId)
                      .then(response => {
                          sendThread(HttpResponse, response.rows[0]);
                      });
              })
              .catch(e=> {
                  console.log(e);
                  sendError(HttpResponse);
              });
      }
};

let sendThread = (HttpResponse, data) => {
    HttpResponse.status(200);
    HttpResponse.json({
        author: data[5],
        created: data[3],
        id: data[0],
        forum: data[6],
        message: data[2],
        slug: data[4],
        title: data[1],
    });
};

let sendThreadNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find thread with slug or id ${HttpRequest.params.slug}`
    })
};

let sendError = (HttpResponse) => {
    HttpResponse.status(500).send('Internal error');
};
