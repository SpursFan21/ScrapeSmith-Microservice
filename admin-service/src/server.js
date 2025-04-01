//admin-service\server.js
import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import adminRoutes from './routes/adminRoutes.js';
import { connectMongo } from './utils/mongoClient.js';

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

app.use('/admin', adminRoutes);

app.get('/', (req, res) => {
  res.json({ message: 'Admin Service is running 🛠️' });
});

const PORT = process.env.PORT || 3005;
app.listen(PORT, () => {
  console.log(`🚀 Admin Service running on port ${PORT}`);
});

connectMongo();
