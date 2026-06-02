const API_BASE = import.meta.env.VITE_API_BASE_URL

export type Player = {
  id: string
  name: string
  created_at: string
  updated_at: string
}

export type PlayerStats = {
  name: string
  id: string
  player_id: string
  match_type: string
  season: number
  wins: number
  losses: number
  otl: number
  streak: number
  longest_streak: number
  scored: number
  conceded: number
  points: number
  played: number
  win_rate: number
  efficiency_rate: number
}

export type Match = {
  id: string
  match_type: string
  season: number
  blue_score: number
  red_score: number
  is_completed: boolean
  created_at: string
  updated_at: string
}

export type MatchPlayer = {
  color: string
  player_id: string
  player_name: string
}

export type MatchHistory = {
  id: string
  match_type: string
  season: number
  blue_score: number
  red_score: number
  created_at: string
  color: string
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

export async function login(email: string, password: string) {
  const res = await fetch(`${API_BASE}/users/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error('Login failed')
  const data = await res.json()
  return data.token as string
}

// ---------------------------------------------------------------------------
// Players
// ---------------------------------------------------------------------------

export async function getPlayers(): Promise<Player[]> {
  const res = await fetch(`${API_BASE}/players`)
  if (!res.ok) throw new Error('Failed to fetch players')
  return res.json()
}

export async function getPlayerStats(
  playerId: string,
  matchType: string,
  season: number
): Promise<PlayerStats> {
  const res = await fetch(
    `${API_BASE}/players/${playerId}?match_type=${matchType}&season=${season}`
  )
  if (!res.ok) throw new Error('Failed to fetch player stats')
  return res.json()
}

export async function getLeaderboard(
  matchType: string,
  season: number
): Promise<PlayerStats[]> {
  const res = await fetch(
    `${API_BASE}/players/leaderboard?match_type=${matchType}&season=${season}`
  )
  if (!res.ok) throw new Error('Failed to fetch leaderboard')
  return res.json()
}

export async function getPlayerHistory(
  playerId: string,
  matchType: string,
  season: number
): Promise<MatchHistory[]> {
  const res = await fetch(
    `${API_BASE}/players/${playerId}/history?match_type=${matchType}&season=${season}`
  )
  if (!res.ok) throw new Error('Failed to fetch player history')
  return res.json()
}

export async function createPlayer(name: string): Promise<Player> {
  const res = await fetch(`${API_BASE}/players`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error('Failed to create player')
  return res.json()
}

export async function updatePlayerName(
  playerId: string,
  name: string
): Promise<Player> {
  const res = await fetch(`${API_BASE}/players/${playerId}`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error('Failed to update player')
  return res.json()
}

export async function deletePlayer(playerId: string): Promise<void> {
  const res = await fetch(`${API_BASE}/players/${playerId}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('Failed to delete player')
}

// ---------------------------------------------------------------------------
// Matches
// ---------------------------------------------------------------------------

export async function createMatch(
  matchType: string,
  blueTeam: string[],
  redTeam: string[]
): Promise<Match> {
  const res = await fetch(`${API_BASE}/matches`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ match_type: matchType, blue_team: blueTeam, red_team: redTeam }),
  })
  if (!res.ok) throw new Error('Failed to create match')
  return res.json()
}

export async function getMatch(matchId: string): Promise<Match> {
  const res = await fetch(`${API_BASE}/matches/${matchId}`)
  if (!res.ok) throw new Error('Failed to fetch match')
  return res.json()
}

export async function getMatchRoster(matchId: string): Promise<MatchPlayer[]> {
  const res = await fetch(`${API_BASE}/matches/${matchId}/roster`)
  if (!res.ok) throw new Error('Failed to fetch match roster')
  return res.json()
}

export async function getUncompletedMatches(): Promise<Match[]> {
  const res = await fetch(`${API_BASE}/matches/uncompleted`)
  if (!res.ok) throw new Error('Failed to fetch uncompleted matches')
  return res.json()
}

export async function getMatchesBySeason(
  matchType: string,
  season: number
): Promise<Match[]> {
  const res = await fetch(
    `${API_BASE}/matches/?match_type=${matchType}&season=${season}`
  )
  if (!res.ok) throw new Error('Failed to fetch matches')
  return res.json()
}

export async function submitMatchResults(
  matchId: string,
  blueScore: number,
  redScore: number
): Promise<Match> {
  const res = await fetch(`${API_BASE}/matches/${matchId}`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify({ blue_score: blueScore, red_score: redScore }),
  })
  if (!res.ok) throw new Error('Failed to submit results')
  return res.json()
}

export async function deleteMatch(matchId: string): Promise<void> {
  const res = await fetch(`${API_BASE}/matches/${matchId}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error('Failed to delete match')
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function authHeaders(): Record<string, string> {
  const token = localStorage.getItem('access_token')
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
}
