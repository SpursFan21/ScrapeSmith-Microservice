//ScrapeSmith\admin-service\src\models\Scrape.js

import mongoose from 'mongoose';

const scrapeSchema = new mongoose.Schema({
  order_id:      { type: String, required: true },
  user_id:       { type: String, required: true },
  created_at:    { type: Date,   default: Date.now },
  url:           { type: String, required: true },
  analysis_type: { type: String, required: true },
  custom_script: { type: String },
  data:          { type: String, required: true },
}, { collection: 'scraped_data' });

export const Scrape = mongoose.model('Scrape', scrapeSchema);
