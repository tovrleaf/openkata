const express = require('express');
const router = express.Router();
const { processPayment } = require('../services/payment-service');

// v2 payment endpoint
router.post('/v2/payment', async (req, res) => {
  const { amount, currency, customerId, paymentMethodId } = req.body;

  if (!amount || !currency || !customerId || !paymentMethodId) {
    return res.status(400).json({ error: 'Missing required fields' });
  }

  try {
    const result = await processPayment({ amount, currency, customerId, paymentMethodId });
    res.json({ success: true, transactionId: result.id });
  } catch (err) {
    res.status(500).json({ error: 'Payment processing failed', details: err.message });
  }
});

module.exports = router;
