//ScrapeSmith\ai-analysis-service\src\utils\analysisPrompts.js
export const analysisTypePrompts = {
    "Sentiment Analysis":
      "Analyze the sentiment of the text and return whether it is positive, negative, or neutral. Provide a short explanation for your conclusion.",
      
    "Keyword Extraction":
      "Extract the most important keywords and key phrases from the text. Present them as a bullet list ranked by relevance.",
      
    "Entity Recognition":
      "Identify all named entities in the text, such as people, organizations, locations, and dates. Present them in a categorized format.",
      
    "Text Summarization":
      "Summarize the content into a concise paragraph. Focus on the main ideas and essential details.",
      
    "Topic Modeling":
      "Identify and list the main topics discussed in the text. Group related ideas under each topic and provide a brief description.",
      
    "Named Entity Linking":
      "Find all entities in the text and attempt to link each one to a real-world entry such as Wikipedia, Wikidata, or a known public figure.",
      
    "Language Detection":
      "Determine the primary language used in the text. If the text includes multiple languages, list all detected and their confidence scores.",
      
    "Content Classification":
      "Classify the text into appropriate categories (e.g., news, blog, scientific, opinion). Justify the classification with a short reasoning.",
      
    "Anomaly Detection":
      "Identify any unusual or unexpected information or language in the text that stands out from the rest. Explain why it appears anomalous.",
      
    "Text Clustering":
      "Divide the text into logical segments or clusters based on themes or topics. Describe each cluster in 1–2 sentences.",
      
    "Custom Analysis":
      "Follow the custom instructions provided by the user to analyze the text accordingly. Perform no additional analysis beyond the user-defined scope.",
  };
  