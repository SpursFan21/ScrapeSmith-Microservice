//admin-service\src\controllers\adminStatsController.js

import { User } from '../models/User.js';
import { Ticket } from '../models/ticketModel.js';
import mongoose from 'mongoose';

export const getAdminStats = async (req, res) => {
  try {
    const db = mongoose.connection.db;

    const [totalUsers, totalAdmins, totalOrders, openTickets] = await Promise.all([
      User.countDocuments(),
      User.countDocuments({ is_admin: true }),
      db.collection('scraped_data').countDocuments(),
      Ticket.countDocuments({ status: 'open' }),
    ]);

    res.json({ totalUsers, totalAdmins, totalOrders, openTickets });
  } catch (error) {
    console.error('Error fetching admin stats:', error);
    res.status(500).json({ error: 'Failed to load admin stats' });
  }
};
