//data-cleaning-service\src\handlers\cleaningHandler.js
import * as cheerio from 'cheerio';
import { CleanedData } from '../models/cleanedData.js';

function cleanHTMLContent(rawHtml) {
  const $ = cheerio.load(rawHtml);
  $('style, script, link, button, nav, footer, header').remove();

  let meaningfulText = '';
  $('h1, h2, h3').each((_, el) => meaningfulText += `${$(el).text().trim()}\n\n`);
  $('p').each((_, el) => meaningfulText += `${$(el).text().trim()}\n\n`);
  $('li').each((_, el) => meaningfulText += `- ${$(el).text().trim()}\n`);
  $('a').each((_, el) => {
    const href = $(el).attr('href');
    const text = $(el).text().trim();
    if (href && text) meaningfulText += `Link: ${text} (${href})\n`;
  });

  return meaningfulText;
}

export const cleanAndStoreData = async (req, res) => {
  const {
    orderId,
    userId,
    url,
    analysisType,
    customScript,
    createdAt,
    rawData
  } = req.body;

  try {
    if (!userId || !orderId || !url || !analysisType || typeof rawData !== 'string') {
      console.warn("❌ Missing or invalid required fields in request body");
      return res.status(400).json({ error: 'Missing or invalid required fields' });
    }

    const exists = await CleanedData.findOne({ orderId, userId });
    if (exists) {
      return res.status(200).json({
        message: 'Cleaned data already exists',
        cleanedOrderId: exists._id,
        cleanedData: exists.cleanedData,
      });
    }

    const cleanedContent = cleanHTMLContent(rawData);

    const cleanedDataEntry = await CleanedData.create({
      orderId,
      userId,
      url,
      analysisType,
      customScript: customScript || null,
      createdAt: createdAt ? new Date(createdAt) : new Date(),
      cleanedData: cleanedContent,
    });

    res.status(200).json({
      message: 'Data cleaned and stored successfully',
      cleanedOrderId: cleanedDataEntry._id,
      cleanedData: cleanedDataEntry.cleanedData,
    });
  } catch (err) {
    console.error('❌ Error during data cleaning:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
};
