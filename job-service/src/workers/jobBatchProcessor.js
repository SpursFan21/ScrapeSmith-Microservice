//ScrapeSmith\job-service\src\workers\jobBatchProcessor.js

import { ScheduledJob } from '../models/ScheduledJob.js';
import { QueuedScrapeJob } from '../models/QueuedScrapeJob.js';

export async function processJobBatchQueue() {
  const now = new Date();

  // Fetch jobs that are scheduled and ready to run
  const jobs = await ScheduledJob.find({
    status: 'scheduled',
    runAt: { $lte: now }
  })
    .sort({ runAt: 1 })
    .limit(3);

  if (jobs.length === 0) {
    return;
  }

  // Process the jobs
  await Promise.all(jobs.map(job => handleSingleJob(job)));
}

async function handleSingleJob(job) {
  try {
    // Lock the job and update its status to "processing"
    await ScheduledJob.findByIdAndUpdate(job._id, {
      $set: {
        status: 'processing',
        lastTriedAt: new Date()
      },
      $inc: { attempts: 1 }
    });

    // Generate a unique order ID if not present (could be UUID or similar)
    const orderId = job.orderId || job._id.toString();

    // Prepare payload for the MongoDB queue
    const mongoQueuePayload = {
      orderId: orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      status: 'pending',  // Status set to "pending" initially
      attempts: 0,
      createdAt: job.createdAt,
    };

    // Insert the job into MongoDB queue (scraping service will handle this job)
    await QueuedScrapeJob.create(mongoQueuePayload);
    console.log(`Job ${orderId} added to MongoDB queue`);

    console.log(`Job ${orderId} dispatched to MongoDB queue and remains in processing state`);
  } catch (err) {
    console.error(`Error dispatching job ${job._id}:`, err.message);

    // If job has failed 3 times, mark it as permanently_failed
    const failStatus = job.attempts >= 3 ? 'permanently_failed' : 'failed';

    await ScheduledJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message
    });
  }
}
