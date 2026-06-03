<script lang="ts">
import { auth } from '../lib/auth.svelte.ts'
import { getPlayers, getPlayerStats, getPlayerHistory, createPlayer, updatePlayerName, deletePlayer, type Player, type PlayerStats, type MatchHistory } from '../lib/api/index.ts'

const currentYear = new Date().getFullYear()

    let activeTab = $state('profile')

    let players: Player[] = $state([])
let loadingPlayers = $state(false)

    async function fetchPlayers() {
        loadingPlayers = true
            try {
                players = await getPlayers()
            } catch (e) {
                console.error(e)
            } finally {
                loadingPlayers = false
            }
    }

let selectedPlayerId = $state('')
    let matchType = $state('indoor')
let season = $state(currentYear)
    let stats: PlayerStats | null = $state(null)
    let prevStats: PlayerStats | null = $state(null)
    let history: MatchHistory[] = $state([])
let loadingStats = $state(false)
    let statsError = $state('')

    async function fetchStats() {
        if (!selectedPlayerId) return
            loadingStats = true
                statsError = ''
                try {
                    stats = await getPlayerStats(selectedPlayerId, matchType, season)
                        history = await getPlayerHistory(selectedPlayerId, matchType, season)
                        try {
                            prevStats = await getPlayerStats(selectedPlayerId, matchType, season - 1)
                        } catch {
                            prevStats = null
                        }
                } catch (e) {
                    statsError = 'Failed to load player stats.'
                        stats = null
                } finally {
                    loadingStats = false
                }
    }

$effect(() => { fetchPlayers() })

    $effect(() => {
            selectedPlayerId
            matchType
            season
            fetchStats()
            })

let newPlayerName = $state('')
let addSuccess = $state('')
let addError = $state('')

async function handleAddPlayer() {
    addSuccess = ''
        addError = ''
        try {
            await createPlayer(newPlayerName)
                addSuccess = `✅ Added ${newPlayerName}`
                newPlayerName = ''
                await fetchPlayers()
        } catch {
            addError = 'Failed to add player.'
        }
}

let managePlayerId = $state('')
let newName = $state('')
let manageSuccess = $state('')
let manageError = $state('')

async function handleRename() {
    manageSuccess = ''
        manageError = ''
        try {
            await updatePlayerName(managePlayerId, newName)
                manageSuccess = '✅ Renamed successfully'
                await fetchPlayers()
        } catch {
            manageError = 'Failed to rename player.'
        }
}

async function handleDelete() {
    manageSuccess = ''
        manageError = ''
        try {
            await deletePlayer(managePlayerId)
                manageSuccess = '✅ Player deleted'
                managePlayerId = ''
                await fetchPlayers()
        } catch {
            manageError = 'Failed to delete player.'
        }
}

function delta(current: number, previous: number | null): string {
    if (previous === null) return ''
        const diff = current - previous
            if (diff === 0) return ''
                return diff > 0 ? `(+${diff})` : `(${diff})`
}

function calculateMmr(avgPoints: number, efficiencyRate: number): number {
    return avgPoints * efficiencyRate
}

function getRank(played: number, points: number, efficiencyRate: number): string {
    if (played < 10) return 'Unranked'

        const avgPoints = played > 0 ? points / played : 0
            const mmr = calculateMmr(avgPoints, efficiencyRate)

            const ranks: [string, number, number][] = [
            ['Iron I',      0,    0.1],
            ['Iron II',     0.1,  0.2],
            ['Iron III',    0.2,  0.3],
            ['Bronze I',    0.3,  0.4],
            ['Bronze II',   0.4,  0.5],
            ['Bronze III',  0.5,  0.6],
            ['Silver I',    0.6,  0.7],
            ['Silver II',   0.7,  0.8],
            ['Silver III',  0.8,  0.9],
            ['Gold I',      0.9,  1.0],
            ['Gold II',     1.0,  1.1],
            ['Gold III',    1.1,  1.2],
            ['Platinum I',  1.2,  1.3],
            ['Platinum II', 1.3,  1.4],
            ['Platinum III',1.4,  1.5],
            ['Diamond I',   1.5,  1.6],
            ['Diamond II',  1.6,  1.7],
            ['Diamond III', 1.7,  1.8],
            ['Spiker',      1.8,  1.9],
            ['Ace',         1.9,  2.0],
            ['Sensei',      2.0,  Infinity],
            ]

                for (const [name, low, high] of ranks) {
                    if (mmr >= low && mmr < high) return name
                }

            return 'Iron I'
}
</script>

<div class="page">
<h1>👥 Players</h1>

<div class="tabs">
<button class:active={activeTab === 'profile'} onclick={() => activeTab = 'profile'}>Player Profile</button>
<button class:active={activeTab === 'add'} onclick={() => activeTab = 'add'}>Add Player</button>
<button class:active={activeTab === 'manage'} onclick={() => activeTab = 'manage'}>Manage Players</button>
</div>

<!-- Profile Tab -->
{#if activeTab === 'profile'}
<div class="tab-content">
{#if loadingPlayers}
<p>Loading players...</p>
{:else if players.length === 0}
<p>No players yet!</p>
{:else}
<div class="filters">
<select bind:value={selectedPlayerId}>
<option value="">Select a player</option>
{#each players as player}
<option value={player.id}>{player.name}</option>
{/each}
</select>
<select bind:value={matchType}>
<option value="indoor">Indoor</option>
<option value="beach">Beach</option>
</select>
<input type="number" min="2023" max={currentYear} bind:value={season} />
</div>

{#if !selectedPlayerId}
<p>Select a player to view their stats.</p>
{:else if loadingStats}
<p>Loading stats...</p>
{:else if statsError}
<p class="error">{statsError}</p>
{:else if stats}
<div class="stats-grid">
<div class="stat-card">
<span class="stat-label">Matches Played</span>
<span class="stat-value">{stats.played} {delta(stats.played, prevStats?.played ?? null)}</span>
</div>
<div class="stat-card">
<span class="stat-label">Wins</span>
<span class="stat-value">{stats.wins} {delta(stats.wins, prevStats?.wins ?? null)}</span>
</div>
<div class="stat-card">
<span class="stat-label">Losses</span>
<span class="stat-value">{stats.losses}</span>
</div>
<div class="stat-card">
<span class="stat-label">OTL</span>
<span class="stat-value">{stats.otl}</span>
</div>
<div class="stat-card">
<span class="stat-label">Points</span>
<span class="stat-value">{stats.points} {delta(stats.points, prevStats?.points ?? null)}</span>
</div>
<div class="stat-card">
<span class="stat-label">Win Rate</span>
<span class="stat-value">
{(stats.win_rate * 100).toFixed(1)}%
{delta(Math.round(stats.win_rate * 100), prevStats ? Math.round(prevStats.win_rate * 100) : null)}
</span>
</div>
<div class="stat-card">
<span class="stat-label">Win Streak</span>
<span class="stat-value">{stats.streak}</span>
</div>
<div class="stat-card">
<span class="stat-label">Longest Streak</span>
<span class="stat-value">{stats.longest_streak} {delta(stats.longest_streak, prevStats?.longest_streak ?? null)}</span>
</div>
</div>

<div class="rank-section">
<h3>Rank: {getRank(stats.played, stats.points, stats.efficiency_rate)}</h3>
<img
src={`/assets/${getRank(stats.played, stats.points, stats.efficiency_rate)}.png`}
alt={getRank(stats.played, stats.points, stats.efficiency_rate)}
width="160"
/>
</div>

<div class="history">
<h3>Match History</h3>
{#if history.length === 0}
<p>No completed matches this season.</p>
{:else}
{#each history as match}
{@const myScore = match.color === 'blue' ? match.blue_score : match.red_score}
{@const theirScore = match.color === 'blue' ? match.red_score : match.blue_score}
{@const won = myScore > theirScore}
{@const isOtl = Math.abs(match.blue_score - match.red_score) === 2 && !won}
<div class="history-row" class:win={won} class:otl={isOtl} class:loss={!won && !isOtl}>
<span>{won ? '✅ Win' : isOtl ? '⚠️ OTL' : '❌ Loss'}</span>
<span>{myScore} : {theirScore}</span>
<span>{match.color} team</span>
<span>{match.created_at.slice(0, 10)}</span>
</div>
{/each}
{/if}
</div>
{/if}
{/if}
</div>

<!-- Add Player Tab -->
{:else if activeTab === 'add'}
<div class="tab-content">
<h2>Add New Player</h2>
{#if auth.isAuthenticated}
<div class="form">
<input type="text" placeholder="Player name" bind:value={newPlayerName} />
<button onclick={handleAddPlayer} disabled={!newPlayerName}>Add Player</button>
</div>
{#if addSuccess}<p class="success">{addSuccess}</p>{/if}
{#if addError}<p class="error">{addError}</p>{/if}
{:else}
<p>You must be <a href="/login">logged in</a> to add players.</p>
{/if}
</div>

<!-- Manage Players Tab -->
{:else if activeTab === 'manage'}
<div class="tab-content">
<h2>Manage Players</h2>
{#if !auth.isAuthenticated}
<p>You must be <a href="/login">logged in</a> to manage players.</p>
{:else if players.length === 0}
<p>No players to manage.</p>
{:else}
<select bind:value={managePlayerId}>
<option value="">Select a player</option>
{#each players as player}
<option value={player.id}>{player.name}</option>
{/each}
</select>

{#if managePlayerId}
<div class="manage-actions">
<div class="form">
<h3>✏️ Rename</h3>
<input type="text" placeholder="New name" bind:value={newName} />
<button onclick={handleRename} disabled={!newName}>Rename</button>
</div>

<div class="danger-zone">
<h3>🗑️ Delete</h3>
<p>This will permanently delete this player.</p>
<button class="danger" onclick={handleDelete}>Confirm Delete</button>
</div>
</div>
{/if}

{#if manageSuccess}<p class="success">{manageSuccess}</p>{/if}
{#if manageError}<p class="error">{manageError}</p>{/if}
{/if}
</div>
{/if}
</div>

<style>
.page { padding: 2rem; }

.tabs {
display: flex;
gap: 0.5rem;
     margin-bottom: 1.5rem;
     border-bottom: 2px solid #eee;
     padding-bottom: 0.5rem;
}

.tabs button {
padding: 0.5rem 1rem;
border: none;
background: none;
cursor: pointer;
        font-size: 1rem;
color: #666;
       border-radius: 4px 4px 0 0;
}

.tabs button.active {
color: #ff6b35;
       font-weight: 600;
       border-bottom: 2px solid #ff6b35;
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

.stats-grid {
display: grid;
         grid-template-columns: repeat(4, 1fr);
gap: 1rem;
     margin-bottom: 2rem;
}

.stat-card {
display: flex;
         flex-direction: column;
padding: 1rem;
border: 1px solid #eee;
        border-radius: 8px;
background: #fafafa;
}

.stat-label {
    font-size: 0.85rem;
color: #666;
       margin-bottom: 0.25rem;
}

.stat-value {
    font-size: 1.5rem;
    font-weight: 600;
}

.rank-section {
display: flex;
         align-items: center;
gap: 1rem;
     margin-bottom: 2rem;
}

.history { margin-top: 1.5rem; }

.history-row {
display: flex;
gap: 2rem;
padding: 0.5rem 1rem;
         border-radius: 4px;
         margin-bottom: 0.5rem;
}

.win { background: #e8f5e9; }
.otl { background: #fff8e1; }
.loss { background: #fce4ec; }

.form {
display: flex;
gap: 1rem;
     align-items: center;
     margin-bottom: 1rem;
}

.manage-actions {
display: flex;
         flex-direction: column;
gap: 2rem;
     margin-top: 1.5rem;
}

.danger-zone {
border: 1px solid #ffcdd2;
        border-radius: 8px;
padding: 1rem;
background: #fff5f5;
}

button {
padding: 0.5rem 1rem;
         border-radius: 4px;
border: none;
cursor: pointer;
background: #ff6b35;
color: white;
       font-size: 1rem;
}

button:disabled {
opacity: 0.5;
cursor: not-allowed;
       }

button.danger { background: #e53935; }

.success { color: green; }
.error { color: red; }
</style>
