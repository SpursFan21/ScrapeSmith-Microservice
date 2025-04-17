//ScrapeSmith\ai-analysis-service\src\routes\analysisRoutes.js
import express from "express";
import { handleAnalysis } from "../handlers/analyzeHandler.js";

const router = express.Router();

router.post("/", handleAnalysis);

export default router;
