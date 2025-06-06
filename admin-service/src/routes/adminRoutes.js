//admin-service\src\routes\adminRoutes.js
import express from 'express';
import {
  getAllUsers,
  editUser,
  deleteUser
} from '../controllers/adminController.js';

import {
  getAllOrders,
  getOrdersByUser
} from '../controllers/ordersController.js';

import * as ticketController from '../controllers/ticketController.js';
import { getAdminStats } from '../controllers/adminStatsController.js';
import { getAllAIResults } from '../controllers/aiController.js';

import { verifyAdmin } from '../middleware/authMiddleware.js';

const router = express.Router();

router.use(verifyAdmin);

// User routes
router.get('/users', getAllUsers);
router.put('/users/:id', editUser);
router.delete('/users/:id', deleteUser);

// Order routes
router.get('/orders', getAllOrders);
router.get('/orders/:userId', getOrdersByUser);

// Ticket routes
router.get('/tickets', ticketController.getAllTickets);
router.get('/tickets/:id', ticketController.getTicketById);
router.post('/tickets/:id/respond', ticketController.respondToTicket);
router.post('/tickets/:id/close', ticketController.closeTicket);

// Stats route
router.get('/stats', getAdminStats);

// AI Analysis route
router.get('/ai-results', getAllAIResults);

export default router;
