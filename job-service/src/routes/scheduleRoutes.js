//ScrapeSmith\job-service\src\routes\scheduleRoutes.js

// src/routes/scheduleRoutes.js
import express from 'express';
import { scheduleJob } from '../handlers/scheduleHandler.js';

const router = express.Router();

router.post('/', scheduleJob);

export default router;
