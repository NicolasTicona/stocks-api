# Stock App API

## Overview
The Stock App API is a backend service designed to provide stock data, recommendations, and analysis. It integrates with external APIs like Finnhub, OpenAI, and News APIs to fetch stock quotes, analyze sentiment, and provide actionable insights.

---

## Technologies Used
- **Programming Language**: Go (Golang)
- **Database**: PostgreSQL (CockroachDB)
- **Caching**: Redis
- **External APIs**:
  - Finnhub API (Stock quotes)
  - OpenAI API (Sentiment analysis)
  - News API (News headlines)
- **Python**: For executing stock analysis scripts
- **Docker**: Containerization for production deployment

---

## Dependencies
### Go Dependencies:
- `github.com/gorilla/mux`: HTTP routing
- `github.com/joho/godotenv`: Environment variable management
- `github.com/go-redis/redis/v9`: Redis client
- `gorm.io/gorm`: ORM for database operations
- `gorm.io/driver/postgres`: PostgreSQL driver

### Python Dependencies:
- `beautifulsoup4`: HTML parsing
- `requests`: HTTP requests
- `openai`: OpenAI API integration
- `python-dotenv`: Environment variable management

---

## Features
1. **Stock Synchronization**:
   - Fetches stock ratings and updates the database with stocks and recommendations.
2. **Stock Data Retrieval**:
   - Provides stock data with stocks saved from swechallenge/list, it allows pagination and filter by ticker name
3. **Recommendations**:
   - Generates top stock recommendations based on target and rating score.
4. **Stock Analysis**:
   - Executes Python scripts for moving average crossover analysis and sentiment analysis based on news headlines.
5. **Caching**:
   - Uses Redis to cache stock analysis results for improved performance.

---

## Endpoints
### 1. **Sync Database**
- **URL**: `/sync`
- **Method**: `POST`
- **Description**: Synchronizes stock ratings from external APIs and updates the database.
It searches for all the stocks in swechallenge/list endpoint and insert them in stocks_ratings table and in the same time, add an score based on the rating and target change to calculate which stocks to recommend 

### 2. **Get Stocks**
- **URL**: `/stocks`
- **Method**: `GET`
- **Params**:
  - `page` (optional): Page number for pagination.
  - `limit` (optional): Number of stocks per page.
  - `filter` (optional): Filter by ticker.
- **Description**: Retrieves stock data

### 3. **Get Recommendations**
- **URL**: `/recommendations`
- **Method**: `GET`
- **Description**: Fetches top stock recommendations based on scores.

### 4. **Analyze Stock**
- **URL**: `/analyze`
- **Method**: `GET`
- **Params**:
  - `stock`: Stock ticker symbol.
- **Description**: Executes Python scripts to analyze stock data and sentiment.

### 5. **Test Redis**
- **URL**: `/test-analyze`
- **Method**: `GET`
- **Params**:
  - `stock`: Stock ticker symbol.
- **Description**: Tests Redis caching for stock analysis.

---

## Example `.env` File
```env
ALPHA_VANTAGE_KEY=
OPENAI_API_KEY=
FINNHUB_API_KEY=
REDIS_USERNAME=
REDIS_HOST=
REDIS_PASSWORD=
REDIS_DB=
POSTGRES_DSN=
IS_ENV=