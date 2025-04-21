//ScrapeSmith\data-cleaning-service\src\workers\cleaningMaintenanceWorker.js

import { QueuedCleanJob } from '../models/QueuedCleanJob.js';

export async function runCleanQueueMaintenance() {
  const now = new Date();

  // 1. Remove "done" jobs older than 1 day
  const deleted = await QueuedCleanJob.deleteMany({
    status: 'done',
    createdAt: { $lt: new Date(now.getTime() - 24 * 60 * 60 * 1000) },
  });
  if (deleted.deletedCount > 0) {
    console.log(`Cleaned ${deleted.deletedCount} old done jobs`);
  }

  // 2. Remove permanently failed jobs
  const failed = await QueuedCleanJob.deleteMany({ status: 'permanently_failed' });
  if (failed.deletedCount > 0) {
    console.log(`Deleted ${failed.deletedCount} permanently failed jobs`);
  }

  // 3. Recover stuck "processing" jobs
  const stuck = await QueuedCleanJob.updateMany(
    {
      status: 'processing',
      lastTriedAt: { $lt: new Date(now.getTime() - 10 * 60 * 1000) },
    },
    { $set: { status: 'failed' } }
  );
  if (stuck.modifiedCount > 0) {
    console.log(`Recovered ${stuck.modifiedCount} stuck jobs`);
  }
}
