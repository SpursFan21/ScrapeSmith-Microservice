//ScrapeSmith\ai-analysis-service\src\utils\db.js

// ScrapeSmith/ai-analysis-service/src/utils/db.js
import mongoose from 'mongoose';
import dotenv from 'dotenv';

dotenv.config();

export const connectMongo = async () => {
  try {
    await mongoose.connect(process.env.MONGO_URI, {
      useNewUrlParser: true,
      useUnifiedTopology: true,
      dbName: process.env.MONGO_DB,
    });

    console.log(`✅ Connected to MongoDB database: ${process.env.MONGO_DB}`);
  } catch (err) {
    console.error('❌ MongoDB connection error:', err.message);
    process.exit(1);
  }
};
