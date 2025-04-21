//ScrapeSmith\ai-analysis-service\src\workers\aiRetryWorker.js

import { QueuedAIJob } from "../models/QueuedAIJob.js";

export async function retryFailedAIJobs() {
  const result = await QueuedAIJob.updateMany(
    {
      status: "failed",
      attempts: { $lt: 3 },
      lastAttemptAt: { $lte: new Date(Date.now() - 30 * 1000) },
    },
    {
      $set: {
        status: "pending",
        lastAttemptAt: new Date(),
      },
      $inc: { attempts: 1 },
    }
  );

  if (result.modifiedCount > 0) {
    console.log(`🔁 Retried ${result.modifiedCount} failed AI job(s)`);
  }
}
