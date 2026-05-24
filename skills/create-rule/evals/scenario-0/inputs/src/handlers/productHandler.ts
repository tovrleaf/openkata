import { Request, Response } from 'express';

export async function getProduct(req: Request, res: Response) {
  try {
    const product = await db.findProduct(req.params.id);
    if (!product) {
      res.status(404).json({ status: 'not_found', message: 'Product not found' });
      return;
    }
    res.json({ result: product, ok: true });
  } catch (error) {
    res.status(500).json({ status: 'error', msg: String(error) });
  }
}

export async function listProducts(req: Request, res: Response) {
  try {
    const products = await db.listProducts();
    res.json({ items: products, count: products.length });
  } catch (error) {
    res.status(500).json({ message: String(error), code: 'INTERNAL_ERROR' });
  }
}

export async function updateProduct(req: Request, res: Response) {
  try {
    const product = await db.updateProduct(req.params.id, req.body);
    if (!product) {
      return res.status(404).json({ err: 'Not found' });
    }
    return res.json({ data: product, success: true });
  } catch (error) {
    return res.status(500).json({ error: String(error) });
  }
}
