from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends, Request
from pydantic import BaseModel
from sqlalchemy import create_engine, Column, Integer, String, Text, DateTime, Boolean, ForeignKey, select
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import session, sessionmaker, Session
from datetime import datetime, timezone
from sqlalchemy.exc import OperationalError
import time
import os

DATABASE_URL = os.getenv(key="DATABASE_URL", default="sqlite:///./test.db")

@asynccontextmanager
async def lifespan(app: FastAPI):
    # before yield is before startup
    while True:
        try:
            engine = create_engine(DATABASE_URL)
            with engine.connect() as connection:
                print("Database is ready!")
                break
        except OperationalError:
            print("Database is not ready yet. Retrying in 5 seconds...")
            time.sleep(5)
    yield
    # after yield is code to cleanup or do something after fastapi is closed
    print("Cleanup or something...")

# SQLAlchemy setup
Base = declarative_base()
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# Models
class Post(Base):
    __tablename__ = "posts"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, nullable=False)
    title = Column(String, nullable=False)
    body = Column(Text, nullable=False)
    date_posted = Column(DateTime, default=datetime.now(timezone.utc))
    date_updated = Column(DateTime, default=datetime.now(timezone.utc), onupdate=datetime.now(timezone.utc))


class Comment(Base):
    __tablename__ = "comments"

    id = Column(Integer, primary_key=True, index=True)
    post_id = Column(Integer, ForeignKey("posts.id", ondelete="CASCADE"), nullable=False)
    username = Column(String, nullable=False)
    content = Column(Text, nullable=False)
    date_posted = Column(DateTime, default=datetime.now(timezone.utc))
    is_edited = Column(Boolean, default=False)
    
# API Schemas
class PostCreate(BaseModel):
    # username: str
    title: str
    body: str

class PostResponse(BaseModel):
    id: int
    username: str
    title: str
    body: str
    date_posted: datetime
    date_updated: datetime

    class Config:
        orm_mode = True

class CommentCreate(BaseModel):
    # username: str
    content: str

class CommentResponse(BaseModel):
    id: int
    post_id: int
    username: str
    content: str
    date_posted: datetime
    is_edited: bool

    class Config:
        orm_mode = True

# Create tables
Base.metadata.create_all(bind=engine)

app = FastAPI()

# Routes
@app.post("/posts/", response_model=PostResponse)
def create_post(request: Request, post: PostCreate, db: Session = Depends(get_db)):
    db_post = Post(
        username=request.headers.get("X-Username", None), 
        title=post.title, 
        body=post.body
    )
    db.add(db_post)
    db.commit()
    db.refresh(db_post)
    return db_post

@app.get("/posts/{post_id}", response_model=PostResponse)
def get_post(post_id: int, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    return db_post

@app.post("/comments/{post_id}", response_model=CommentResponse)
def create_comment(request: Request, post_id: int, comment: CommentCreate, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    db_comment = Comment(
        post_id=post_id, 
        username=request.headers.get("X-Username", None), 
        content=comment.content
    )
    db.add(db_comment)
    db.commit()
    db.refresh(db_comment)
    return db_comment

@app.get("/comments/{post_id}", response_model=list[CommentResponse])
def get_comments(post_id: int, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    comments = db.query(Comment).filter(Comment.post_id == post_id).all()
    return comments

@app.put("/posts/{post_id}", response_model=PostResponse)
def update_post(request: Request, post_id: int, post: PostCreate, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    username = request.headers.get("X-Username")
    if username != db_post.username:
        raise HTTPException(status_code=401, detail="Can't edit posts of others")
    db_post.title = post.title
    db_post.body = post.body
    db.commit()
    db.refresh(db_post)
    return db_post

@app.delete("/posts/{post_id}", status_code=204)
def delete_post(request: Request, post_id: int, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    username = request.headers.get("X-Username")
    if username != db_post.username:
        raise HTTPException(status_code=401, detail="Can't delete posts of others")
    db.delete(db_post)
    db.commit()
    return None

@app.put("/comments/{comment_id}", response_model=CommentResponse)
def update_comment(request: Request, comment_id: int, comment: CommentCreate, db: Session = Depends(get_db)):
    db_comment = db.query(Comment).filter(Comment.id == comment_id).first()
    if not db_comment:
        raise HTTPException(status_code=404, detail="Comment not found")
    username = request.headers.get("X-Username")
    if username != db_comment.username:
        raise HTTPException(status_code=401, detail="Unauthorized for this action")
    db_comment.content = comment.content
    db_comment.is_edited = True
    db.commit()
    db.refresh(db_comment)
    return db_comment

@app.delete("/comments/{comment_id}", status_code=204)
def delete_comment(request: Request, comment_id: int, db: Session = Depends(get_db)):
    db_comment = db.query(Comment).filter(Comment.id == comment_id).first()
    if not db_comment:
        raise HTTPException(status_code=404, detail="Comment not found")
    username = request.headers.get("X-Username")
    if username != db_comment.username:
        raise HTTPException(status_code=401, detail="Can't delete comments of others")
    db.delete(db_comment)
    db.commit()
    return None
