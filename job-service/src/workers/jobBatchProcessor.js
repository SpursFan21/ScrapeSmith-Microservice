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

    const orderId = job.orderId || job._id.toString();

    const mongoQueuePayload = {
      orderId: orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      status: 'pending',
      attempts: 0,
      createdAt: job.createdAt,
    };

    await QueuedScrapeJob.create(mongoQueuePayload);

    console.log(`Job ${orderId} added to [scraping-service MongoDB queue]`);
    console.log(`Job ${orderId} dispatched to [scraping queue] and marked as processing in scheduled jobs`);
  } catch (err) {
    console.error(`Error dispatching job ${job._id}:`, err.message);

    const failStatus = job.attempts >= 3 ? 'permanently_failed' : 'failed';

    await ScheduledJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message
    });
  }
}

