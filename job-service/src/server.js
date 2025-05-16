//ScrapeSmith/job-service/src/server.js

import express from 'express';
import dotenv from 'dotenv';
import { connectMongo } from './utils/db.js';
import scheduleRoutes from './routes/scheduleRoutes.js';
import { processJobBatchQueue } from './workers/jobBatchProcessor.js';
import { retryFailedJobs } from './workers/jobRetryWorker.js';
import { runJobQueueMaintenance } from './workers/jobMaintenanceWorker.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3007;

app.use(express.json({ limit: '10mb' }));

connectMongo();

app.use('/', scheduleRoutes);

app.get('/', (req, res) => {
  res.send('Job Scheduler Service is running...');
});

app.listen(PORT, () => {
  console.log(`Job Scheduler Service running on port ${PORT}`);
});

// Batch poller every 3s
setInterval(() => {
  processJobBatchQueue().catch(err =>
    console.error('Batch delivery worker error:', err.message)
  );
}, 3000);

// Retry failed jobs every 30s
setInterval(() => {
  retryFailedJobs().catch(err =>
    console.error('Retry worker error:', err.message)
  );
}, 30000);

// Queue maintenance every 30min
setInterval(() => {
  runJobQueueMaintenance().catch(err =>
    console.error('Maintenance worker error:', err.message)
  );
}, 30 * 60 * 1000);
