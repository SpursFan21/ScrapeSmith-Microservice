//admin-service\controllers\ordersController.js

import { Scrape } from '../models/ScrapeData.js';
import { CleanedData } from '../models/CleanedData.js';
import { AnalyzedData } from '../models/AnalyzedData.js';

// GET /admin/orders
export const getAllOrders = async (req, res) => {
  try {
    const [scrapes, cleaned, analyzed] = await Promise.all([
      Scrape.find(),
      CleanedData.find(),
      AnalyzedData.find()
    ]);

    res.json({ scrapes, cleaned, analyzed });
  } catch (err) {
    console.error("Error fetching all orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};

// GET /admin/orders/:userId
export const getOrdersByUser = async (req, res) => {
  const userId = req.params.userId;

  try {
    const [scrapes, cleaned, analyzed] = await Promise.all([
      Scrape.find({ user_id: userId }),
      CleanedData.find({ userId }),
      AnalyzedData.find({ userId }) 
    ]);

    res.json({ scrapes, cleaned, analyzed });
  } catch (err) {
    console.error("Error fetching user orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};
