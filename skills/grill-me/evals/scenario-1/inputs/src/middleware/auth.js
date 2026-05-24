/**
 * Auth middleware — verifies session is valid and attaches user to req.
 * Used on all protected routes.
 */
function requireAuth(req, res, next) {
  if (!req.session || !req.session.userId) {
    return res.status(401).json({ error: 'Unauthorized' });
  }

  // Re-validate session hasn't been invalidated server-side
  // (e.g. after password change or forced logout)
  const isRevoked = req.app.locals.sessionBlocklist.has(req.session.id);
  if (isRevoked) {
    return res.status(401).json({ error: 'Session revoked' });
  }

  req.user = { id: req.session.userId, role: req.session.role };
  next();
}

module.exports = { requireAuth };
