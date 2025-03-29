import express from 'express';
import cleaningRoutes from './routes/cleaningRoutes.js';
import { connectMongo } from './utils/db.js';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3004;

app.use(express.json());

// Connect to MongoDB
connectMongo();

// Routes
app.use('/api/clean', cleaningRoutes);

app.listen(PORT, () => {
  console.log(`Data Cleaning Service running on port ${PORT}`);
});
