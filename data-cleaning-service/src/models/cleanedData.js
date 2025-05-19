// data-cleaning-service/src/models/cleanedData.js

import mongoose from 'mongoose';

const cleanedDataSchema = new mongoose.Schema({
  orderId:      { type: String, required: true },
  userId:       { type: String, required: true },
  createdAt:    { type: Date,   required: true, default: Date.now },
  url:          { type: String, required: true },
  analysisType: { type: String, required: true },
  customScript: { type: String },
  cleanedData:  { type: String, required: true }
}, {
  collection: 'cleaned_data'
});

export const CleanedData = mongoose.model('CleanedData', cleanedDataSchema);
