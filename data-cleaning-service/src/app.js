// data-cleaning-service/src/app.js
import express from 'express';
import cleaningRoutes from './routes/cleaningRoutes.js';
import { connectMongo } from './utils/db.js';
import dotenv from 'dotenv';
import queueRoutes from './routes/queueRoutes.js';
import { processCleanBatchQueue } from './workers/cleaningBatchProcessor.js';
import { retryFailedCleanJobs } from './workers/cleaningRetryWorker.js';
import { runCleanQueueMaintenance } from './workers/cleaningMaintenanceWorker.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3004;

app.use(express.json({ limit: '10mb' }));

connectMongo();

app.use('/api/clean', cleaningRoutes);
app.use('/api/clean/queue', queueRoutes);

// Batch cleaner: every 3 seconds
setInterval(() => {
  processCleanBatchQueue().catch(err => console.error('Batch clean queue error:', err));
}, 3000);

// Retry worker: every 30 seconds
setInterval(() => {
  retryFailedCleanJobs().catch(err => console.error('Retry worker error:', err));
}, 30000);

// Maintenance: every 30 minutes
setInterval(() => {
  runCleanQueueMaintenance().catch(err => console.error('Maintenance worker error:', err));
}, 30 * 60 * 1000);

app.listen(PORT, () => {
  console.log(`Data Cleaning Service running on port ${PORT}`);
});
