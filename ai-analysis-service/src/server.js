//ScrapeSmith\ai-analysis-service\src\server.js
import express from  'express';
import dotenv from 'dotenv';
import mongoose from 'mongoose';
import analysisRoutes from './routes/analysisRoutes.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3006;

app.use(express.json());
app.use('/api/analyize', analysisRoutes);

app.get('/', (req, res) => {
    res.send('AI Analysis Service is running...');
});

mongoose.connect(process.env.MONGO_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
})
.then(() => {
    console.log('Connected to MongoDB');
    app.listen(PORT, () => console.log(`Server is running on PORT ${PORT}`));
})
.catch(err => {
    console.error(' MongoDB connection error: ', err);
});