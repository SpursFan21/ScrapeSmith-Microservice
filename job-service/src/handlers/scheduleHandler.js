//ScrapeSmith\job-service\src\handlers\scheduleHandler.js

import { ScheduledJob } from '../models/ScheduledJob.js';

export const scheduleJob = async (req, res) => {
  try {
    const jobDataArray = req.body; // Expecting an array of jobs

    // Step 1: Ensure jobDataArray is an array
    if (!Array.isArray(jobDataArray)) {
      return res.status(400).json({ error: 'Expected an array of job data' });
    }

    // Step 2: Validate each job object
    const validJobs = [];
    const invalidJobs = [];

    for (const [index, jobData] of jobDataArray.entries()) {
      const { url, userId, analysisType, runAt, customScript } = jobData;

      if (!url || !userId || !analysisType || !runAt) {
        invalidJobs.push({
          jobIndex: index,
          error: 'Missing required fields (url, userId, analysisType, runAt)',
        });
        continue;
      }

      if (typeof userId !== 'string') {
        invalidJobs.push({
          jobIndex: index,
          error: 'userId must be a string',
        });
        continue;
      }

      const parsedRunAt = new Date(runAt);
      if (isNaN(parsedRunAt.getTime())) {
        invalidJobs.push({
          jobIndex: index,
          error: `Invalid runAt date format for job ${index}`,
        });
        continue;
      }

      if (customScript && typeof customScript !== 'string') {
        invalidJobs.push({
          jobIndex: index,
          error: 'customScript must be a string if provided',
        });
        continue;
      }

      // Force runAt to UTC ISO format
      const utcRunAt = new Date(parsedRunAt.toISOString());

      validJobs.push({
        url,
        userId,
        analysisType,
        customScript,
        runAt: utcRunAt,
      });
    }

    if (invalidJobs.length > 0) {
      return res.status(400).json({
        error: 'Some jobs have validation errors',
        details: invalidJobs,
      });
    }

    // Step 3: Insert all valid jobs into the ScheduledJob collection
    const newJobs = await ScheduledJob.insertMany(validJobs);

    // Step 4: Return response with job details
    res.status(201).json({
      message: `${newJobs.length} job(s) scheduled successfully`,
      jobIds: newJobs.map((job) => job._id),
    });
  } catch (err) {
    console.error('Failed to schedule jobs:', err.message);
    res.status(500).json({ error: 'Failed to schedule jobs' });
  }
};
