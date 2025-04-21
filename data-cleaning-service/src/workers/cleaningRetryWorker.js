//ScrapeSmith\data-cleaning-service\src\workers\cleaningRetryWorker.js

import { QueuedCleanJob } from '../models/QueuedCleanJob.js';

export async function retryFailedCleanJobs() {
  const result = await QueuedCleanJob.updateMany(
    {
      status: 'failed',
      attempts: { $lt: 3 },
      lastTriedAt: { $lte: new Date(Date.now() - 30 * 1000) },
    },
    {
      $set: {
        status: 'pending',
        lastTriedAt: new Date(),
      },
      $inc: { attempts: 1 },
    }
  );

  if (result.modifiedCount > 0) {
    console.log(`Retried ${result.modifiedCount} failed clean job(s)`);
  }
}
