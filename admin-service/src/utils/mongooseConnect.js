// src/utils/mongooseConnect.js
import mongoose from 'mongoose';
import dotenv from 'dotenv';
dotenv.config();

export const connectMongoose = async () => {
  if (mongoose.connection.readyState >= 1) return;

  try {
    await mongoose.connect(process.env.MONGO_URI, {
      dbName: process.env.MONGO_DB,
    });
    console.log("✅ Mongoose connected in admin-service");
  } catch (error) {
    console.error("❌ Mongoose connection error:", error);
    throw error;
  }
};
