//ScrapeSmith\job-service\src\handlers\scheduleHandler.js

import { ScheduledJob } from '../models/ScheduledJob.js';

export const scheduleJob = async (req, res) => {
  try {
    const jobDataArray = req.body;  // Expecting an array of jobs

    // Step 1: Ensure jobDataArray is an array
    if (!Array.isArray(jobDataArray)) {
      return res.status(400).json({ error: 'Expected an array of job data' });
    }

    // Step 2: Validate each job object
    const invalidJobs = [];
    for (const [index, jobData] of jobDataArray.entries()) {
      if (!jobData.url || !jobData.userId || !jobData.analysisType || !jobData.runAt) {
        invalidJobs.push({
          jobIndex: index,
          error: 'Missing required fields (url, userId, analysisType, runAt)'
        });
      }

      // Validate runAt date format
      const runAt = new Date(jobData.runAt);
      if (isNaN(runAt)) {
        invalidJobs.push({
          jobIndex: index,
          error: `Invalid runAt date format for job ${index}`
        });
      }

      // Validate customScript (optional, must be a string if provided)
      if (jobData.customScript && typeof jobData.customScript !== 'string') {
        invalidJobs.push({
          jobIndex: index,
          error: 'customScript must be a string if provided'
        });
      }
    }

    if (invalidJobs.length > 0) {
      return res.status(400).json({
        error: 'Some jobs have validation errors',
        details: invalidJobs,
      });
    }

    // Step 3: Insert all valid jobs into the ScheduledJob collection
    const newJobs = await ScheduledJob.insertMany(jobDataArray);

    // Step 4: Return response with job details
    res.status(201).json({
      message: `${newJobs.length} job(s) scheduled successfully`,
      jobIds: newJobs.map(job => job._id),
    });
  } catch (err) {
    console.error('Failed to schedule jobs:', err.message);
    res.status(500).json({ error: 'Failed to schedule jobs' });
  }
};
