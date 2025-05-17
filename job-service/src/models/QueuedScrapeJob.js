// ScrapeSmith/job-service/src/models/QueuedScrapeJob.js

import mongoose from 'mongoose';

const queuedScrapeJobSchema = new mongoose.Schema({
  orderId: { type: String, required: true, unique: true },
  userId: { type: String, required: true },
  url: { type: String, required: true },
  analysisType: { type: String, required: true },
  customScript: { type: String, default: null },
  createdAt: { type: Date, required: true, default: Date.now },
  status: {
    type: String,
    enum: ['pending', 'processing', 'done', 'failed'],
    default: 'pending'
  },
  attempts: { type: Number, default: 0 },
  lastTriedAt: { type: Date, default: null }
});

export const QueuedScrapeJob = mongoose.model('QueuedScrapeJob', queuedScrapeJobSchema);
