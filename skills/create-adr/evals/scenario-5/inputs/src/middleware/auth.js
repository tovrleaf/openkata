const session = require('express-session');
const RedisStore = require('connect-redis').default;
const redis = require('redis');

const redisClient = redis.createClient({
  url: process.env.REDIS_URL || 'redis://localhost:6379'
});

redisClient.connect().catch(console.error);

const sessionMiddleware = session({
  store: new RedisStore({ client: redisClient }),
  secret: process.env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  cookie: {
    secure: process.env.NODE_ENV === 'production',
    httpOnly: true,
    maxAge: 30 * 60 * 1000  // 30 minutes
  }
});

function requireAuth(req, res, next) {
  if (!req.session || !req.session.serviceId) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  next();
}

module.exports = { sessionMiddleware, requireAuth };
