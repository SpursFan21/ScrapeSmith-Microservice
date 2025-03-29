import * as cheerio from 'cheerio';
import { CleanedData } from '../models/cleanedData.js';

// Cleaning function replicating the Python logic
function cleanHTMLContent(rawHtml) {
  const $ = cheerio.load(rawHtml);

  // Remove unwanted tags
  $('style, script, link, button, nav, footer, header').remove();

  let meaningfulText = '';

  // Collect header tags (h1, h2, h3)
  $('h1, h2, h3').each((_, element) => {
    meaningfulText += `${$(element).text().trim()}\n\n`;
  });

  // Collect paragraph tags
  $('p').each((_, element) => {
    meaningfulText += `${$(element).text().trim()}\n\n`;
  });

  // Collect list items (li)
  $('li').each((_, element) => {
    meaningfulText += `- ${$(element).text().trim()}\n`;
  });

  // Add <a> links with valuable anchor text and URLs
  $('a').each((_, element) => {
    const href = $(element).attr('href');
    const text = $(element).text().trim();
    if (href && text) {
      meaningfulText += `Link: ${text} (${href})\n`;
    }
  });

  return meaningfulText;
}

// Handler function to clean and store data
export const cleanAndStoreData = async (req, res) => {
  const { orderId, userId, rawData } = req.body;

  try {
    // Check if data already exists for the order and user
    const existingEntry = await CleanedData.findOne({ orderId, userId });
    if (existingEntry) {
      return res.status(200).json({
        message: 'Cleaned data already exists',
        cleanedOrderId: existingEntry._id,
        cleanedData: existingEntry.cleanedContent,
      });
    }

    // Clean the raw HTML content
    const cleanedContent = cleanHTMLContent(rawData);

    // Store cleaned data in MongoDB
    const cleanedDataEntry = await CleanedData.create({
      userId,
      orderId,
      cleanedContent,
    });

    res.status(200).json({
      message: 'Data cleaned and stored successfully',
      cleanedOrderId: cleanedDataEntry._id,
      cleanedData: cleanedDataEntry.cleanedContent,
    });
  } catch (err) {
    console.error('Error during data cleaning:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
};
