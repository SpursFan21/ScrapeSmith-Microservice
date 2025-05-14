// admin-service/server.js

import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import adminRoutes from './routes/adminRoutes.js';
import { connectMongoose } from './utils/mongooseConnect.js';

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

// Admin routes
app.use('/admin', adminRoutes);

// Health check
app.get('/', (req, res) => {
  res.json({ message: 'Admin Service is running ' });
});

const PORT = process.env.PORT || 3005;

app.listen(PORT, async () => {
  console.log(`Admin Service running on port ${PORT}`);
  await connectMongoose();
});
