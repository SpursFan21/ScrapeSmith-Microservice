//ScrapeSmith\job-service\src\workers\jobRetryWorker.js

import { ScheduledJob } from '../models/ScheduledJob.js';

export async function retryFailedJobs() {
  const result = await ScheduledJob.updateMany(
    {
      status: 'failed',
      attempts: { $lt: 3 },
      lastTriedAt: { $lte: new Date(Date.now() - 30 * 1000) } // Retry after 30s
    },
    {
      $set: {
        status: 'scheduled',
        lastTriedAt: new Date()
      },
      $inc: { attempts: 1 }
    }
  );

  if (result.modifiedCount > 0) {
    console.log(`Retried ${result.modifiedCount} failed job(s)`);
  }
}