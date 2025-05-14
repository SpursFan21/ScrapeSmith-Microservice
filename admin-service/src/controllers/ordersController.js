//admin-service\controllers\ordersController.js

import { Scrape } from '../models/Scrape.js';
import { CleanedData } from '../models/CleanedData.js';

// GET /admin/orders
export const getAllOrders = async (req, res) => {
  try {
    const [scrapes, cleaned] = await Promise.all([
      Scrape.find(),
      CleanedData.find()
    ]);

    res.json({ scrapes, cleaned });
  } catch (err) {
    console.error("Error fetching all orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};

// GET /admin/orders/:userId
export const getOrdersByUser = async (req, res) => {
  const userId = req.params.userId;

  try {
    const [scrapes, cleaned] = await Promise.all([
      Scrape.find({ user_id: userId }),
      CleanedData.find({ userId })
    ]);

    res.json({ scrapes, cleaned });
  } catch (err) {
    console.error("Error fetching user orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};
