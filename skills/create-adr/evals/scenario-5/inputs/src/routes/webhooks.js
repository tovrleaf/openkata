const express = require('express');
const router = express.Router();
const { requireAuth } = require('../middleware/auth');

router.use(requireAuth);

router.post('/register', async (req, res) => {
  const { url, events } = req.body;
  // Register webhook endpoint for calling service
  const serviceId = req.session.serviceId;
  res.json({ webhookId: 'wh_001', serviceId, url, events });
});

router.delete('/:id', async (req, res) => {
  res.json({ deleted: true });
});

module.exports = router;
