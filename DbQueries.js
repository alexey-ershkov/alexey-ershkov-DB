module.exports.getUserByNickname = {
    text: 'SELECT * FROM usr WHERE nickname = $1',
    rowMode: 'array',
};

module.exports.insertUser = {
    text: 'INSERT INTO usr (email, fullname, nickname, about) VALUES ($1, $2, $3, $4)',
    rowMode: 'array',
};

module.exports.getUser = {
    text: 'SELECT * FROM usr WHERE nickname = $1 OR email = $2',
    rowMode: 'array',
};

module.exports.updateUser = {
    rowMode: 'array',
};

module.exports.createForum = {
  rowMode: 'array',
  text: 'INSERT INTO forum (slug, title, usr) VALUES ($1, $2, $3)',
};

module.exports.getForumBySlug = {
    rowMode: 'array',
    text: 'SELECT count(p), f.slug, (SELECT count(*) FROM forum f2 JOIN thread t2 on f2.slug = t2.forum WHERE f2.slug = $1), f.title, u.nickname FROM forum f\n' +
        'LEFT JOIN thread t on f.slug = t.forum\n' +
        'LEFT JOIN post p on t.id = p.thread\n' +
        'JOIN usr u on f.usr = u.nickname\n'+
        'WHERE f.slug = $1' +
        'GROUP BY f.slug, u.nickname'
};

module.exports.getForumBySlugSimple = {
    rowMode: 'array',
    text: 'SELECT f.slug, f.title, u.nickname FROM forum f\n' +
        'JOIN usr u on f.usr = u.nickname\n'+
        'WHERE f.slug = $1'
};

module.exports.createThread = {
    rowMode: 'array',
    text: 'INSERT INTO thread (usr, created, forum, message, title, slug) VALUES ($1, $2, $3, $4, $5, $6)'
};

module.exports.getThread = {
    rowMode: 'array',
    text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  FROM thread t JOIN forum f on t.forum = f.slug\n' +
        'WHERE t.usr = $1 AND t.forum = $2 AND t.message = $3 AND t.title = $4'
};

module.exports.getThreadBySlug = {
  rowMode: 'array',
  text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  FROM thread t JOIN forum f on t.forum = f.slug\n' +
        'WHERE t.slug = $1'
};

module.exports.getThreadBySlugOrId = {
    rowMode: 'array',
    text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug  FROM thread t JOIN forum f on t.forum = f.slug\n' +
        'WHERE t.slug = $1 OR t.id::citext = $1'
};

module.exports.getThreadsInForumBySlagDESC = {
    rowMode: 'array',
    text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug FROM thread t\n' +
        'JOIN forum f on t.forum = f.slug\n' +
        'WHERE f.slug = $1 AND t.created <=  $2::timestamp AT TIME ZONE \'0\'\n' +
        'ORDER BY t.created DESC'
};


module.exports.getThreadsInForumBySlagASC = {
    rowMode: 'array',
    text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug FROM thread t\n' +
        'JOIN forum f on t.forum = f.slug\n' +
        'WHERE f.slug = $1 AND t.created >=  $2::timestamp AT TIME ZONE \'0\'\n' +
        'ORDER BY t.created'
};

module.exports.getThreadBySlugOrIdWithVotes = {
    rowMode: 'array',
    text: 'SELECT t.id, t.title, t.message, t.created, t.slug, t.usr, f.slug, SUM(v.vote)::integer FROM thread t\n' +
        'JOIN forum f on t.forum = f.slug\n' +
        'LEFT JOIN vote v on t.id = v.thread\n' +
        'WHERE t.slug = $1 OR t.id::citext = $1\n' +
        'GROUP BY f.slug, t.id'
};

module.exports.createVote = {
    rowMode: 'array',
    text: 'INSERT INTO vote (vote, usr, thread) VALUES ($1::integer , $2, $3)'
};

module.exports.updateVote = {
    rowMode: 'array',
    text: 'UPDATE vote SET vote = $1 WHERE usr = $2 AND thread = $3\n' +
        'RETURNING thread'
};

module.exports.createSinglePost = {
    rowMode: 'array',
    text: 'INSERT INTO post (usr, message,  parent, thread, created) VALUES ($1, $2, $3, $4, $5) RETURNING id'
};

module.exports.createPosts = {
    rowMode: 'array',
    text: 'INSERT INTO post (usr, message,  parent, thread, created) VALUES'
};

module.exports.getPostsByIds = {
    rowMode: 'array',
    text: 'SELECT u.nickname, p.created, f.slug, p.id, p.isEdited, p.message, p.parent, t.id FROM post p\n' +
        'JOIN thread t on p.thread = t.id\n' +
        'JOIN forum f on t.forum = f.slug\n' +
        'JOIN usr u on p.usr = u.nickname\n' +
        'WHERE p.id IN '
};

module.exports.getParentThread = {
    rowMode: 'array',
    text: 'SELECT thread FROM post\n' +
        'WHERE id = $1 '
};

module.exports.getInfo = {
    rowMode:'array',
    text: 'SELECT count(*) FROM forum\n' +
        'UNION ALL\n' +
        'SELECT count(*) FROM post\n' +
        'UNION ALL\n' +
        'SELECT count(*) FROM thread\n' +
        'UNION ALL\n' +
        'SELECT count(*) FROM usr'
};

module.exports.clearDB = {
  text: 'DELETE FROM usr'
};

module.exports.getTimestamp = {
    text: 'SELECT current_timestamp'
};
