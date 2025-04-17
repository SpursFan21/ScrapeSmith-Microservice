//ScrapeSmith\ai-analysis-service\src\server.js
import express from 'express';
import dotenv from 'dotenv';
import mongoose from 'mongoose';
import analysisRoutes from './routes/analysisRoutes.js';
import queueRoutes from './routes/queueRoutes.js';
import { processQueue } from './workers/queueProcessor.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3006;

// Middleware
app.use(express.json());

// Routes
app.use('/api/analyize', analysisRoutes);
app.use('/api/queue', queueRoutes);

// Health check
app.get('/', (req, res) => {
  res.send('AI Analysis Service is running...');
});

// DB Connection & Queue Poller
mongoose
  .connect(process.env.MONGO_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  })
  .then(() => {
    console.log('Connected to MongoDB');

    // Start Express
    app.listen(PORT, () => {
      console.log(`AI Analysis Service running on port ${PORT}`);
    });

    // Start polling queue
    setInterval(() => {
      processQueue().catch((err) =>
        console.error('Queue processing error:', err.message)
      );
    }, 5000);
  })
  .catch((err) => {
    console.error('MongoDB connection error:', err.message);
  });
