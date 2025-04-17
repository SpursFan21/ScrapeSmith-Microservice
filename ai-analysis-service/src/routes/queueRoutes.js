import express from "express";
import { QueuedAIJob } from "../models/QueuedAIJob.js";

const router = express.Router();

router.post("/", async (req, res) => {
  try {
    const job = await QueuedAIJob.create(req.body);
    res.status(201).json({ message: "Job queued successfully", job });
  } catch (err) {
    console.error("Failed to queue job:", err);
    res.status(500).json({ error: "Failed to queue job" });
  }
});

export default router;
