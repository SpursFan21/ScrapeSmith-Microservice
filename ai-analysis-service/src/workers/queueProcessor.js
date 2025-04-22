//ScrapeSmith\ai-analysis-service\src\workers\queueProcessor.js

import { QueuedAIJob } from "../models/QueuedAIJob.js";
import { AnalyzedData } from "../models/AnalyzedData.js";
import { analysisTypePrompts } from "../utils/analysisPrompts.js";
import { buildPrompt } from "../utils/promptBuilder.js";
import { analyzeText } from "../utils/openaiClient.js";

const basePrompt = `You are an expert data analyst. Analyze the following cleaned web data based on the prompt instructions provided.`;

export async function processQueue() {
  const job = await QueuedAIJob.findOneAndUpdate(
    { status: "pending" },
    { status: "processing", lastAttemptAt: new Date(), $inc: { attempts: 1 } },
    { sort: { createdAt: 1 }, new: true }
  );

  if (!job) return; // No job to process

  try {
    const analysisTypePrompt = analysisTypePrompts[job.analysisType] || "Provide a general analysis of the data.";
    const prompt = buildPrompt({
      basePrompt,
      analysisTypePrompt,
      customScript: job.customScript,
      cleanedData: job.cleanedData,
    });

    const analysis = await analyzeText(prompt);

    await AnalyzedData.create({
      orderId: job.orderId,
      userId: job.userId,
      createdAt: job.createdAt,
      url: job.url,
      analysisType: job.analysisType,
      customScript: job.customScript,
      analysisData: analysis,
    });

    await QueuedAIJob.findByIdAndUpdate(job._id, { status: "done", error: null });

    console.log(`[AI Analysis Complete] Job ${job.orderId}`);
  } catch (err) {
    console.error(`[AI Analysis Failed] Job ${job.orderId}`, err.message);
    await QueuedAIJob.findByIdAndUpdate(job._id, {
      status: job.attempts >= 3 ? "failed" : "pending",
      error: err.message,
    });
  }
}

