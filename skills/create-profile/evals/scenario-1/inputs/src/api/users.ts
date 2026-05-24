import express from 'express';
const router = express.Router();

router.get('/', async (req, res) => {
  res.json({ users: [] });
});

router.post('/', async (req, res) => {
  res.status(201).json({ created: true });
});

export default router;
