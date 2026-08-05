import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  getPlayers,
  getPlayerStats,
  getPlayerHistory,
  createPlayer,
  updatePlayerName,
  deletePlayer,
  type Player,
  type PlayerStats,
  type MatchHistory,
} from '../lib/api'
import Pagination from '../components/Pagination'
import { useAuth } from '../lib/AuthContext'

const currentYear = new Date().getFullYear()

type Tab = 'profile' | 'add' | 'manage'

function delta(
  current: number,
  previous: number | null
): { text: string; type: 'positive' | 'negative' | 'neutral' } {
  if (previous === null) return { text: '', type: 'neutral' }
  const diff = current - previous
  if (diff === 0) return { text: '±0', type: 'neutral' }
  return diff > 0 ? { text: `+${diff}`, type: 'positive' } : { text: `${diff}`, type: 'negative' }
}

const deltaColor = {
  positive: 'text-emerald-600',
  negative: 'text-red-600',
  neutral: 'text-slate-400',
}

export default function Players() {
  const auth = useAuth()
  const [activeTab, setActiveTab] = useState<Tab>('profile')

  const [players, setPlayers] = useState<Player[]>([])
  const [loadingPlayers, setLoadingPlayers] = useState(false)

  async function fetchPlayers() {
    setLoadingPlayers(true)
    try {
      setPlayers(await getPlayers())
    } catch (e) {
      console.error(e)
    } finally {
      setLoadingPlayers(false)
    }
  }

  useEffect(() => {
    fetchPlayers()
  }, [])

  // --- Profile tab state ---
  const [selectedPlayerId, setSelectedPlayerId] = useState('')
  const [matchType, setMatchType] = useState('indoor')
  const [season, setSeason] = useState(currentYear)
  const [stats, setStats] = useState<PlayerStats | null>(null)
  const [prevStats, setPrevStats] = useState<PlayerStats | null>(null)
  const [history, setHistory] = useState<MatchHistory[]>([])
  const [historyPage, setHistoryPage] = useState(1)
  const [historyHasMore, setHistoryHasMore] = useState(false)
  const [loadingHistory, setLoadingHistory] = useState(false)
  const [loadingStats, setLoadingStats] = useState(false)
  const [statsError, setStatsError] = useState('')

  useEffect(() => {
    if (!selectedPlayerId) return
    let cancelled = false

    async function fetchStats() {
      setLoadingStats(true)
      setStatsError('')
      try {
        const s = await getPlayerStats(selectedPlayerId, matchType, season)
        if (cancelled) return
        setStats(s)
        try {
          const p = await getPlayerStats(selectedPlayerId, matchType, season - 1)
          if (!cancelled) setPrevStats(p)
        } catch {
          if (!cancelled) setPrevStats(null)
        }
      } catch {
        if (!cancelled) {
          setStatsError('Failed to load player stats.')
          setStats(null)
        }
      } finally {
        if (!cancelled) setLoadingStats(false)
      }
    }

    fetchStats()
    return () => {
      cancelled = true
    }
  }, [selectedPlayerId, matchType, season])

  // Reset to page 1 whenever the player/filters change
  useEffect(() => {
    setHistoryPage(1)
  }, [selectedPlayerId, matchType, season])

  useEffect(() => {
    if (!selectedPlayerId) return
    let cancelled = false

    async function fetchHistory() {
      setLoadingHistory(true)
      try {
        const { items, hasMore } = await getPlayerHistory(
          selectedPlayerId,
          matchType,
          season,
          historyPage,
          5
        )
        if (cancelled) return
        setHistory(items)
        setHistoryHasMore(hasMore)
      } catch {
        if (!cancelled) setHistory([])
      } finally {
        if (!cancelled) setLoadingHistory(false)
      }
    }

    fetchHistory()
    return () => {
      cancelled = true
    }
  }, [selectedPlayerId, matchType, season, historyPage])

  // --- Add player tab state ---
  const [newPlayerName, setNewPlayerName] = useState('')
  const [addSuccess, setAddSuccess] = useState('')
  const [addError, setAddError] = useState('')

  async function handleAddPlayer() {
    setAddSuccess('')
    setAddError('')
    try {
      await createPlayer(newPlayerName)
      setAddSuccess(`Added ${newPlayerName}`)
      setNewPlayerName('')
      await fetchPlayers()
    } catch {
      setAddError('Failed to add player.')
    }
  }

  // --- Manage players tab state ---
  const [managePlayerId, setManagePlayerId] = useState('')
  const [newName, setNewName] = useState('')
  const [manageSuccess, setManageSuccess] = useState('')
  const [manageError, setManageError] = useState('')

  async function handleRename() {
    setManageSuccess('')
    setManageError('')
    try {
      await updatePlayerName(managePlayerId, newName)
      setManageSuccess('Renamed successfully')
      await fetchPlayers()
    } catch {
      setManageError('Failed to rename player.')
    }
  }

  async function handleDelete() {
    setManageSuccess('')
    setManageError('')
    try {
      await deletePlayer(managePlayerId)
      setManageSuccess('Player deleted')
      setManagePlayerId('')
      await fetchPlayers()
    } catch {
      setManageError('Failed to delete player.')
    }
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold tracking-tight text-slate-900">
        Players
      </h1>

      <div className="mb-6 flex gap-1 border-b border-slate-200">
        {(
          [
            ['profile', 'Player Profile'],
            ['add', 'Add Player'],
            ['manage', 'Manage Players'],
          ] as [Tab, string][]
        ).map(([tab, label]) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={[
              'rounded-t-lg px-4 py-2 text-sm font-medium transition-colors',
              activeTab === tab
                ? 'border-b-2 border-brand-accent text-brand-accent'
                : 'text-slate-500 hover:text-slate-800',
            ].join(' ')}
          >
            {label}
          </button>
        ))}
      </div>

      {activeTab === 'profile' && (
        <div>
          {loadingPlayers ? (
            <p className="text-sm text-slate-500">Loading players...</p>
          ) : players.length === 0 ? (
            <p className="text-sm text-slate-500">No players yet!</p>
          ) : (
            <>
              <div className="mb-6 flex flex-col gap-3 sm:flex-row">
                <select
                  value={selectedPlayerId}
                  onChange={(e) => setSelectedPlayerId(e.target.value)}
                  className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                >
                  <option value="">Select a player</option>
                  {players.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
                <select
                  value={matchType}
                  onChange={(e) => setMatchType(e.target.value)}
                  className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                >
                  <option value="indoor">Indoor</option>
                  <option value="beach">Beach</option>
                </select>
                <input
                  type="number"
                  min={2023}
                  max={currentYear}
                  value={season}
                  onChange={(e) => setSeason(Number(e.target.value))}
                  className="w-28 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                />
              </div>

              {!selectedPlayerId ? (
                <p className="text-sm text-slate-500">Select a player to view their stats.</p>
              ) : loadingStats ? (
                <p className="text-sm text-slate-500">Loading stats...</p>
              ) : statsError ? (
                <p className="text-sm text-red-600">{statsError}</p>
              ) : stats ? (
                <>
                  <div className="mb-8 grid grid-cols-2 gap-4 md:grid-cols-4">
                    {[
                      { label: 'Matches Played', value: stats.played, prev: prevStats?.played ?? null },
                      { label: 'Wins', value: stats.wins, prev: prevStats?.wins ?? null },
                      { label: 'Losses', value: stats.losses, prev: prevStats?.losses ?? null },
                      { label: 'OTL', value: stats.otl, prev: prevStats?.otl ?? null },
                      { label: 'Points', value: stats.points, prev: prevStats?.points ?? null },
                      {
                        label: 'Win Rate',
                        value: Math.round(stats.win_rate * 100),
                        prev: prevStats ? Math.round(prevStats.win_rate * 100) : null,
                        suffix: '%',
                      },
                      { label: 'Win Streak', value: stats.streak, prev: null },
                      {
                        label: 'Longest Streak',
                        value: stats.longest_streak,
                        prev: prevStats?.longest_streak ?? null,
                      },
                    ].map((stat) => {
                      const d = delta(stat.value, stat.prev)
                      return (
                        <div
                          key={stat.label}
                          className="flex flex-col rounded-xl bg-brand-bg p-4 text-white"
                        >
                          <span className="text-xs text-slate-400">{stat.label}</span>
                          <span className="text-2xl font-semibold">
                            {stat.value}
                            {stat.suffix ?? ''}
                          </span>
                          {d.text && (
                            <span className={`mt-1 text-xs font-medium ${deltaColor[d.type]}`}>
                              {d.text}
                            </span>
                          )}
                        </div>
                      )
                    })}
                  </div>

                  <div>
                    <h3 className="mb-3 text-sm font-medium text-slate-500">Match History</h3>
                    {history.length === 0 ? (
                      <p className="text-sm text-slate-500">No completed matches this season.</p>
                    ) : (
                      <div className="space-y-2">
                        {history.map((match) => {
                          const myScore = match.color === 'blue' ? match.blue_score : match.red_score
                          const theirScore = match.color === 'blue' ? match.red_score : match.blue_score
                          const won = myScore > theirScore
                          const isOtl = Math.abs(match.blue_score - match.red_score) === 2 && !won
                          const rowClass = won
                            ? 'bg-emerald-50 text-emerald-700'
                            : isOtl
                              ? 'bg-amber-50 text-amber-700'
                              : 'bg-red-50 text-red-700'
                          return (
                            <div
                              key={match.id}
                              className={`flex flex-wrap items-center gap-x-6 gap-y-1 rounded-lg px-4 py-2 text-sm ${rowClass}`}
                            >
                              <span>{won ? 'Win' : isOtl ? 'OTL' : 'Loss'}</span>
                              <span>
                                {myScore} : {theirScore}
                              </span>
                              <span className="capitalize">{match.color} team</span>
                              <span>{match.created_at.slice(0, 10)}</span>
                            </div>
                          )
                        })}
                      </div>
                    )}
                    {(history.length > 0 || historyPage > 1) && (
                      <Pagination
                        page={historyPage}
                        hasMore={historyHasMore}
                        onPageChange={setHistoryPage}
                        loading={loadingHistory}
                      />
                    )}
                  </div>
                </>
              ) : null}
            </>
          )}
        </div>
      )}

      {activeTab === 'add' && (
        <div className="max-w-md rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold text-slate-900">Add New Player</h2>
          {auth.isAuthenticated ? (
            <>
              <div className="flex gap-3">
                <input
                  type="text"
                  placeholder="Player name"
                  value={newPlayerName}
                  onChange={(e) => setNewPlayerName(e.target.value)}
                  className="flex-1 rounded-lg border border-slate-200 px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                />
                <button
                  onClick={handleAddPlayer}
                  disabled={!newPlayerName}
                  className="rounded-lg bg-brand-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-accent-hover disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Add Player
                </button>
              </div>
              {addSuccess && <p className="mt-3 text-sm text-emerald-600">{addSuccess}</p>}
              {addError && <p className="mt-3 text-sm text-red-600">{addError}</p>}
            </>
          ) : (
            <p className="text-sm text-slate-500">
              You must be{' '}
              <Link to="/login" className="text-brand-accent hover:underline">
                logged in
              </Link>{' '}
              to add players.
            </p>
          )}
        </div>
      )}

      {activeTab === 'manage' && (
        <div className="max-w-md rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold text-slate-900">Manage Players</h2>
          {!auth.isAuthenticated ? (
            <p className="text-sm text-slate-500">
              You must be{' '}
              <Link to="/login" className="text-brand-accent hover:underline">
                logged in
              </Link>{' '}
              to manage players.
            </p>
          ) : players.length === 0 ? (
            <p className="text-sm text-slate-500">No players to manage.</p>
          ) : (
            <>
              <select
                value={managePlayerId}
                onChange={(e) => setManagePlayerId(e.target.value)}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
              >
                <option value="">Select a player</option>
                {players.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>

              {managePlayerId && (
                <div className="mt-6 space-y-6">
                  <div>
                    <h3 className="mb-2 text-sm font-medium text-slate-700">Rename</h3>
                    <div className="flex gap-3">
                      <input
                        type="text"
                        placeholder="New name"
                        value={newName}
                        onChange={(e) => setNewName(e.target.value)}
                        className="flex-1 rounded-lg border border-slate-200 px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                      />
                      <button
                        onClick={handleRename}
                        disabled={!newName}
                        className="rounded-lg bg-brand-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-accent-hover disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        Rename
                      </button>
                    </div>
                  </div>

                  <div className="rounded-lg border border-red-200 bg-red-50 p-4">
                    <h3 className="mb-1 text-sm font-medium text-red-700">Delete</h3>
                    <p className="mb-3 text-sm text-red-600">
                      This will permanently delete this player.
                    </p>
                    <button
                      onClick={handleDelete}
                      className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
                    >
                      Confirm Delete
                    </button>
                  </div>
                </div>
              )}

              {manageSuccess && <p className="mt-4 text-sm text-emerald-600">{manageSuccess}</p>}
              {manageError && <p className="mt-4 text-sm text-red-600">{manageError}</p>}
            </>
          )}
        </div>
      )}
    </div>
  )
}
