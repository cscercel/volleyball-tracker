<script lang="ts">
import { getLeaderboard, type PlayerStats } from '../lib/api/index.ts'

const currentYear = new Date().getFullYear()

    let matchType = $state('indoor')
let season = $state(currentYear)
    let leaderboard: PlayerStats[] = $state([])
let loading = $state(false)
    let error = $state('')

    async function fetchLeaderboard() {
        loading = true
            error = ''
            try {
                leaderboard = await getLeaderboard(matchType, season)
            } catch (e) {
                error = 'Failed to load leaderboard.'
            } finally {
                loading = false
            }
    }

// Fetch whenever matchType or season changes
$effect(() => {
        matchType
        season
        fetchLeaderboard()
        })

function calculateMmr(avgPoints: number, efficiencyRate: number): number {
    return avgPoints * efficiencyRate
}

function getRank(played: number, points: number, efficiencyRate: number): string {
    if (played < 10) return 'Unranked'
        const avgPoints = played > 0 ? points / played : 0
            const mmr = calculateMmr(avgPoints, efficiencyRate)
            const ranks: [string, number, number][] = [
            ['Iron I',       0,    0.1],
            ['Iron II',      0.1,  0.2],
            ['Iron III',     0.2,  0.3],
            ['Bronze I',     0.3,  0.4],
            ['Bronze II',    0.4,  0.5],
            ['Bronze III',   0.5,  0.6],
            ['Silver I',     0.6,  0.7],
            ['Silver II',    0.7,  0.8],
            ['Silver III',   0.8,  0.9],
            ['Gold I',       0.9,  1.0],
            ['Gold II',      1.0,  1.1],
            ['Gold III',     1.1,  1.2],
            ['Platinum I',   1.2,  1.3],
            ['Platinum II',  1.3,  1.4],
            ['Platinum III', 1.4,  1.5],
            ['Diamond I',    1.5,  1.6],
            ['Diamond II',   1.6,  1.7],
            ['Diamond III',  1.7,  1.8],
            ['Spiker',       1.8,  1.9],
            ['Ace',          1.9,  2.0],
            ['Sensei',       2.0,  Infinity],
            ]
                for (const [name, low, high] of ranks) {
                    if (mmr >= low && mmr < high) return name
                }
            return 'Iron I'
}
</script>

<div class="page">
<h1>🏆 Leaderboard</h1>

<div class="filters">
<select bind:value={matchType}>
<option value="indoor">Indoor</option>
<option value="beach">Beach</option>
</select>

<input
type="number"
min="2023"
max={currentYear}
bind:value={season}
/>
</div>

{#if loading}
<p>Loading...</p>
{:else if error}
<p class="error">{error}</p>
{:else if leaderboard.length === 0}
<p>No stats for {matchType} in season {season}.</p>
{:else}
<table>
<thead>
<tr>
<th>#</th>
<th>Player</th>
<th>Played</th>
<th>Wins</th>
<th>Losses</th>
<th>OTL</th>
<th>Points</th>
<th>Win Rate</th>
<th>Rank</th>
</tr>
</thead>
<tbody>
{#each leaderboard as player, index}
<tr>
<td>{index + 1}</td>
<td>{player.name}</td>
<td>{player.played}</td>
<td>{player.wins}</td>
<td>{player.losses}</td>
<td>{player.otl}</td>
<td>{player.points}</td>
<td>{(player.win_rate * 100).toFixed(1)}%</td>
<td>{getRank(player.played, player.points, player.efficiency_rate)}</td>
</tr>
{/each}
</tbody>
</table>
{/if}
</div>

<style>
.page {
padding: 2rem;
}

.filters {
display: flex;
gap: 1rem;
     margin-bottom: 1.5rem;
}

select, input {
padding: 0.5rem;
         font-size: 1rem;
         border-radius: 4px;
border: 1px solid #ccc;
}

table {
width: 100%;
       border-collapse: collapse;
}

th, td {
    text-align: left;
padding: 0.75rem 1rem;
         border-bottom: 1px solid #eee;
}

th {
    font-weight: 600;
background: #1a1a2e;
}

tr:hover {
background: #f5f5f5;
   }

.error {
color: red;
}
</style>
