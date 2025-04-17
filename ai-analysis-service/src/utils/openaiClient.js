//ScrapeSmith\ai-analysis-service\src\utils\openaiClient.js
import { Configuration, OpenAIApi } from "openai";

const config = new Configuration({
  apiKey: process.env.OPENAI_API_KEY,
});

const openai = new OpenAIApi(config);

export async function analyzeText(prompt) {
  const response = await openai.createChatCompletion({
    model: "gpt-4o",
    messages: [{ role: "user", content: prompt }],
  });

  const usage = response.data.usage;
  console.log(`[OpenAI Usage] Input Tokens: ${usage.prompt_tokens}, Output Tokens: ${usage.completion_tokens}, Total: ${usage.total_tokens}`);

  return response.data.choices[0].message.content;
}
