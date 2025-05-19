//ScrapeSmith\job-service\src\handlers\scheduleHandler.js

import { ScheduledJob } from '../models/ScheduledJob.js';
import { v4 as uuidv4 } from 'uuid';

export const scheduleJob = async (req, res) => {
  try {
    const jobDataArray = req.body;

    // Ensure the request body is an array of jobs
    if (!Array.isArray(jobDataArray)) {
      return res.status(400).json({ error: 'Expected an array of job data' });
    }

    const validJobs = [];
    const invalidJobs = [];

    // Validate each job entry in the request
    for (const [index, jobData] of jobDataArray.entries()) {
      const { url, userId, analysisType, runAt, customScript } = jobData;

      // Check for required fields
      if (!url || !userId || !analysisType || !runAt) {
        invalidJobs.push({
          jobIndex: index,
          error: 'Missing required fields (url, userId, analysisType, runAt)',
        });
        continue;
      }

      // Validate userId type
      if (typeof userId !== 'string') {
        invalidJobs.push({
          jobIndex: index,
          error: 'userId must be a string',
        });
        continue;
      }

      // Validate runAt is a proper date
      const parsedRunAt = new Date(runAt);
      if (isNaN(parsedRunAt.getTime())) {
        invalidJobs.push({
          jobIndex: index,
          error: `Invalid runAt date format for job ${index}`,
        });
        continue;
      }

      // Validate optional customScript type
      if (customScript && typeof customScript !== 'string') {
        invalidJobs.push({
          jobIndex: index,
          error: 'customScript must be a string if provided',
        });
        continue;
      }

      // Convert runAt to UTC date and add to valid jobs
      const utcRunAt = new Date(parsedRunAt.toISOString());

      validJobs.push({
        orderId: uuidv4(), // Generate unique orderId for tracking the job
        url,
        userId,
        analysisType,
        customScript,
        runAt: utcRunAt,
      });
    }

    // If any jobs were invalid, return early with detailed error info
    if (invalidJobs.length > 0) {
      return res.status(400).json({
        error: 'Some jobs have validation errors',
        details: invalidJobs,
      });
    }

    // Insert valid jobs into MongoDB using ScheduledJob model
    const newJobs = await ScheduledJob.insertMany(validJobs);

    // Prepare a user-friendly response with both UTC and local time
    const responseJobs = newJobs.map(job => ({
      id: job._id,
      orderId: job.orderId,
      url: job.url,
      analysisType: job.analysisType,
      runAtUTC: job.runAt,
      runAtLocal: job.runAt.toLocaleString(),
    }));

    // Respond with success and job details
    res.status(201).json({
      message: `${newJobs.length} job(s) scheduled successfully`,
      jobs: responseJobs,
    });
  } catch (err) {
    console.error('Failed to schedule jobs:', err.message);
    res.status(500).json({ error: 'Failed to schedule jobs' });
  }
};
