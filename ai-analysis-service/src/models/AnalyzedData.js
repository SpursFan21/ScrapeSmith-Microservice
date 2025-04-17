//ScrapeSmith\ai-analysis-service\src\models\AnalyzedData.js
import mongoose from "mongoose";

const analyzedDataSchema = new mongoose.Schema({
    orderId: {type: String, required: true},
    userId: {type: String, required: true},
    createdAt: {type: Date, default: Date.now},
    url: {type: String, required: true},
    analysisType: {type: String, required: true},
    customScript: {type: String},
    analysisData: {type: String, required: true},
});

export const AnalyzedData = mongoose.model('AnalyzedData', analyzedDataSchema);