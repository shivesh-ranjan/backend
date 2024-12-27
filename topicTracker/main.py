from typing import Dict
from fastapi import FastAPI, HTTPException
from pika import connection
from pydantic import BaseModel
import os
import pika
import json
from tavily import TavilyClient
from env import TAVILY_API_KEY
from langchain_community.chat_models import ChatOllama
from langchain.prompts import PromptTemplate
from langchain_core.output_parsers import StrOutputParser

connection = pika.BlockingConnection(
    pika.ConnectionParameters(host='localhost'))
channel = connection.channel()

channel.exchange_declare(exchange='llm_comms', exchange_type='topic')

# TAVILY_API_KEY=os.getenv(key="TAVILY_API_KEY", default="")
tavily_client = TavilyClient(api_key=TAVILY_API_KEY)


async def fetch_articles(topic: str):
    """Fetch top articles for a given topic using Tavily API."""
    response = tavily_client.search(topic)
    return response['results']

async def summarize_articles(articles):
    """Summarize articles using the LLM."""
    links = []
    content = ""
    for article in articles:
        links.append(article['url'])
        content += article['content']
    llm = ChatOllama(model="phi3", temperature=0)
    prompt = PromptTemplate(
        template="""You are an expert at summarizing snippets of news articles. You are provided some text having three news articles concatenated one after another. Summarize them!
        Text to summarize: {text}""",
        input_variables=["text"],
    )
    summarizer = prompt | llm | StrOutputParser()
    summary = summarizer.invoke({"text" : content})
    print("Links: ", links)
    print("Content: ", content)
    print("Summary: ", summary)
    return links, summary

async def send_to_rabbitmq(articles, email) -> None:
    """Send summaries along with article links to RabbitMQ."""
    links, summary = await summarize_articles(articles)
    message = {}
    message['email'] = email
    message['links'] = links
    message['summary'] = summary
    channel.basic_publish(
        exchange='llm_comms', routing_key='llm.request', body=json.dumps(message)
    )
    return

##################################################################################
app = FastAPI()

class TopicRequest(BaseModel):
    topic: str
    email: str

class GptRequest(BaseModel):
    query: str

@app.post("/search-and-summarize")
async def search_and_summarize(request: TopicRequest):
    """Search the web for the topic, summarize it, and send to RabbitMQ."""
    articles = await fetch_articles(request.topic)
    send_to_rabbitmq(articles, request.email)
    return {"message": "Request Accepted, we will respond via email."}

@app.post("/phi3-gpt")
async def phi3_gpt(request: GptRequest):
    return
