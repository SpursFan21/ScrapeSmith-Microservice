//ScrapeSmith\data-cleaning-service\src\workers\cleaningProcessor.js

import { QueuedCleanJob } from '../models/QueuedCleanJob.js';
import { CleanedData } from '../models/cleanedData.js';
import * as cheerio from 'cheerio';

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

export async function processCleanQueue() {
  const job = await QueuedCleanJob.findOneAndUpdate(
    { status: 'pending' },
    { $set: { status: 'processing', lastTriedAt: new Date() }, $inc: { attempts: 1 } },
    { sort: { createdAt: 1 }, new: true }
  );

  if (!job) return;

  try {
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

    // Forward to AI Analysis Service
    const aiQueueURL = process.env.AI_QUEUE_URL || "http://ai-analysis-service:3006/api/queue";

    const aiPayload = {
      orderId: job.orderId,
      userId: job.userId,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      createdAt: job.createdAt,
      cleanedData: cleanedContent,
    };

    const resp = await fetch(aiQueueURL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(aiPayload),
    });

    if (!resp.ok) {
      console.error(`Failed to queue AI analysis: ${resp.statusText}`);
    } else {
      console.log(`Cleaned job forwarded to AI queue: ${job.orderId}`);
    }

  } catch (err) {
    console.error(`Cleaning error for ${job.orderId}:`, err);
    const failStatus = job.attempts >= 3 ? 'failed' : 'pending';
    await QueuedCleanJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message,
    });
  }
}
