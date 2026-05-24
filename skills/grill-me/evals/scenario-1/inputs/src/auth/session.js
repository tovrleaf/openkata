const session = require('express-session');
const RedisStore = require('connect-redis')(session);
const redis = require('../lib/redis');

// Sessions expire after 24 hours of inactivity
const SESSION_TTL_SECONDS = 86400;

module.exports = session({
  store: new RedisStore({
    client: redis,
    ttl: SESSION_TTL_SECONDS,
  }),
  secret: process.env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  cookie: {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: SESSION_TTL_SECONDS * 1000,
  },
});
