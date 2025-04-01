//admin-service\models\ticketModel.js
import mongoose from 'mongoose';

const responseSchema = new mongoose.Schema({
  adminId: { type: String },
  message: { type: String, required: true },
  timestamp: { type: Date, default: Date.now },
});

const ticketSchema = new mongoose.Schema({
  userId: { type: String, required: true },
  subject: { type: String, required: true },
  message: { type: String, required: true },
  status: { type: String, enum: ['open', 'closed'], default: 'open' },
  responses: [responseSchema],
  createdAt: { type: Date, default: Date.now },
});

export const Ticket = mongoose.model('Ticket', ticketSchema);

