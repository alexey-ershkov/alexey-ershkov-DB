const pg = require('pg');

const client = new pg.Client({
    user: 'farcoad',
    host: 'localhost',
    database: 'forum',
    password: 'postgres',
    port: 5432,
});

client.connect(err => {
    if (err) {
        console.error('[ERROR] Database connection error', err.stack)
    } else {
        console.log('[INFO] Database connected')
    }
});

module.exports = client;
