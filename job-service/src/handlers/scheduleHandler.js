//ScrapeSmith\job-service\src\handlers\scheduleHandler.js

import { ScheduledJob } from '../models/ScheduledJob.js';

export const scheduleJob = async (req, res) => {
  try {
    const jobData = req.body;

    if (!jobData.url || !jobData.userId || !jobData.analysisType || !jobData.runAt) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const newJob = await ScheduledJob.create(jobData);

    res.status(201).json({
      message: 'Job scheduled successfully',
      jobId: newJob._id,
    });
  } catch (err) {
    console.error('Failed to schedule job:', err.message);
    res.status(500).json({ error: 'Failed to schedule job' });
  }
};
