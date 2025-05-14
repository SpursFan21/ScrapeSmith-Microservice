
// admin-service/models/User.js
import mongoose from 'mongoose';

const userSchema = new mongoose.Schema(
  {
    email: { type: String, required: true },
    username: { type: String, required: true },
    hashed_password: { type: String, required: true, select: false }, // excluded by default
    is_admin: { type: Boolean, default: false },
    created_at: { type: Date, default: Date.now },
  },
  { collection: 'users' }
);

export const User = mongoose.model('User', userSchema);
