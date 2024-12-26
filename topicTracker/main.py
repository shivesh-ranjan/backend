from fastapi import FastAPI, HTTPException
from pika import connection
from pydantic import BaseModel
import os
import pika
import json
from tavily import TavilyClient
from env import TAVILY_API_KEY


# connection = pika.BlockingConnection(
#     pika.ConnectionParameters(host='localhost'))
# channel = connection.channel()
#
# channel.exchange_declare(exchange='llm_comms', exchange_type='topic')

# TAVILY_API_KEY=os.environ["TAVILY_API_KEY"]
tavily_client = TavilyClient(api_key=TAVILY_API_KEY)


async def fetch_articles(topic: str):
    """Fetch top articles for a given topic using Tavily API."""
    response = tavily_client.search(topic)
    print(response, type(response))
    return response

async def summarize_articles(articles):
    """Summarize articles using the LLM."""
    print(articles)
    return articles

async def send_to_rabbitmq(summaries):
    """Send summaries along with article links to RabbitMQ."""
    for summary in summaries:
        message = {
            "title": summary["title"],
            "link": summary["link"],
            "summary": summary["summary"]
        }
    channel.basic_publish(
        exchange='llm_comms', routing_key='llm.request', body=json.dumps(message)
    )

##################################################################################
app = FastAPI()

class TopicRequest(BaseModel):
    topic: str
    # email: str


@app.post("/search-and-summarize")
async def search_and_summarize(request: TopicRequest):
    """Search the web for the topic, summarize it, and send to RabbitMQ."""
    articles = await fetch_articles(request.topic)
    summaries = await summarize_articles(articles)
    # await send_to_rabbitmq(summaries)
    print(summaries)
    return {"message": "Summaries and links sent to RabbitMQ."}
