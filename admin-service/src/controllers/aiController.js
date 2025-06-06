//ScrapeSmith\admin-service\src\controllers\aiController.js

import { AnalyzedData } from '../models/AnalyzedData.js';

export const getAllAIResults = async (req, res) => {
  try {
    const results = await AnalyzedData.find().sort({ createdAt: -1 });
    res.json(results);
  } catch (error) {
    console.error("Failed to fetch AI analysis results:", error);
    res.status(500).json({ error: "Failed to fetch AI analysis results" });
  }
};
