// admin-service\config\config.js
import dotenv from 'dotenv';
dotenv.config();

export const config = {
  db: {
    host: process.env.DB_HOST,
    port: process.env.DB_PORT,
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    name: process.env.DB_NAME,
    ssl: process.env.SSL_MODE === 'require',
  },
  jwtSecret: process.env.JWT_SECRET_KEY,
};
