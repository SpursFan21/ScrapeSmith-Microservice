//ScrapeSmith\ai-analysis-service\src\server.js

import express from 'express';
import dotenv from 'dotenv';
import { connectMongo } from './utils/db.js';
import analysisRoutes from './routes/analysisRoutes.js';
import queueRoutes from './routes/queueRoutes.js';
import { processAIBatchQueue } from "./workers/aiBatchProcessor.js";
import { retryFailedAIJobs } from "./workers/aiRetryWorker.js";
import { runAIMaintenance } from "./workers/aiMaintenanceWorker.js";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3006;

app.use(express.json());
app.use('/api/analyize', analysisRoutes);
app.use('/api/queue', queueRoutes);

// Health check
app.get('/', (req, res) => {
  res.send('AI Analysis Service is running...');
});

// Connect to Mongo and then start server and workers
connectMongo().then(() => {
  app.listen(PORT, () => {
    console.log(`AI Analysis Service running on port ${PORT}`);
  });

  // Batch poller every 3s
  setInterval(() => {
    processAIBatchQueue().catch(err =>
      console.error("Batch AI worker error:", err.message)
    );
  }, 3000);

  // Retry failed jobs every 30s
  setInterval(() => {
    retryFailedAIJobs().catch(err =>
      console.error("Retry worker error:", err.message)
    );
  }, 30000);

  // Maintenance every 30 minutes
  setInterval(() => {
    runAIMaintenance().catch(err =>
      console.error("Maintenance worker error:", err.message)
    );
  }, 30 * 60 * 1000);
});
