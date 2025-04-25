//ScrapeSmith\ai-analysis-service\src\workers\aiBatchProcessor.js

import { QueuedAIJob } from "../models/QueuedAIJob.js";
import { AnalyzedData } from "../models/AnalyzedData.js";
import { analysisTypePrompts } from "../utils/analysisPrompts.js";
import { buildPrompt } from "../utils/promptBuilder.js";
import { analyzeText } from "../utils/openaiClient.js";

const basePrompt = `You are an expert data analyst. Analyze the following cleaned web data based on the prompt instructions provided.`;

export async function processAIBatchQueue() {
  const jobs = await QueuedAIJob.find({ status: "pending" })
    .sort({ createdAt: 1 })
    .limit(3);

  if (jobs.length === 0) {
    return;
  }

  await Promise.all(jobs.map(job => handleAIJob(job)));
}

async function handleAIJob(job) {
  try {
    await QueuedAIJob.findByIdAndUpdate(job._id, {
      $set: { status: "processing", lastAttemptAt: new Date() },
      $inc: { attempts: 1 },
    });

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
    console.log(`AI Analysis Complete: ${job.orderId}`);
  } catch (err) {
    console.error(`AI Job ${job.orderId} Failed:`, err.message);
    const failStatus = job.attempts >= 3 ? "failed" : "pending";
    await QueuedAIJob.findByIdAndUpdate(job._id, {
      status: failStatus,
      error: err.message,
    });
  }
}
