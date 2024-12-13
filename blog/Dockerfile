FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt
RUN apt update && apt install -y netcat-traditional

COPY main.py main.py
COPY wait-for.sh .
COPY start.sh .

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
ENTRYPOINT ["/app/start.sh"]
