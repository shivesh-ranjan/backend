import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from main import app, Base, get_db

# Use an in-memory SQLite database for testing
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Create a new database and tables for testing
Base.metadata.create_all(bind=engine)

def override_get_db():
    db = TestingSessionLocal()
    try:
        yield db
    finally:
        db.close()

# Override the dependency in the app
app.dependency_overrides[get_db] = override_get_db

client = TestClient(app)

@pytest.fixture(scope="function")
def setup_test_database():
    # Ensure the database is empty before each test
    Base.metadata.drop_all(bind=engine)
    Base.metadata.create_all(bind=engine)

@pytest.mark.usefixtures("setup_test_database")
def test_create_and_get_post():
    response = client.post("/posts/", json={
        "username": "testuser",
        "title": "Test Title",
        "body": "Test Body"
    })
    assert response.status_code == 200
    post = response.json()
    assert post["username"] == "testuser"
    assert post["title"] == "Test Title"
    assert post["body"] == "Test Body"

    post_id = post["id"]
    response = client.get(f"/posts/{post_id}")
    assert response.status_code == 200
    fetched_post = response.json()
    assert fetched_post == post

def test_create_comment():
    # Create a post first
    post_response = client.post("/posts/", json={
        "username": "testuser",
        "title": "Test Title",
        "body": "Test Body"
    })
    post_id = post_response.json()["id"]

    # Add a comment to the post
    response = client.post(f"/comments/{post_id}", json={
        "username": "commenter",
        "content": "This is a comment."
    })
    assert response.status_code == 200
    comment = response.json()
    assert comment["username"] == "commenter"
    assert comment["content"] == "This is a comment."
    assert comment["post_id"] == post_id

def test_update_post():
    # Create a post first
    response = client.post("/posts/", json={
        "username": "testuser",
        "title": "Test Title",
        "body": "Test Body"
    })
    post_id = response.json()["id"]

    # Update the post
    response = client.put(f"/posts/{post_id}", json={
        "username": "updateduser",
        "title": "Updated Title",
        "body": "Updated Body"
    })
    assert response.status_code == 200
    updated_post = response.json()
    assert updated_post["username"] == "updateduser"
    assert updated_post["title"] == "Updated Title"
    assert updated_post["body"] == "Updated Body"

def test_delete_post():
    # Create a post first
    response = client.post("/posts/", json={
        "username": "testuser",
        "title": "Test Title",
        "body": "Test Body"
    })
    post_id = response.json()["id"]

    # Delete the post
    response = client.delete(f"/posts/{post_id}")
    assert response.status_code == 204

    # Ensure the post no longer exists
    response = client.get(f"/posts/{post_id}")
    assert response.status_code == 404

def test_get_comments():
    # Create a post first
    post_response = client.post("/posts/", json={
        "username": "testuser",
        "title": "Test Title",
        "body": "Test Body"
    })
    post_id = post_response.json()["id"]

    # Add comments to the post
    client.post(f"/comments/{post_id}", json={
        "username": "commenter1",
        "content": "Comment 1"
    })
    client.post(f"/comments/{post_id}", json={
        "username": "commenter2",
        "content": "Comment 2"
    })

    # Fetch comments
    response = client.get(f"/comments/{post_id}")
    assert response.status_code == 200
    comments = response.json()
    assert len(comments) == 2
    assert comments[0]["content"] == "Comment 1"
    assert comments[1]["content"] == "Comment 2"

