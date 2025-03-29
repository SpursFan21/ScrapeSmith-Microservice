import express from 'express';
import { cleanAndStoreData } from '../handlers/cleaningHandler.js';

const router = express.Router();

router.post('/', cleanAndStoreData);

export default router;
