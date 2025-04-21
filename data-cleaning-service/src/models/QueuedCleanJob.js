//ScrapeSmith\data-cleaning-service\src\models\QueuedCleanJob.js
import mongoose from 'mongoose';

const queuedCleanJobSchema = new mongoose.Schema({
  orderId: { type: String, required: true, unique: true },
  userId: { type: String, required: true },
  url: { type: String, required: true },
  analysisType: { type: String, required: true },
  customScript: { type: String },
  createdAt: { type: Date, default: Date.now },
  rawData: { type: String, required: true },
  status: { type: String, enum: ['pending', 'processing', 'done', 'failed'], default: 'pending' },
  attempts: { type: Number, default: 0 },
  lastTriedAt: { type: Date },
  error: { type: String },
});

export const QueuedCleanJob = mongoose.model('QueuedCleanJob', queuedCleanJobSchema);
