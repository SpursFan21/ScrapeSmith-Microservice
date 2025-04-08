// data-cleaning-service/src/utils/db.js
import mongoose from 'mongoose';

export const connectMongo = async () => {
  try {
    // use URI as-is from .env — it's complete already
    await mongoose.connect(process.env.MONGO_URI, {
      useNewUrlParser: true,
      useUnifiedTopology: true,
      dbName: process.env.MONGO_DB,
    });

    console.log(`✅ Connected to MongoDB database: ${process.env.MONGO_DB}`);
  } catch (err) {
    console.error('❌ MongoDB connection error:', err);
    process.exit(1);
  }
};
