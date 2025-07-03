from openai import OpenAI
import os
from urllib.request import urlopen, Request
from bs4 import BeautifulSoup
import requests
from dotenv import load_dotenv
import json

load_dotenv()

def moving_average(prices, period):
    return sum(map(float, prices[:period])) / period


def moving_average_crossover(closes):
    if len(closes) < 20:
        return "-"  # not enough data

    short_ma = moving_average(closes, 5)
    long_ma = moving_average(closes, 20)
    today = float(closes[0])

    if short_ma > long_ma and today > short_ma:
        return "BUY"
    elif short_ma < long_ma and today < short_ma:
        return "SELL"
    else:
        return "HOLD"


def fetch_news(symbol):
    """
    Fetch news headlines from Finviz for a given stock symbol.
    """
    website_url = 'https://finviz.com/quote.ashx?t='
    url = website_url + symbol
    req = Request(url=url, headers={'user-agent': 'my-scrape'})
    response = urlopen(req)
    html = BeautifulSoup(response, 'html.parser')
    news_data = html.find(id='news-table')

    headlines = []
    if news_data:
        for row in news_data.find_all('tr')[:20]:
            if row.a:
                title = row.a.text.strip()
                headlines.append(title)
    return headlines


def get_news_sentiment(symbol, headlines):
    """
    Analyze sentiment of the given headlines using OpenAI.
    """
    if not headlines:
        return "No news available"

    client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
    prompt = f"Analyze the sentiment of the following news headlines about {symbol}:\n" + "\n".join(
        headlines) + " \n Provide a sentiment score from 0 (very negative) to 100 (very positive). Return only the score, No summary, No explanation."
    try:
        resp = client.chat.completions.create(
            model="gpt-4o",
            messages=[{"role": "user", "content": prompt}]
        )
        return resp.choices[0].message.content.strip()
    except Exception as e:
        return f"Error analyzing sentiment: {e}"


def main(event):
    symbol = event.get("symbol", "AAPL")

    print(f"Fetching data for symbol: {symbol}")

    key = os.getenv("ALPHA_VANTAGE_API_KEY")

    url = f"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol={symbol}&apikey={key}"
    data = requests.get(url).json()

    if 'Time Series (Daily)' not in data:
        print(f"Error: 'Time Series (Daily)' key not found in the API response.")
        print(f"Response: {data}")
        return

    prices = list(data['Time Series (Daily)'].values())
    closes = [d['4. close'] for d in prices]

    reco = moving_average_crossover(closes)
    headlines = fetch_news(symbol)
    sentiment = get_news_sentiment(symbol, headlines)

    output = {
        "symbol": symbol,
        "moving_average_crossover": reco,
        "sentiment": sentiment,
        "news": headlines
    }

    print(json.dumps(output, indent=2))

    return output