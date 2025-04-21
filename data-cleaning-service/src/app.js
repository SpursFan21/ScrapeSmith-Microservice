// data-cleaning-service/src/app.js
import express from 'express';
import cleaningRoutes from './routes/cleaningRoutes.js';
import { connectMongo } from './utils/db.js';
import dotenv from 'dotenv';
import queueRoutes from './routes/queueRoutes.js';
import { processCleanQueue } from './workers/cleaningProcessor.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3004;

app.use(express.json({ limit: '10mb' }));

connectMongo();

app.use('/api/clean', cleaningRoutes);
app.use('/api/clean/queue', queueRoutes);

setInterval(() => {
  processCleanQueue().catch(err => console.error('Queue error:', err));
}, 5000);

app.listen(PORT, () => {
  console.log(`Data Cleaning Service running on port ${PORT}`);
});
