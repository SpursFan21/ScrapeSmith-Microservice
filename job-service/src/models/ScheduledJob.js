//ScrapeSmith\job-service\src\models\ScheduledJob.js

import mongoose from "mongoose";

const scheduledJobSchema = new mongoose.Schema({
  orderId: {
    type: String,
    required: true,
    unique: true
  },
  userId: { type: String, required: true },
  url: { type: String, required: true },
  analysisType: { type: String, required: true },
  customScript: { type: String },
  runAt: { type: Date, required: true },
  status: {
    type: String,
    enum: ["scheduled", "processing", "done", "failed"],
    default: "scheduled"
  },
  attempts: { type: Number, default: 0 },
  lastTriedAt: { type: Date },
  error: { type: String }
});

export const ScheduledJob = mongoose.model("ScheduledJob", scheduledJobSchema);
