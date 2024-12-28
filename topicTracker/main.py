from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os
import aio_pika
import json
from tavily import TavilyClient
from langchain_community.chat_models import ChatOllama
from langchain.prompts import PromptTemplate
from langchain_core.output_parsers import StrOutputParser

# Configuration
TAVILY_API_KEY = os.getenv("TAVILY_API_KEY", "")
RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost/")

# Initialize Tavily Client
tavily_client = TavilyClient(api_key=TAVILY_API_KEY)

# FastAPI app
app = FastAPI()

class TopicRequest(BaseModel):
    topic: str
    email: str

async def fetch_articles(topic: str):
    """Fetch top articles for a given topic using Tavily API."""
    try:
        response = tavily_client.search(topic)
        return response.get("results", [])
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching articles: {str(e)}")

async def summarize_articles(articles):
    """Summarize articles using the LLM."""
    links = []
    content = ""
    for article in articles:
        links.append(article["url"])
        content += article["content"]
    
    llm = ChatOllama(model="phi3", temperature=0)
    prompt = PromptTemplate(
        template="""You are an expert at summarizing snippets of news articles. You are provided some text having three news articles concatenated one after another. Summarize them!
        Text to summarize: {text}""",
        input_variables=["text"],
    )
    summarizer = prompt | llm | StrOutputParser()
    summary = summarizer.invoke({"text": content})
    
    return links, summary

async def send_to_rabbitmq(links, summary, email):
    """Send summaries and links to RabbitMQ."""
    try:
        message = {
            "email": email,
            "links": links,
            "summary": summary,
        }
        connection = await aio_pika.connect_robust(RABBITMQ_URL)
        async with connection:
            channel = await connection.channel()
            await channel.default_exchange.publish(
                aio_pika.Message(body=json.dumps(message).encode()),
                routing_key="llm.request"
            )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error sending message to RabbitMQ: {str(e)}")

@app.post("/search-and-summarize")
async def search_and_summarize(request: TopicRequest):
    """Search the web for the topic, summarize it, and send to RabbitMQ."""
    articles = await fetch_articles(request.topic)
    links, summary = await summarize_articles(articles)
    await send_to_rabbitmq(links, summary, request.email)
    return {"message": "Request Accepted. We will respond via email."}
