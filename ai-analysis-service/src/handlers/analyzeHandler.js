//ScrapeSmith\ai-analysis-service\src\handlers\analyzeHandler.js
import { buildPrompt } from "../utils/promptBuilder.js";
import { analyzeText } from "../utils/openaiClient.js";
import { AnalyzedData } from "../models/AnalyzedData.js";
import { analysisTypePrompts } from "../utils/analysisPrompts.js";

const basePrompt = `You are an expert data analyst. Analyze the following cleaned web data based on the prompt instructions provided.`;

export const handleAnalysis = async (req, res) => {
  try {
    const {
      orderId,
      userId,
      createdAt,
      url,
      analysisType,
      customScript,
      cleanedData,
    } = req.body;

    const analysisTypePrompt =
      analysisTypePrompts[analysisType] || "Provide a general analysis of the data.";

    const finalPrompt = buildPrompt({
      basePrompt,
      analysisTypePrompt,
      customScript,
      cleanedData,
    });

    // Log the final prompt for debugging
    console.log(`[Prompt Used] (${analysisType})\n${finalPrompt}`);

    const result = await analyzeText(finalPrompt);

    const saved = await AnalyzedData.create({
      orderId,
      userId,
      createdAt,
      url,
      analysisType,
      customScript,
      analysisData: result,
    });

    res.status(201).json(saved);
  } catch (err) {
    console.error("Analysis Error:", err);
    res.status(500).json({ error: "Failed to analyze data" });
  }
};
