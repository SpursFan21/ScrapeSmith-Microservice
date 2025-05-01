//ScrapeSmith\data-cleaning-service\src\workers\cleaningBatchProcessor.js

import { QueuedCleanJob } from '../models/QueuedCleanJob.js';
import { CleanedData } from '../models/cleanedData.js';
import { QueuedAIJob } from '../../ai-analysis-service/src/models/QueuedAIJob.js';
import * as cheerio from 'cheerio';

// Function to clean HTML content
function cleanHTMLContent(rawHtml) {
  const $ = cheerio.load(rawHtml);
  $('style, script, link, button, nav, footer, header').remove();

  let meaningfulText = '';
  $('h1, h2, h3').each((_, el) => meaningfulText += `${$(el).text().trim()}\n\n`);
  $('p').each((_, el) => meaningfulText += `${$(el).text().trim()}\n\n`);
  $('li').each((_, el) => meaningfulText += `- ${$(el).text().trim()}\n`);
  $('a').each((_, el) => {
    const href = $(el).attr('href');
    const text = $(el).text().trim();
    if (href && text) meaningfulText += `Link: ${text} (${href})\n`;
  });

  return meaningfulText;
}

export async function processCleanBatchQueue() {
  const jobs = await QueuedCleanJob.find({ status: 'pending' })
    .sort({ createdAt: 1 })
    .limit(3);

  if (jobs.length === 0) {
    return;
  }

  await Promise.all(jobs.map(job => handleCleanJob(job)));
}

async function handleCleanJob(job) {
  try {
    // Lock the job
    await QueuedCleanJob.findByIdAndUpdate(job._id, {
      $set: { status: 'processing', lastTriedAt: new Date() },
      $inc: { attempts: 1 }
    });

    // Skip if already cleaned
    const exists = await CleanedData.findOne({ orderId: job.orderId, userId: job.userId });
    if (exists) {
      await QueuedCleanJob.findByIdAndUpdate(job._id, { status: 'done' });
      return;
    }

    const cleanedContent = cleanHTMLContent(job.rawData);

    await CleanedData.create({
      orderId: job.orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript || null,
      createdAt: job.createdAt || new Date(),
      cleanedData: cleanedContent,
    });

    await QueuedCleanJob.findByIdAndUpdate(job._id, { status: 'done', error: null });
    console.log(`Cleaned job: ${job.orderId}`);

    // Add cleaned data directly to the AI queue in MongoDB
    const aiJob = {
      orderId: job.orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      createdAt: job.createdAt,
      cleanedData: cleanedContent,
    };

    // Insert the AI job into the AI service's MongoDB Atlas collection (AI queue)
    await QueuedAIJob.create(aiJob);  // Direct insertion into the AI service's queue

    console.log(`Job added to AI queue in MongoDB: ${job.orderId}`);

  } catch (err) {
    console.error(`Error in cleaning job ${job.orderId}:`, err.message);
    const failStatus = job.attempts >= 3 ? 'failed' : 'pending';
    await QueuedCleanJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message,
    });
  }
}
