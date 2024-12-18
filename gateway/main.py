from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import JWTError, jwt
from pydantic import BaseModel
from requests.sessions import Request
from sqlalchemy import create_engine, Column, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
import os
import requests
from datetime import datetime, timedelta, timezone

from sqlalchemy.sql.functions import current_user

# Environment Variables
BLOG_SERVICE_HOST = os.getenv("BLOG_SERVICE_HOST", "http://localhost:8000")
DATABASE_URL = os.getenv("DATABASE_URL", "sqlite:///./test.db")
SECRET_KEY = os.getenv("SECRET_KEY", "derekshawisback")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 1440

# Database Setup
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# Models
class User(Base):
    __tablename__ = "users"

    username = Column(String, primary_key=True, index=True)
    hashed_password = Column(String)

Base.metadata.create_all(bind=engine)

# Authentication
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

class Token(BaseModel):
    access_token: str
    token_type: str

class TokenData(BaseModel):
    username: str | None = None

class UserInDB(User):
    hashed_password: str

def fake_hash_password(password: str):
    return "fakehashed" + password

def verify_password(plain_password: str, hashed_password: str):
    return hashed_password == fake_hash_password(plain_password)

def authenticate_user(db: Session, username: str, password: str):
    user = db.query(User).filter(User.username == username).first()
    if not user or not verify_password(password, user.hashed_password.__str__()):
        return None
    return user

def create_access_token(data: dict, expires_delta: timedelta | None = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now() + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

def get_current_user(token: str = Depends(oauth2_scheme), db: Session = Depends(get_db)):
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username = payload.get("sub")
        if username is None:
            raise credentials_exception
        token_data = TokenData(username=username)
    except JWTError:
        raise credentials_exception
    user = db.query(User).filter(User.username == token_data.username).first()
    if user is None:
        raise credentials_exception
    return user

# App Setup
app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Routes
@app.post("/token", response_model=Token)
def login_for_access_token(form_data: OAuth2PasswordRequestForm = Depends(), db: Session = Depends(get_db)):
    user = authenticate_user(db, form_data.username, form_data.password)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user.username}, expires_delta=access_token_expires
    )
    return {"access_token": access_token, "token_type": "bearer"}

@app.get("/users/me")
def read_users_me(current_user: User = Depends(get_current_user)):
    return current_user

def gateway(full_path: str, request: Request, username: str):
    json_data = None
    if request.method in ("POST", "PUT", "PATCH"):
        json_data = request.json()
        json_data["username"] = username
    response: requests.Response = requests.Response()
    try:
        response = requests.request(
                request.method, f"{full_path}", headers=request.headers, params=request.params, json=json_data
        )
        response.raise_for_status()
        return response.json() if response.content else {"detail": "success"}
    except:
        raise HTTPException(status_code=response.status_code, detail=response.text)

def simpgateway(full_path: str, request:Request):
    response = requests.Response()
    try:
        response = requests.request(
                request.method, f"{full_path}", headers=request.headers, params=request.params, json=request.json
        )
        response.raise_for_status()
        return response.json() if response.content else {"detail": "success"}
    except:
        raise HTTPException(status_code=response.status_code, detail=response.text)

@app.api_route("/{full_path:path}", methods=["POST", "PUT", "DELETE", "PATCH"])
async def blogGateway(full_path: str, request: Request, current_user: User = Depends(get_current_user)):
    return gateway(full_path=BLOG_SERVICE_HOST + "/" + full_path, request=request, username=current_user.username.__str__())

@app.api_route("/{full_path:path}", methods=["GET"])
async def simpBlogGateway(full_path: str, request: Request):
    return simpgateway(full_path=BLOG_SERVICE_HOST + "/" + full_path, request=request)
