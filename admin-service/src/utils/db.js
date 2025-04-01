// utils/db.js
import { Pool } from 'pg';
import { config } from '../config/config.js';

const pool = new Pool({
  host: config.db.host,
  port: config.db.port,
  user: config.db.user,
  password: config.db.password,
  database: config.db.name,
  ssl: config.db.ssl,
});

pool.on('connect', () => {
  console.log('✅ admin-service connected to PostgreSQL');
});

export default pool;
