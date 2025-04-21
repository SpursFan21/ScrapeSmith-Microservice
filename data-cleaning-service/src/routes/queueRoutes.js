//ScrapeSmith\data-cleaning-service\src\routes\queueRoutes.js
import express from 'express';
import { QueuedCleanJob } from '../models/QueuedCleanJob.js';

const router = express.Router();

router.post('/', async (req, res) => {
  try {
    const job = await QueuedCleanJob.create(req.body);
    res.status(201).json({ message: 'Job queued for cleaning', jobId: job._id });
  } catch (err) {
    console.error('Failed to queue clean job:', err);
    res.status(500).json({ error: 'Failed to queue cleaning job' });
  }
});

export default router;
