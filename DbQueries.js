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
    text: 'SELECT count(p), f.slug, count(t), f.title, f.usr FROM forum f\n' +
        'LEFT JOIN post p on f.slug = p.forum\n' +
        'LEFT JOIN thread t on f.slug = t.forum\n' +
        'WHERE f.slug = $1' +
        'GROUP BY f.slug'
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
