//admin-service\controllers\ordersController.js
import { getMongoClient } from "../utils/mongoClient.js";

export const getAllOrders = async (req, res) => {
  try {
    const client = await getMongoClient();
    const db = client.db(process.env.MONGO_DB);

    const scrapes = await db.collection("scrapes").find().toArray();
    const cleaned = await db.collection("cleaneddatas").find().toArray();

    res.json({
      scrapes,
      cleaned,
    });
  } catch (err) {
    console.error("Error fetching all orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};

export const getOrdersByUser = async (req, res) => {
  const userId = req.params.userId;

  try {
    const client = await getMongoClient();
    const db = client.db(process.env.MONGO_DB);

    const scrapes = await db.collection("scrapes").find({ userId }).toArray();
    const cleaned = await db.collection("cleaneddatas").find({ userId }).toArray();

    res.json({
      scrapes,
      cleaned,
    });
  } catch (err) {
    console.error("Error fetching user orders:", err);
    res.status(500).json({ error: "Internal server error" });
  }
};
