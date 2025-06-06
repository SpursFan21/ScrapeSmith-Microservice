//ScrapeSmith\admin-service\src\models\ScrapeData.js

import mongoose from 'mongoose';

const scrapeSchema = new mongoose.Schema({
  orderId:      { type: String, required: true },
  userId:       { type: String, required: true },
  createdAt:    { type: Date,   default: Date.now },
  url:          { type: String, required: true },
  analysisType: { type: String, required: true },
  customScript: { type: String },
  data:         { type: String, required: true },
}, { collection: 'scraped_data' });


export const Scrape = mongoose.model('Scrape', scrapeSchema);
