import express from 'express';
const router = express.Router();

router.post('/login', async (req, res) => {
  const { email, password } = req.body;
  res.json({ token: 'placeholder' });
});

router.post('/logout', async (req, res) => {
  res.json({ success: true });
});

export default router;
