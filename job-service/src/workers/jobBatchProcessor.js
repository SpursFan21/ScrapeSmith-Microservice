//ScrapeSmith\job-service\src\workers\jobBatchProcessor.js

// ScrapeSmith/job-service/src/workers/jobBatchProcessor.js

import { ScheduledJob } from '../models/ScheduledJob.js';
import fetch from 'node-fetch';

export async function processJobBatchQueue() {
  const now = new Date();

  const jobs = await ScheduledJob.find({
    status: 'scheduled',
    runAt: { $lte: now }
  })
    .sort({ runAt: 1 })
    .limit(3);

  if (jobs.length === 0) {
    console.log("No scheduled jobs ready for dispatch");
    return;
  }

  await Promise.all(jobs.map(job => handleSingleJob(job)));
}

async function handleSingleJob(job) {
  try {
    await ScheduledJob.findByIdAndUpdate(job._id, {
      $set: {
        status: 'processing',
        lastTriedAt: new Date()
      },
      $inc: { attempts: 1 }
    });

    const payload = {
      orderId: job.orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      createdAt: job.createdAt
    };

    const scrapeUrl = process.env.SCRAPING_QUEUE_URL || 'http://scraping-service:3003/scrape/queue';

    const response = await fetch(scrapeUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      throw new Error(`Failed to queue job for scraping: ${response.statusText}`);
    }

    await ScheduledJob.findByIdAndUpdate(job._id, {
      status: 'done',
      error: null
    });

    console.log(`Dispatched job ${job.orderId} to scraping queue`);
  } catch (err) {
    console.error(`Error dispatching job ${job.orderId}:`, err.message);

    const failStatus = job.attempts >= 3 ? 'permanently_failed' : 'failed';

    await ScheduledJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message
    });
  }
}
