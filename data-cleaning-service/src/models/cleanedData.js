import mongoose from 'mongoose';

const cleanedDataSchema = new mongoose.Schema({
  userId: { type: String, required: true },
  orderId: { type: String, required: true },
  cleanedContent: { type: String, required: true },
});

export const CleanedData = mongoose.model('CleanedData', cleanedDataSchema);
