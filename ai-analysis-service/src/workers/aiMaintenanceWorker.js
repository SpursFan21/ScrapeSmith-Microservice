//ScrapeSmith\ai-analysis-service\src\workers\aiMaintenanceWorker.js

import { QueuedAIJob } from "../models/QueuedAIJob.js";

export async function runAIMaintenance() {
  const now = new Date();

  // Delete old "done" jobs (older than 1 day)
  const deleted = await QueuedAIJob.deleteMany({
    status: "done",
    createdAt: { $lt: new Date(now.getTime() - 24 * 60 * 60 * 1000) },
  });
  if (deleted.deletedCount > 0) {
    console.log(`Deleted ${deleted.deletedCount} old done jobs`);
  }

  // Delete "permanently_failed" jobs
  const failed = await QueuedAIJob.deleteMany({ status: "permanently_failed" });
  if (failed.deletedCount > 0) {
    console.log(`Deleted ${failed.deletedCount} permanently failed jobs`);
  }

  // Recover stuck jobs older than 10 min
  const rescued = await QueuedAIJob.updateMany(
    {
      status: "processing",
      lastAttemptAt: { $lt: new Date(now.getTime() - 10 * 60 * 1000) },
    },
    { $set: { status: "failed" } }
  );
  if (rescued.modifiedCount > 0) {
    console.log(`Recovered ${rescued.modifiedCount} stuck AI jobs`);
  }
}
