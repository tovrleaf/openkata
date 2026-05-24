import { Request, Response } from 'express';

export async function getUser(req: Request, res: Response) {
  try {
    const userId = req.params.id;
    const user = await db.findUser(userId);
    if (!user) {
      return res.status(404).json({ message: 'User not found', code: 404 });
    }
    return res.json({ data: user });
  } catch (err) {
    return res.status(500).json({ error: 'Internal error', details: String(err) });
  }
}

export async function createUser(req: Request, res: Response) {
  try {
    const { name, email } = req.body;
    if (!name || !email) {
      return res.status(400).json({ msg: 'Missing required fields' });
    }
    const user = await db.createUser({ name, email });
    return res.status(201).json(user);
  } catch (err) {
    return res.status(500).json({ error: String(err) });
  }
}

export async function deleteUser(req: Request, res: Response) {
  try {
    const { id } = req.params;
    await db.deleteUser(id);
    return res.status(200).json({ ok: true });
  } catch (err) {
    return res.status(500).json({ status: 'error', message: String(err) });
  }
}
