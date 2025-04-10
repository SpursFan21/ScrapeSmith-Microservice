
# ScrapeSmith Microservices

This repository contains all backend microservices powering **ScrapeSmith** — a cloud-based SaaS platform for automated web scraping, data cleaning, job scheduling, and AI-powered analysis.

> “A blacksmith for the web — forge your data, sharpen your insights.”

---

## Microservices Overview

| Service                  | Description                                                                                  |
|--------------------------|----------------------------------------------------------------------------------------------|
| `auth-service`           | Handles user registration, login, JWT auth, and role-based permissions                       |
| `user-service`           | Manages user profiles, dashboards, and genral frontend crud operations                       |
| `scraping-service`       | Scrapes websites using ScrapeNinja API and stores raw data in MongoDB                        |
| `data-cleaning-service`  | Cleans scraped JS and HTML using Cheerio and stores cleaned data                             |
| `ai-analysis-service`    | (WIP) Performs AI analysis (e.g., sentiment, clustering) on cleaned data                     |
| `job-scheduling-service` | Handles single and bulk scheduled scrape/clean/analyze jobs with triggers and status tracking|
| `payment-service`        | Integrates Stripe for one-time and subscription payments + voucher support                   |
| `admin-service`          | Admin-only access to manage users, orders, support tickets, and platform stats               |

All services communicate via **RESTful APIs** behind a shared **Kong API Gateway** 

---

## Tech Stack

- **Languages**: Go with Fiber, Node.js (Express), TypeScript
- **Databases**: PostgreSQL, MongoDB Atlas
- **Gateway**: Kong API Gateway
- **Containerization**: Docker
- **Auth**: JWT access + refresh tokens
- **Payments**: Stripe API + voucher support
- **Scraping**: ScrapeNinja API
- **Job Scheduling**: Custom service with queue support (cron/timer-based)
- **AI Analysis**: analysis via ChatGPT third-party APIs

---

## License

This project is proprietary. Please do not use any part of this codebase without prior written permission from the author.

For licensing inquiries, contact https://github.com/SpursFan21



