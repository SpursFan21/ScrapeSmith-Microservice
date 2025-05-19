//ScrapeSmith\ai-analysis-service\src\models\QueuedAIJob.js

import mongoose from "mongoose";

const queuedAIJobSchema = new mongoose.Schema({
  orderId:       { type: String, required: true, unique: true },
  userId:        { type: String, required: true },
  createdAt:     { type: Date, default: Date.now },
  url:           { type: String, required: true },
  analysisType:  { type: String, required: true },
  customScript:  { type: String },
  cleanedData:   { type: String, required: true },
  status:        { type: String, enum: ["pending", "processing", "done", "failed"], default: "pending" },
  attempts:      { type: Number, default: 0 },
  lastAttemptAt: { type: Date },
  error:         { type: String },
}, {
  collection: 'queued_ai_jobs'
});

export const QueuedAIJob = mongoose.model("QueuedAIJob", queuedAIJobSchema);
