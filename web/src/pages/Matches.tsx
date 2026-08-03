import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  getPlayers,
  createMatch,
  getMatchRoster,
  getUncompletedMatches,
  getMatchesBySeason,
  submitMatchResults,
  deleteMatch,
  type Player,
  type Match,
  type MatchPlayer,
} from '../lib/api'
import Court from '../components/Court'
import { useAuth } from '../lib/AuthContext'

const currentYear = new Date().getFullYear()

type Tab = 'create' | 'drafts' | 'completed'

export default function Matches() {
  const auth = useAuth()
  const [activeTab, setActiveTab] = useState<Tab>('create')

  const [players, setPlayers] = useState<Player[]>([])

  useEffect(() => {
    getPlayers()
      .then(setPlayers)
      .catch((e) => console.error(e))
  }, [])

  // ---------------------------------------------------------------------
  // TAB 1 - Create Match
  // ---------------------------------------------------------------------
  const [matchType, setMatchType] = useState('indoor')
  const [draftType, setDraftType] = useState<'random' | 'manual'>('random')
  const [blueTeam, setBlueTeam] = useState<string[]>([])
  const [redTeam, setRedTeam] = useState<string[]>([])
  const [createError, setCreateError] = useState('')
  const [createSuccess, setCreateSuccess] = useState('')
  const [createdMatchRoster, setCreatedMatchRoster] = useState<MatchPlayer[]>([])

  const [blueSelected, setBlueSelected] = useState('')
  const [redSelected, setRedSelected] = useState('')

  const availableForBlue = useMemo(
    () => players.filter((p) => !blueTeam.includes(p.name) && !redTeam.includes(p.name)),
    [players, blueTeam, redTeam]
  )
  const availableForRed = useMemo(
    () => players.filter((p) => !redTeam.includes(p.name) && !blueTeam.includes(p.name)),
    [players, blueTeam, redTeam]
  )

  function addToBlue(name: string) {
    if (name && !blueTeam.includes(name)) {
      setBlueTeam((t) => [...t, name])
      setBlueSelected('')
    }
  }

  function addToRed(name: string) {
    if (name && !redTeam.includes(name)) {
      setRedTeam((t) => [...t, name])
      setRedSelected('')
    }
  }

  function removeFromBlue(name: string) {
    setBlueTeam((t) => t.filter((p) => p !== name))
  }

  function removeFromRed(name: string) {
    setRedTeam((t) => t.filter((p) => p !== name))
  }

  async function handleCreateMatch() {
    setCreateError('')
    setCreateSuccess('')
    setCreatedMatchRoster([])

    if (blueTeam.length === 0 || redTeam.length === 0) {
      setCreateError('Both teams need at least one player.')
      return
    }

    let finalBlue = blueTeam
    let finalRed = redTeam

    if (draftType === 'random') {
      const all = [...blueTeam, ...redTeam].sort(() => Math.random() - 0.5)
      const mid = Math.floor(all.length / 2)
      finalBlue = all.slice(0, mid)
      finalRed = all.slice(mid)
    }

    try {
      const match = await createMatch(matchType, finalBlue, finalRed)
      setCreatedMatchRoster(await getMatchRoster(match.id))
      setCreateSuccess('✅ Match created!')
      setBlueTeam([])
      setRedTeam([])
    } catch {
      setCreateError('Failed to create match.')
    }
  }

  const createdBluePlayers = createdMatchRoster
    .filter((p) => p.color === 'blue')
    .map((p) => p.player_name)
  const createdRedPlayers = createdMatchRoster
    .filter((p) => p.color === 'red')
    .map((p) => p.player_name)

  // ---------------------------------------------------------------------
  // TAB 2 - Draft Matches
  // ---------------------------------------------------------------------
  const [drafts, setDrafts] = useState<Match[]>([])
  const [draftRosters, setDraftRosters] = useState<Record<string, MatchPlayer[]>>({})
  const [loadingDrafts, setLoadingDrafts] = useState(false)
  const [blueScores, setBlueScores] = useState<Record<string, number | ''>>({})
  const [redScores, setRedScores] = useState<Record<string, number | ''>>({})

  async function fetchDrafts() {
    setLoadingDrafts(true)
    try {
      const fetchedDrafts = await getUncompletedMatches()
      setDrafts(fetchedDrafts)
      const rosters: Record<string, MatchPlayer[]> = {}
      for (const match of fetchedDrafts) {
        rosters[match.id] = await getMatchRoster(match.id)
      }
      setDraftRosters(rosters)
    } catch {
      console.error('Failed to load drafts')
    } finally {
      setLoadingDrafts(false)
    }
  }

  async function handleSubmitResults(matchId: string) {
    const blue = blueScores[matchId]
    const red = redScores[matchId]
    if (blue === undefined || blue === '' || red === undefined || red === '') {
      alert('Enter scores for both teams.')
      return
    }
    if (blue === red) {
      alert('Scores cannot be equal.')
      return
    }
    try {
      await submitMatchResults(matchId, blue, red)
      await fetchDrafts()
    } catch {
      alert('Failed to submit results.')
    }
  }

  async function handleDeleteDraft(matchId: string) {
    try {
      await deleteMatch(matchId)
      await fetchDrafts()
    } catch {
      alert('Failed to delete draft.')
    }
  }

  // ---------------------------------------------------------------------
  // TAB 3 - Completed Matches
  // ---------------------------------------------------------------------
  const [completedType, setCompletedType] = useState('indoor')
  const [completedSeason, setCompletedSeason] = useState(currentYear)
  const [completed, setCompleted] = useState<Match[]>([])
  const [completedRosters, setCompletedRosters] = useState<Record<string, MatchPlayer[]>>({})
  const [loadingCompleted, setLoadingCompleted] = useState(false)

  async function fetchCompleted() {
    setLoadingCompleted(true)
    try {
      const fetchedCompleted = await getMatchesBySeason(completedType, completedSeason)
      setCompleted(fetchedCompleted)
      const rosters: Record<string, MatchPlayer[]> = {}
      for (const match of fetchedCompleted) {
        rosters[match.id] = await getMatchRoster(match.id)
      }
      setCompletedRosters(rosters)
    } catch {
      console.error('Failed to load completed matches')
    } finally {
      setLoadingCompleted(false)
    }
  }

  useEffect(() => {
    if (activeTab === 'drafts') fetchDrafts()
    if (activeTab === 'completed') fetchCompleted()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab])

  useEffect(() => {
    if (activeTab === 'completed') fetchCompleted()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [completedType, completedSeason])

  return (
    <div>
      <h1 className="mb-6 text-2xl font-semibold tracking-tight text-slate-900">
        🏐 Matches
      </h1>

      <div className="mb-6 flex gap-1 border-b border-slate-200">
        {(
          [
            ['create', 'Create Match'],
            ['drafts', 'Draft Matches'],
            ['completed', 'Completed Matches'],
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

      {/* TAB 1 - Create Match */}
      {activeTab === 'create' && (
        <div>
          <h2 className="mb-4 text-lg font-semibold text-slate-900">Create New Match</h2>
          {players.length < 2 ? (
            <p className="text-sm text-slate-500">⚠️ Need at least 2 players to create a match.</p>
          ) : !auth.isAuthenticated ? (
            <p className="text-sm text-slate-500">
              You must be{' '}
              <Link to="/login" className="text-brand-accent hover:underline">
                logged in
              </Link>{' '}
              to create matches.
            </p>
          ) : (
            <>
              <div className="mb-6 flex flex-wrap items-center gap-6">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-slate-600">Match Type</span>
                  {(['indoor', 'beach'] as const).map((t) => (
                    <button
                      key={t}
                      onClick={() => setMatchType(t)}
                      className={[
                        'rounded-full border px-3 py-1 text-sm capitalize transition-colors',
                        matchType === t
                          ? 'border-brand-accent bg-brand-accent text-white'
                          : 'border-slate-300 bg-white text-slate-600 hover:border-brand-accent/50',
                      ].join(' ')}
                    >
                      {t}
                    </button>
                  ))}
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-slate-600">Draft Type</span>
                  {(['random', 'manual'] as const).map((t) => (
                    <button
                      key={t}
                      onClick={() => setDraftType(t)}
                      className={[
                        'rounded-full border px-3 py-1 text-sm capitalize transition-colors',
                        draftType === t
                          ? 'border-brand-accent bg-brand-accent text-white'
                          : 'border-slate-300 bg-white text-slate-600 hover:border-brand-accent/50',
                      ].join(' ')}
                    >
                      {t}
                    </button>
                  ))}
                </div>
              </div>

              <div className="mb-4 grid grid-cols-2 gap-4">
                <div className="rounded-xl border border-blue-900 bg-blue-950 p-4">
                  <h3 className="mb-3 font-semibold text-white">🔵 Blue Team</h3>
                  <div className="mb-3 flex min-h-[60px] flex-wrap gap-2">
                    {blueTeam.length === 0 ? (
                      <p className="text-sm italic text-white/40">No players added yet</p>
                    ) : (
                      blueTeam.map((name) => (
                        <div
                          key={name}
                          className="flex items-center gap-1.5 rounded-full bg-white/15 px-3 py-1 text-sm text-white"
                        >
                          <span>{name}</span>
                          <button
                            onClick={() => removeFromBlue(name)}
                            className="text-white/60 hover:text-white"
                          >
                            ✕
                          </button>
                        </div>
                      ))
                    )}
                  </div>
                  <select
                    value={blueSelected}
                    onChange={(e) => {
                      setBlueSelected(e.target.value)
                      addToBlue(e.target.value)
                    }}
                    className="w-full rounded-lg px-2 py-1.5 text-sm text-slate-800"
                  >
                    <option value="">Add a player...</option>
                    {availableForBlue.map((p) => (
                      <option key={p.id} value={p.name}>
                        {p.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="rounded-xl border border-red-900 bg-red-950 p-4">
                  <h3 className="mb-3 font-semibold text-white">🔴 Red Team</h3>
                  <div className="mb-3 flex min-h-[60px] flex-wrap gap-2">
                    {redTeam.length === 0 ? (
                      <p className="text-sm italic text-white/40">No players added yet</p>
                    ) : (
                      redTeam.map((name) => (
                        <div
                          key={name}
                          className="flex items-center gap-1.5 rounded-full bg-white/15 px-3 py-1 text-sm text-white"
                        >
                          <span>{name}</span>
                          <button
                            onClick={() => removeFromRed(name)}
                            className="text-white/60 hover:text-white"
                          >
                            ✕
                          </button>
                        </div>
                      ))
                    )}
                  </div>
                  <select
                    value={redSelected}
                    onChange={(e) => {
                      setRedSelected(e.target.value)
                      addToRed(e.target.value)
                    }}
                    className="w-full rounded-lg px-2 py-1.5 text-sm text-slate-800"
                  >
                    <option value="">Add a player...</option>
                    {availableForRed.map((p) => (
                      <option key={p.id} value={p.name}>
                        {p.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <button
                onClick={handleCreateMatch}
                className="rounded-lg bg-brand-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-accent-hover"
              >
                Create Match
              </button>

              {createSuccess && <p className="mt-3 text-sm text-emerald-600">{createSuccess}</p>}
              {createError && <p className="mt-3 text-sm text-red-600">{createError}</p>}

              {createdMatchRoster.length > 0 && (
                <>
                  <div className="mt-4 grid grid-cols-2 gap-4">
                    <div className="rounded-xl bg-blue-950 p-4">
                      <h3 className="mb-2 font-semibold text-white">🔵 Blue Team</h3>
                      {createdBluePlayers.map((p) => (
                        <p key={p} className="text-sm text-white/80">
                          {p}
                        </p>
                      ))}
                    </div>
                    <div className="rounded-xl bg-red-950 p-4">
                      <h3 className="mb-2 font-semibold text-white">🔴 Red Team</h3>
                      {createdRedPlayers.map((p) => (
                        <p key={p} className="text-sm text-white/80">
                          {p}
                        </p>
                      ))}
                    </div>
                  </div>
                  <Court blueTeam={createdBluePlayers} redTeam={createdRedPlayers} />
                </>
              )}
            </>
          )}
        </div>
      )}

      {/* TAB 2 - Draft Matches */}
      {activeTab === 'drafts' && (
        <div>
          <h2 className="mb-4 text-lg font-semibold text-slate-900">Draft Matches</h2>
          {loadingDrafts ? (
            <p className="text-sm text-slate-500">Loading...</p>
          ) : drafts.length === 0 ? (
            <p className="text-sm text-slate-500">
              No draft matches. Create one in the Create Match tab!
            </p>
          ) : (
            <div className="space-y-4">
              {drafts.map((match) => {
                const roster = draftRosters[match.id] ?? []
                const blue = roster.filter((p) => p.color === 'blue')
                const red = roster.filter((p) => p.color === 'red')
                return (
                  <div
                    key={match.id}
                    className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
                  >
                    <div className="mb-3 flex gap-4 text-sm font-medium text-slate-500">
                      <span>Match {match.id.slice(0, 8)}</span>
                      <span>{match.match_type.toUpperCase()}</span>
                    </div>

                    <div className="mb-3 grid grid-cols-2 gap-4">
                      <div className="rounded-xl bg-blue-950 p-4">
                        <h3 className="mb-2 font-semibold text-white">🔵 Blue Team</h3>
                        {blue.map((p) => (
                          <p key={p.player_id} className="text-sm text-white/80">
                            {p.player_name}
                          </p>
                        ))}
                      </div>
                      <div className="rounded-xl bg-red-950 p-4">
                        <h3 className="mb-2 font-semibold text-white">🔴 Red Team</h3>
                        {red.map((p) => (
                          <p key={p.player_id} className="text-sm text-white/80">
                            {p.player_name}
                          </p>
                        ))}
                      </div>
                    </div>

                    {auth.isAuthenticated && (
                      <div className="flex flex-wrap items-center gap-3">
                        <input
                          type="number"
                          placeholder="Blue score"
                          value={blueScores[match.id] ?? ''}
                          onChange={(e) =>
                            setBlueScores((s) => ({
                              ...s,
                              [match.id]: e.target.value === '' ? '' : Number(e.target.value),
                            }))
                          }
                          className="w-32 rounded-lg border border-slate-200 px-3 py-1.5 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                        />
                        <input
                          type="number"
                          placeholder="Red score"
                          value={redScores[match.id] ?? ''}
                          onChange={(e) =>
                            setRedScores((s) => ({
                              ...s,
                              [match.id]: e.target.value === '' ? '' : Number(e.target.value),
                            }))
                          }
                          className="w-32 rounded-lg border border-slate-200 px-3 py-1.5 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
                        />
                        <button
                          onClick={() => handleSubmitResults(match.id)}
                          className="rounded-lg bg-brand-accent px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-brand-accent-hover"
                        >
                          Submit Results
                        </button>
                        <button
                          onClick={() => handleDeleteDraft(match.id)}
                          className="rounded-lg bg-red-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-red-700"
                        >
                          Delete Draft
                        </button>
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </div>
      )}

      {/* TAB 3 - Completed Matches */}
      {activeTab === 'completed' && (
        <div>
          <h2 className="mb-4 text-lg font-semibold text-slate-900">Completed Matches</h2>

          <div className="mb-6 flex gap-3">
            <select
              value={completedType}
              onChange={(e) => setCompletedType(e.target.value)}
              className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
            >
              <option value="indoor">Indoor</option>
              <option value="beach">Beach</option>
            </select>
            <input
              type="number"
              min={2023}
              max={currentYear}
              value={completedSeason}
              onChange={(e) => setCompletedSeason(Number(e.target.value))}
              className="w-28 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
            />
          </div>

          {loadingCompleted ? (
            <p className="text-sm text-slate-500">Loading...</p>
          ) : completed.length === 0 ? (
            <p className="text-sm text-slate-500">No completed matches for this filter.</p>
          ) : (
            <div className="space-y-4">
              {completed.map((match) => {
                const isOtl = Math.abs(match.blue_score - match.red_score) === 2
                const winner = match.blue_score > match.red_score ? 'blue' : 'red'
                const roster = completedRosters[match.id] ?? []
                const blue = roster.filter((p) => p.color === 'blue')
                const red = roster.filter((p) => p.color === 'red')
                return (
                  <div
                    key={match.id}
                    className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
                  >
                    <div className="mb-3 flex flex-wrap gap-4 text-sm font-medium text-slate-500">
                      <span>
                        {winner === 'blue' ? '🔵' : '🔴'} {match.blue_score} – {match.red_score}{' '}
                        {isOtl ? '⏱️ OT' : ''}
                      </span>
                      <span>{match.match_type.toUpperCase()}</span>
                      <span>{match.created_at.slice(0, 10)}</span>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div className="rounded-xl bg-blue-950 p-4">
                        <h3 className="mb-2 font-semibold text-white">
                          🔵 Blue — {match.blue_score}
                        </h3>
                        {blue.map((p) => (
                          <p key={p.player_id} className="text-sm text-white/80">
                            {p.player_name}
                          </p>
                        ))}
                      </div>
                      <div className="rounded-xl bg-red-950 p-4">
                        <h3 className="mb-2 font-semibold text-white">
                          🔴 Red — {match.red_score}
                        </h3>
                        {red.map((p) => (
                          <p key={p.player_id} className="text-sm text-white/80">
                            {p.player_name}
                          </p>
                        ))}
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
