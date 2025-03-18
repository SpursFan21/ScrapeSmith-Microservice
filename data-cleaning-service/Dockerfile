# Use official Python image
FROM python:3.11

# Set the working directory
WORKDIR /app

# Copy dependencies first for caching
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the app
COPY . .

# Expose FastAPI default port
EXPOSE 8000

# Run FastAPI with Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
