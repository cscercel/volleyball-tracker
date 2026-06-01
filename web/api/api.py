import streamlit as st
import requests
from typing import List, Dict

API_BASE = st.secrets["API_BASE_URL"]


# ---------------------------------------------------------------------------
# Auth
# ---------------------------------------------------------------------------

def login(email: str, password: str) -> Dict:
    """POST /api/v1/users/login  →  { access_token, ... }"""
    response = requests.post(
        f"{API_BASE}/users/login",
        json={"email": email, "password": password},
        headers={"Content-Type": "application/json"},
    )
    response.raise_for_status()
    return response.json()


def get_headers() -> Dict:
    headers = {"Content-Type": "application/json"}
    if "token" in st.session_state:
        headers["Authorization"] = f"Bearer {st.session_state.access_token}"
    return headers


# ---------------------------------------------------------------------------
# Players
# ---------------------------------------------------------------------------

def get_players() -> List[Dict]:
    """GET /api/v1/players  →  []Player  (id, name, created_at, updated_at)"""
    response = requests.get(f"{API_BASE}/players")
    response.raise_for_status()
    return response.json()


def get_player_stats(player_id: str, match_type: str, season: int) -> Dict:
    """GET /api/v1/players/{id}?match_type=&season=  →  GetPlayerStatsByIDRow"""
    response = requests.get(
        f"{API_BASE}/players/{player_id}",
        params={"match_type": match_type, "season": season},
    )
    response.raise_for_status()
    return response.json()


def get_leaderboard(match_type: str, season: int) -> List[Dict]:
    """GET /api/v1/players/leaderboard?match_type=&season=  →  []GetLeaderboardRow"""
    response = requests.get(
        f"{API_BASE}/players/leaderboard",
        params={"match_type": match_type, "season": season},
    )
    response.raise_for_status()
    return response.json()


def get_player_history(player_id: str, match_type: str, season: int) -> List[Dict]:
    """GET /api/v1/players/{id}/history?match_type=&season=  →  []GetPlayerSeasonalMatchesRow"""
    response = requests.get(
        f"{API_BASE}/players/{player_id}/history",
        params={"match_type": match_type, "season": season},
    )
    response.raise_for_status()
    return response.json()


def create_player(name: str) -> Dict:
    """POST /api/v1/players  →  Player"""
    response = requests.post(
        f"{API_BASE}/players",
        json={"name": name},
        headers=get_headers(),
    )
    response.raise_for_status()
    return response.json()


def update_player_name(player_id: str, new_name: str) -> Dict:
    """PUT /api/v1/players/{id}  →  Player"""
    response = requests.put(
        f"{API_BASE}/players/{player_id}",
        json={"name": new_name},
        headers=get_headers(),
    )
    response.raise_for_status()
    return response.json()


def delete_player(player_id: str) -> None:
    """DELETE /api/v1/players/{id}"""
    response = requests.delete(
        f"{API_BASE}/players/{player_id}",
        headers=get_headers(),
    )
    response.raise_for_status()


# ---------------------------------------------------------------------------
# Matches
# ---------------------------------------------------------------------------

def create_match(match_data: Dict) -> Dict:
    """POST /api/v1/matches  →  Match
    match_data: { match_type, blue_team: [name, ...], red_team: [name, ...] }
    """
    response = requests.post(
        f"{API_BASE}/matches",
        json=match_data,
        headers=get_headers(),
    )
    response.raise_for_status()
    return response.json()


def get_match(match_id: str) -> Dict:
    """GET /api/v1/matches/{id}  →  Match"""
    response = requests.get(f"{API_BASE}/matches/{match_id}")
    response.raise_for_status()
    return response.json()


def get_match_roster(match_id: str) -> List[Dict]:
    """GET /api/v1/matches/{id}/roster  →  []GetMatchPlayersRow
    Each row: { color, player_id, player_name }
    """
    response = requests.get(f"{API_BASE}/matches/{match_id}/roster")
    response.raise_for_status()
    return response.json()


def get_uncompleted_matches() -> List[Dict]:
    """GET /api/v1/matches/uncompleted  →  []Match"""
    response = requests.get(f"{API_BASE}/matches/uncompleted")
    response.raise_for_status()
    return response.json()


def get_matches_by_season(match_type: str, season: int) -> List[Dict]:
    """GET /api/v1/matches?match_type=&season=  →  []Match (completed only)"""
    response = requests.get(
        f"{API_BASE}/matches/",
        params={"match_type": match_type, "season": season},
    )
    response.raise_for_status()
    return response.json()


def submit_match_results(match_id: str, blue_score: int, red_score: int) -> Dict:
    """PUT /api/v1/matches/{id}  →  Match"""
    response = requests.put(
        f"{API_BASE}/matches/{match_id}",
        json={"blue_score": blue_score, "red_score": red_score},
        headers=get_headers(),
    )
    response.raise_for_status()
    return response.json()


def delete_match(match_id: str) -> None:
    """DELETE /api/v1/matches/{id}"""
    response = requests.delete(
        f"{API_BASE}/matches/{match_id}",
        headers=get_headers(),
    )
    response.raise_for_status()
