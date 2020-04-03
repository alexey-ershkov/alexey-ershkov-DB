const client = require('../connectDB');
const queries = require('../DbQueries');

module.exports = async (HttpRequest, HttpResponse) => {
    queries.getPostById.values = [
        HttpRequest.params.id,
    ];

    let postResponse = await client.query(queries.getPostById);
    if (postResponse.rows.length !== 0) {
        let answer = {};
        let postData = postResponse.rows[0];
        let threadData = undefined;
        if (HttpRequest.query.related){
            let related = HttpRequest.query.related.split(',');
            if (related.indexOf('user') !== -1) {
                queries.getUserByNickname.values = [
                    postData[0],
                ];

                let userResponse = await client.query(queries.getUserByNickname);
                let userData = userResponse.rows[0];
                answer.author = {
                    about: userData[3],
                    email: userData[0],
                    fullname: userData[1],
                    nickname: userData[2],
                }
            }

            if (related.indexOf('thread') !== -1) {
                queries.getThreadBySlugOrIdWithVotes.values = [
                    postData[7],
                ];

                let threadResponse = await client.query(queries.getThreadBySlugOrIdWithVotes);
                threadData = threadResponse.rows[0];
            }

            if (related.indexOf('forum') !== -1) {
                queries.getForumBySlug.values = [
                    postData[2],
                ];

                let forumResponse = await client.query(queries.getForumBySlug);
                let forumData = forumResponse.rows[0];
                answer.forum = {
                    posts: Number(forumData[0]),
                    slug: forumData[1],
                    threads: Number(forumData[2]),
                    title: forumData[3],
                    user: forumData[4]
                }
            }

        }

        answer.post = {
            author: postData[0],
            created: postData[1],
            forum: postData[2],
            id: postData[3],
            isEdited: Boolean(postData[4]),
            message: postData[5],
            parent: postData[6],
            thread: postData[7],
        };

        if (threadData) {
            answer.thread = {
                author: threadData[5],
                created: threadData[3],
                forum: threadData[6],
                id: threadData[0],
                message: threadData[2],
                slug: threadData[4],
                title: threadData[1],
            };
        }

        HttpResponse.status(200).json(answer);
    } else {
        sendPostNotFound(HttpRequest, HttpResponse);
    }
};

let sendPostNotFound = (HttpRequest, HttpResponse) => {
    HttpResponse.status(404);
    HttpResponse.json({
        message: `Can't find post with id ${HttpRequest.params.id}`
    })
};
