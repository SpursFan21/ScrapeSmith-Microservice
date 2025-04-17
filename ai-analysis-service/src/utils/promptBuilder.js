//ScrapeSmith\ai-analysis-service\src\utils\promptBuilder.js
export function buildPrompt({ basePrompt, analysisTypePrompt, customScript, cleanedData }) {
    return `
  ${basePrompt}
  
  ${analysisTypePrompt}
  
  ${customScript ? `Custom Script Instructions:\n${customScript}` : ''}
  
  Here is the cleaned data to analyze:
  ${cleanedData}
  `;
  }
  