from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import requests
import aio_pika

app = FastAPI()

# Configuration
TAVILY_API_KEY = "your_tavily_api_key"
TAVILY_SEARCH_URL = "https://api.tavily.com/v1/search"
LLM_API_URL = "http://localhost:11434/summarize"
RABBITMQ_URL = "amqp://user:password@localhost/"

# Pydantic models
class TopicRequest(BaseModel):
    topic: str

async def fetch_articles(topic: str):
    """Fetch top articles for a given topic using Tavily API."""
    headers = {"Authorization": f"Bearer {TAVILY_API_KEY}"}
    params = {"query": topic, "count": 3}
    response = requests.get(TAVILY_SEARCH_URL, headers=headers, params=params)
    if response.status_code != 200:
        raise HTTPException(status_code=500, detail="Error fetching articles from Tavily API.")
    return response.json().get("results", [])

async def summarize_articles(articles: list):
    """Summarize articles using the LLM."""
    summaries = []
    for article in articles:
        response = requests.post(LLM_API_URL, json={"text": article["content"]})
        if response.status_code == 200:
            summary = response.json().get("summary", "")
            summaries.append({"title": article["title"], "link": article["url"], "summary": summary})
    return summaries

async def send_to_rabbitmq(summaries):
    """Send summaries along with article links to RabbitMQ."""
    async with aio_pika.connect(RABBITMQ_URL) as connection:
        async with connection.channel() as channel:
            queue = await channel.declare_queue("updates")
            for summary in summaries:
                message = {
                    "title": summary["title"],
                    "link": summary["link"],
                    "summary": summary["summary"]
                }
                await queue.publish(
                    aio_pika.Message(body=str(message).encode()),
                    routing_key="updates"
                )

@app.post("/search-and-summarize")
async def search_and_summarize(request: TopicRequest):
    """Search the web for the topic, summarize it, and send to RabbitMQ."""
    articles = await fetch_articles(request.topic)
    summaries = await summarize_articles(articles)
    await send_to_rabbitmq(summaries)
    return {"message": "Summaries and links sent to RabbitMQ."}

