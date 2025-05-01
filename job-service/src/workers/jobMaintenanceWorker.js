//ScrapeSmith\job-service\src\workers\jobMaintenanceWorker.js

import { ScheduledJob } from '../models/ScheduledJob.js';

export async function runJobQueueMaintenance() {
  const now = new Date();

  // 1. Clean "done" jobs older than 24h
  const cleaned = await ScheduledJob.deleteMany({
    status: 'done',
    createdAt: { $lt: new Date(now.getTime() - 24 * 60 * 60 * 1000) }
  });
  if (cleaned.deletedCount > 0) {
    console.log(`Deleted ${cleaned.deletedCount} old done jobs`);
  }

  // 2. Clean permanently failed jobs
  const failed = await ScheduledJob.deleteMany({
    status: 'permanently_failed'
  });
  if (failed.deletedCount > 0) {
    console.log(`Deleted ${failed.deletedCount} permanently failed jobs`);
  }

  // 3. Recover stuck jobs older than 10 mins
  const recovered = await ScheduledJob.updateMany(
    {
      status: 'processing',
      lastTriedAt: { $lt: new Date(now.getTime() - 10 * 60 * 1000) }
    },
    { $set: { status: 'failed' } }
  );
  if (recovered.modifiedCount > 0) {
    console.log(`Recovered ${recovered.modifiedCount} stuck jobs and marked as failed`);
  }
}

