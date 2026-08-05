import { useEffect, useMemo, useState } from 'react'
import { getLeaderboard, type PlayerStats } from '../lib/api'
import { getRank } from '../lib/rank'

const currentYear = new Date().getFullYear()

type SortKey = keyof PlayerStats
type SortDir = 'asc' | 'desc'

const columns: { key: SortKey; label: string }[] = [
  { key: 'name', label: 'Player' },
  { key: 'played', label: 'Played' },
  { key: 'wins', label: 'Wins' },
  { key: 'losses', label: 'Losses' },
  { key: 'otl', label: 'OTL' },
  { key: 'points', label: 'Points' },
  { key: 'win_rate', label: 'Win Rate' },
]

export default function Home() {
  const [matchType, setMatchType] = useState('indoor')
  const [season, setSeason] = useState(currentYear)
  const [leaderboard, setLeaderboard] = useState<PlayerStats[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [sortColumn, setSortColumn] = useState<SortKey>('points')
  const [sortDir, setSortDir] = useState<SortDir>('desc')

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError('')
    getLeaderboard(matchType, season)
      .then((data) => {
        if (!cancelled) setLeaderboard(data)
      })
      .catch(() => {
        if (!cancelled) setError('Failed to load leaderboard.')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [matchType, season])

  function toggleSort(col: SortKey) {
    if (sortColumn === col) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortColumn(col)
      setSortDir('desc')
    }
  }

  function sortIcon(col: SortKey) {
    if (sortColumn !== col) return '↕'
    return sortDir === 'asc' ? '↑' : '↓'
  }

  const sorted = useMemo(() => {
    return [...leaderboard].sort((a, b) => {
      const aVal = a[sortColumn]
      const bVal = b[sortColumn]
      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return sortDir === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal)
      }
      return sortDir === 'asc'
        ? (aVal as number) - (bVal as number)
        : (bVal as number) - (aVal as number)
    })
  }, [leaderboard, sortColumn, sortDir])

  return (
    <div>
      <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-semibold tracking-tight text-slate-900">
          Leaderboard
        </h1>

        <div className="flex gap-3">
          <select
            value={matchType}
            onChange={(e) => setMatchType(e.target.value)}
            className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
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
            className="w-28 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-sm focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
          />
        </div>
      </div>

      <div className="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
        <div className="min-w-[720px]">
        {loading ? (
          <p className="p-6 text-sm text-slate-500">Loading...</p>
        ) : error ? (
          <p className="p-6 text-sm text-red-600">{error}</p>
        ) : leaderboard.length === 0 ? (
          <p className="p-6 text-sm text-slate-500">
            No stats for {matchType} in season {season}.
          </p>
        ) : (
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr className="bg-brand-bg text-white">
                <th className="px-4 py-3 text-left font-medium">#</th>
                {columns.map((col) => (
                  <th
                    key={col.key}
                    onClick={() => toggleSort(col.key)}
                    className="cursor-pointer select-none px-4 py-3 text-left font-medium transition-colors hover:bg-brand-bg-hover"
                  >
                    {col.label}{' '}
                    <span className="ml-0.5 text-xs opacity-70">{sortIcon(col.key)}</span>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {sorted.map((player, index) => (
                <tr
                  key={player.id}
                  className="border-b border-slate-100 last:border-0 hover:bg-slate-50"
                >
                  <td className="px-4 py-3 text-slate-500">{index + 1}</td>
                  <td className="px-4 py-3 font-medium text-slate-900">{player.name}</td>
                  <td className="px-4 py-3">{player.played}</td>
                  <td className="px-4 py-3">{player.wins}</td>
                  <td className="px-4 py-3">{player.losses}</td>
                  <td className="px-4 py-3">{player.otl}</td>
                  <td className="px-4 py-3">{player.points}</td>
                  <td className="px-4 py-3">{(player.win_rate * 100).toFixed(1)}%</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        </div>
      </div>
    </div>
  )
}
