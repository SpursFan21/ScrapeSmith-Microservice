//ScrapeSmith\job-service\src\models\QueuedScrapeJob.js


const mongoose = require('mongoose');

// Define the schema for QueuedScrapeJob
const queuedScrapeJobSchema = new mongoose.Schema({
  orderId: {
    type: String,
    required: true,
    unique: true
  },
  userId: {
    type: String,
    required: true
  },
  createdAt: {
    type: Date,
    required: true,
    default: Date.now
  },
  url: {
    type: String,
    required: true
  },
  analysisType: {
    type: String,
    required: true
  },
  customScript: {
    type: String,
    default: null
  },
  status: {
    type: String,
    enum: ['pending', 'processing', 'done', 'failed'],
    default: 'pending'
  },
  attempts: {
    type: Number,
    default: 0
  },
  lastTriedAt: {
    type: Date,
    default: null
  }
});

// Create and export the model based on the schema
const QueuedScrapeJob = mongoose.model('QueuedScrapeJob', queuedScrapeJobSchema);

module.exports = QueuedScrapeJob;
