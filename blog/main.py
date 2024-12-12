from fastapi import FastAPI, HTTPException, Depends
from pydantic import BaseModel
from sqlalchemy import create_engine, Column, Integer, String, Text, DateTime, Boolean, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from datetime import datetime

DATABASE_URL = "postgresql://postgres:S&Shivesh72@localhost/blogDB"

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
    date_posted = Column(DateTime, default=datetime.utcnow)
    date_updated = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class Comment(Base):
    __tablename__ = "comments"

    id = Column(Integer, primary_key=True, index=True)
    post_id = Column(Integer, ForeignKey("posts.id", ondelete="CASCADE"), nullable=False)
    username = Column(String, nullable=False)
    content = Column(Text, nullable=False)
    date_posted = Column(DateTime, default=datetime.utcnow)
    is_edited = Column(Boolean, default=False)
    
# API Schemas
class PostCreate(BaseModel):
    username: str
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
    username: str
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
def create_post(post: PostCreate, db: Session = Depends(get_db)):
    db_post = Post(username=post.username, title=post.title, body=post.body)
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
def create_comment(post_id: int, comment: CommentCreate, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    db_comment = Comment(post_id=post_id, username=comment.username, content=comment.content)
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
def update_post(post_id: int, post: PostCreate, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    db_post.username = post.username
    db_post.title = post.title
    db_post.body = post.body
    db.commit()
    db.refresh(db_post)
    return db_post

@app.delete("/posts/{post_id}", status_code=204)
def delete_post(post_id: int, db: Session = Depends(get_db)):
    db_post = db.query(Post).filter(Post.id == post_id).first()
    if not db_post:
        raise HTTPException(status_code=404, detail="Post not found")
    db.delete(db_post)
    db.commit()
    return None

@app.put("/comments/{comment_id}", response_model=CommentResponse)
def update_comment(comment_id: int, comment: CommentCreate, db: Session = Depends(get_db)):
    db_comment = db.query(Comment).filter(Comment.id == comment_id).first()
    if not db_comment:
        raise HTTPException(status_code=404, detail="Comment not found")
    db_comment.username = comment.username
    db_comment.content = comment.content
    db_comment.is_edited = True
    db.commit()
    db.refresh(db_comment)
    return db_comment

@app.delete("/comments/{comment_id}", status_code=204)
def delete_comment(comment_id: int, db: Session = Depends(get_db)):
    db_comment = db.query(Comment).filter(Comment.id == comment_id).first()
    if not db_comment:
        raise HTTPException(status_code=404, detail="Comment not found")
    db.delete(db_comment)
    db.commit()
    return None
