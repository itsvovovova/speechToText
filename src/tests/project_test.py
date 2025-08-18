import pytest
import requests
import time

BASE_URL = "http://server:8080"

@pytest.fixture(scope="session")
def api_client():
    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})
    yield session
    session.close()

@pytest.fixture(scope="session")
def unique_user():
    return {
        "username": f"test_{int(time.time())}",
        "password": "testpass"
    }

def test_register_success(api_client):
    new_user = {
        "username": f"test_reg_{int(time.time())}_{hash(time.time())}",
        "password": "testpass"
    }
    response = api_client.post(f"{BASE_URL}/register", json=new_user)
    print(f"\nDEBUG: Valid user response status: {response.status_code}")
    print(f"DEBUG: Valid user response body: {response.text}")

    assert response.status_code == 200
    data = response.json()
    assert data["result"] == "ok"

@pytest.fixture(scope="session")
def registered_user(api_client, unique_user):
    response = api_client.post(f"{BASE_URL}/register", json=unique_user)
    assert response.status_code == 200
    return unique_user

@pytest.fixture(scope="session")
def auth_client(api_client, registered_user):
    response = api_client.post(f"{BASE_URL}/login", json=registered_user)
    assert response.status_code == 200
    data = response.json()
    assert data["result"] == "ok"
    assert "token" in data
    
    api_client.cookies.update(response.cookies)
    
    return api_client

def test_register_duplicate(api_client, registered_user):
    response = api_client.post(f"{BASE_URL}/register", json=registered_user)
    assert response.status_code == 400

def test_login_success(api_client, registered_user):
    response = api_client.post(f"{BASE_URL}/login", json=registered_user)

    assert response.status_code == 200
    data = response.json()
    assert data["result"] == "ok"
    assert "token" in data
    assert len(data["token"]) > 0

def test_login_wrong_pass(api_client, registered_user):
    wrong_data = {"username": registered_user["username"], "password": "wrong"}
    response = api_client.post(f"{BASE_URL}/login", json=wrong_data)
    assert response.status_code == 401

def test_login_nonexistent_user(api_client):
    nonexistent_user = {
        "username": f"nonexistent_{int(time.time())}",
        "password": "testpass"
    }
    response = api_client.post(f"{BASE_URL}/login", json=nonexistent_user)
    assert response.status_code == 401

def test_protected_endpoints(auth_client):
    audio_data = {"audio": "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"}
    audio_response = auth_client.post(f"{BASE_URL}/audio", json=audio_data)
    assert audio_response.status_code == 200
    
    status_response = auth_client.get(f"{BASE_URL}/status")
    assert status_response.status_code == 200
    assert status_response.text in ["in progress", "completed", "failed"]

    result_response = auth_client.get(f"{BASE_URL}/result")
    assert result_response.status_code == 200

def test_unauthorized_access(api_client):
    # Очищаем cookies перед тестом
    api_client.cookies.clear()
    
    endpoints = [
        ("GET", "/status"),
        ("GET", "/result"),
        ("POST", "/audio", {"audio": "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"})
    ]

    for method, path, *data in endpoints:
        if method == "GET":
            response = api_client.get(f"{BASE_URL}{path}")
        else:
            response = api_client.post(f"{BASE_URL}{path}", json=data[0])

        assert response.status_code == 401, f"Should be unauthorized for {method} {path}"