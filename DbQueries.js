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
