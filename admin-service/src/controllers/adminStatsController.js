//admin-service\src\controllers\adminStatsController.js
import pool from '../utils/db.js';
import { getMongoClient } from '../utils/mongoClient.js';

export const getAdminStats = async (req, res) => {
  try {
    // 1. PostgreSQL queries
    const userCountRes = await pool.query('SELECT COUNT(*) FROM users');
    const adminCountRes = await pool.query('SELECT COUNT(*) FROM users WHERE is_admin = TRUE');

    const totalUsers = parseInt(userCountRes.rows[0].count, 10);
    const totalAdmins = parseInt(adminCountRes.rows[0].count, 10);

    // 2. MongoDB queries
    const client = await getMongoClient();
    const db = client.db(process.env.MONGO_DB);

    const totalOrders = await db.collection('scrapes').countDocuments();
    const openTickets = await db.collection('tickets').countDocuments({ status: 'open' });

    // 3. Response
    res.json({
      totalUsers,
      totalAdmins,
      totalOrders,
      openTickets,
    });
  } catch (error) {
    console.error('Error fetching admin stats:', error);
    res.status(500).json({ error: 'Failed to load admin stats' });
  }
};
