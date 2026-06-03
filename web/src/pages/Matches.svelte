<script lang="ts">
import { auth } from '../lib/auth.svelte.ts'
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
} from '../lib/api/index.ts'
import Court from '../lib/components/Court.svelte'

const currentYear = new Date().getFullYear()

    let activeTab = $state('create')

    let players: Player[] = $state([])

    async function fetchPlayers() {
        try {
            players = await getPlayers()
        } catch (e) {
            console.error(e)
        }
    }

$effect(() => { fetchPlayers() })

    // ---------------------------------------------------------------------------
    // TAB 1 - Create Match
    // ---------------------------------------------------------------------------
    let matchType = $state('indoor')
    let draftType = $state('random')
    let blueTeam: string[] = $state([])
    let redTeam: string[] = $state([])
    let createError = $state('')
    let createSuccess = $state('')
    let createdMatchRoster: MatchPlayer[] = $state([])

    let blueSelected = $state('')
    let redSelected = $state('')

    const availableForBlue = $derived(
            players.filter(p => !blueTeam.includes(p.name) && !redTeam.includes(p.name))
            )

    const availableForRed = $derived(
            players.filter(p => !redTeam.includes(p.name) && !blueTeam.includes(p.name))
            )

    function addToBlue() {
        if (blueSelected && !blueTeam.includes(blueSelected)) {
            blueTeam = [...blueTeam, blueSelected]
                blueSelected = ''
        }
    }

function addToRed() {
    if (redSelected && !redTeam.includes(redSelected)) {
        redTeam = [...redTeam, redSelected]
            redSelected = ''
    }
}

function removeFromBlue(name: string) {
    blueTeam = blueTeam.filter(p => p !== name)
}

function removeFromRed(name: string) {
    redTeam = redTeam.filter(p => p !== name)
}

async function handleCreateMatch() {
    createError = ''
        createSuccess = ''
        createdMatchRoster = []

        if (blueTeam.length === 0 || redTeam.length === 0) {
            createError = 'Both teams need at least one player.'
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
            createdMatchRoster = await getMatchRoster(match.id)
            createSuccess = '✅ Match created!'
            blueTeam = []
            redTeam = []
    } catch {
        createError = 'Failed to create match.'
    }
}

// ---------------------------------------------------------------------------
// TAB 2 - Draft Matches
// ---------------------------------------------------------------------------
let drafts: Match[] = $state([])
let draftRosters: Record<string, MatchPlayer[]> = $state({})
let loadingDrafts = $state(false)
    let blueScores: Record<string, number> = $state({})
    let redScores: Record<string, number> = $state({})

    async function fetchDrafts() {
        loadingDrafts = true
            try {
                drafts = await getUncompletedMatches()
                    for (const match of drafts) {
                        draftRosters[match.id] = await getMatchRoster(match.id)
                    }
            } catch {
                console.error('Failed to load drafts')
            } finally {
                loadingDrafts = false
            }
    }

async function handleSubmitResults(matchId: string) {
    const blue = blueScores[matchId]
        const red = redScores[matchId]
        if (blue === undefined || red === undefined) {
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

// ---------------------------------------------------------------------------
// TAB 3 - Completed Matches
// ---------------------------------------------------------------------------
    let completedType = $state('indoor')
let completedSeason = $state(currentYear)
    let completed: Match[] = $state([])
    let completedRosters: Record<string, MatchPlayer[]> = $state({})
let loadingCompleted = $state(false)

    async function fetchCompleted() {
        loadingCompleted = true
            try {
                completed = await getMatchesBySeason(completedType, completedSeason)
                    for (const match of completed) {
                        completedRosters[match.id] = await getMatchRoster(match.id)
                    }
            } catch {
                console.error('Failed to load completed matches')
            } finally {
                loadingCompleted = false
            }
    }

$effect(() => {
        if (activeTab === 'drafts') fetchDrafts()
        if (activeTab === 'completed') fetchCompleted()
        })

$effect(() => {
        completedType
        completedSeason
        if (activeTab === 'completed') fetchCompleted()
        })
</script>

<div class="page">
<h1>🏐 Matches</h1>

<div class="tabs">
<button class:active={activeTab === 'create'} onclick={() => activeTab = 'create'}>Create Match</button>
<button class:active={activeTab === 'drafts'} onclick={() => activeTab = 'drafts'}>Draft Matches</button>
<button class:active={activeTab === 'completed'} onclick={() => activeTab = 'completed'}>Completed Matches</button>
</div>

<!-- TAB 1 - Create Match -->
{#if activeTab === 'create'}
<div class="tab-content">
<h2>Create New Match</h2>
{#if players.length < 2}
<p>⚠️ Need at least 2 players to create a match.</p>
{:else if !auth.isAuthenticated}
<p>You must be <a href="/login">logged in</a> to create matches.</p>
{:else}
<div class="filters">
<div class="pill-group">
<span>Match Type</span>
<button class:pill-active={matchType === 'indoor'} onclick={() => matchType = 'indoor'}>Indoor</button>
<button class:pill-active={matchType === 'beach'} onclick={() => matchType = 'beach'}>Beach</button>
</div>
<div class="pill-group">
<span>Draft Type</span>
<button class:pill-active={draftType === 'random'} onclick={() => draftType = 'random'}>Random</button>
<button class:pill-active={draftType === 'manual'} onclick={() => draftType = 'manual'}>Manual</button>
</div>
</div>

<div class="teams">
<!-- Blue Team -->
<div class="team blue-team">
<h3>🔵 Blue Team</h3>

<div class="team-players">
{#if blueTeam.length === 0}
<p class="empty">No players added yet</p>
{:else}
{#each blueTeam as name}
<div class="player-tag">
<span>{name}</span>
<button class="remove" onclick={() => removeFromBlue(name)}>✕</button>
</div>
{/each}
{/if}
</div>

<select bind:value={blueSelected} onchange={addToBlue}>
<option value="">Add a player...</option>
{#each availableForBlue as player}
<option value={player.name}>{player.name}</option>
{/each}
</select>
</div>

<!-- Red Team -->
<div class="team red-team">
<h3>🔴 Red Team</h3>

<div class="team-players">
{#if redTeam.length === 0}
<p class="empty">No players added yet</p>
{:else}
{#each redTeam as name}
<div class="player-tag">
<span>{name}</span>
<button class="remove" onclick={() => removeFromRed(name)}>✕</button>
</div>
{/each}
{/if}
</div>

<select bind:value={redSelected} onchange={addToRed}>
<option value="">Add a player...</option>
{#each availableForRed as player}
<option value={player.name}>{player.name}</option>
{/each}
</select>
</div>
</div>

<button onclick={handleCreateMatch}>Create Match</button>

{#if createSuccess}<p class="success">{createSuccess}</p>{/if}
{#if createError}<p class="error">{createError}</p>{/if}

{#if createdMatchRoster.length > 0}
{@const bluePlayers = createdMatchRoster.filter(p => p.color === 'blue').map(p => p.player_name)}
{@const redPlayers = createdMatchRoster.filter(p => p.color === 'red').map(p => p.player_name)}
<div class="teams">
<div class="team result blue-team">
<h3>🔵 Blue Team</h3>
{#each bluePlayers as p}<p>{p}</p>{/each}
</div>
<div class="team result red-team">
<h3>🔴 Red Team</h3>
{#each redPlayers as p}<p>{p}</p>{/each}
</div>
</div>
<Court blueTeam={bluePlayers} redTeam={redPlayers} />
{/if}
{/if}
</div>

<!-- TAB 2 - Draft Matches -->
{:else if activeTab === 'drafts'}
<div class="tab-content">
<h2>Draft Matches</h2>
{#if loadingDrafts}
<p>Loading...</p>
{:else if drafts.length === 0}
<p>No draft matches. Create one in the Create Match tab!</p>
{:else}
{#each drafts as match}
<div class="match-card">
<div class="match-header">
<span>Match {match.id.slice(0, 8)}</span>
<span>{match.match_type.toUpperCase()}</span>
</div>

<div class="teams">
<div class="team blue-team">
<h3>🔵 Blue Team</h3>
{#each (draftRosters[match.id] ?? []).filter(p => p.color === 'blue') as p}
<p>{p.player_name}</p>
{/each}
</div>
<div class="team red-team">
<h3>🔴 Red Team</h3>
{#each (draftRosters[match.id] ?? []).filter(p => p.color === 'red') as p}
<p>{p.player_name}</p>
{/each}
</div>
</div>

{#if auth.isAuthenticated}
<div class="submit-results">
<input type="number" placeholder="Blue score" bind:value={blueScores[match.id]} />
<input type="number" placeholder="Red score" bind:value={redScores[match.id]} />
<button onclick={() => handleSubmitResults(match.id)}>Submit Results</button>
<button class="danger" onclick={() => handleDeleteDraft(match.id)}>Delete Draft</button>
</div>
{/if}
</div>
{/each}
{/if}
</div>

<!-- TAB 3 - Completed Matches -->
{:else if activeTab === 'completed'}
<div class="tab-content">
<h2>Completed Matches</h2>

<div class="filters">
<select bind:value={completedType}>
<option value="indoor">Indoor</option>
<option value="beach">Beach</option>
</select>
<input type="number" min="2023" max={currentYear} bind:value={completedSeason} />
</div>

{#if loadingCompleted}
<p>Loading...</p>
{:else if completed.length === 0}
<p>No completed matches for this filter.</p>
{:else}
{#each completed as match}
{@const isOtl = Math.abs(match.blue_score - match.red_score) === 2}
{@const winner = match.blue_score > match.red_score ? 'blue' : 'red'}
<div class="match-card">
<div class="match-header">
<span>{winner === 'blue' ? '🔵' : '🔴'} {match.blue_score} – {match.red_score} {isOtl ? '⏱️ OT' : ''}</span>
<span>{match.match_type.toUpperCase()}</span>
<span>{match.created_at.slice(0, 10)}</span>
</div>

<div class="teams">
<div class="team blue-team">
<h3>🔵 Blue — {match.blue_score}</h3>
{#each (completedRosters[match.id] ?? []).filter(p => p.color === 'blue') as p}
<p>{p.player_name}</p>
{/each}
</div>
<div class="team red-team">
<h3>🔴 Red — {match.red_score}</h3>
{#each (completedRosters[match.id] ?? []).filter(p => p.color === 'red') as p}
<p>{p.player_name}</p>
{/each}
</div>
</div>
</div>
{/each}
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
     align-items: center;
}

.pill-group {
display: flex;
         align-items: center;
gap: 0.5rem;
}

.pill-group button {
padding: 0.35rem 0.85rem;
         border-radius: 999px;
border: 1px solid #ccc;
background: white;
cursor: pointer;
color: #444;
}

.pill-group button.pill-active {
background: #ff6b35;
color: white;
       border-color: #ff6b35;
}

.teams {
display: grid;
         grid-template-columns: 1fr 1fr;
gap: 1rem;
     margin-bottom: 1rem;
}

.team {
padding: 1rem;
         border-radius: 8px;
border: 1px solid #eee;
}

.team h3 {
margin: 0 0 0.75rem 0;
color: white;
}

.blue-team {
background: #0d2b5e;
            border-color: #0d2b5e;
}

.red-team {
background: #6b0d0d;
            border-color: #6b0d0d;
}

.blue-team p,
    .red-team p {
color: #ddd;
margin: 0.25rem 0;
    }

.team-players {
    min-height: 60px;
    margin-bottom: 0.75rem;
display: flex;
         flex-wrap: wrap;
gap: 0.5rem;
}

.empty {
color: rgba(255,255,255,0.4) !important;
       font-style: italic;
       font-size: 0.9rem;
}

.player-tag {
display: flex;
         align-items: center;
gap: 0.4rem;
padding: 0.3rem 0.7rem;
         border-radius: 999px;
         font-size: 0.9rem;
color: white;
background: rgba(255,255,255,0.15);
}

.remove {
background: none;
border: none;
color: rgba(255,255,255,0.6);
cursor: pointer;
padding: 0;
         font-size: 0.8rem;
         line-height: 1;
}

.remove:hover { color: white; }

.team select {
width: 100%;
padding: 0.4rem;
         border-radius: 4px;
border: none;
        font-size: 0.9rem;
cursor: pointer;
}

.match-card {
border: 1px solid #eee;
        border-radius: 8px;
padding: 1rem;
         margin-bottom: 1rem;
}

.match-header {
display: flex;
gap: 1rem;
     font-weight: 600;
     margin-bottom: 0.75rem;
color: #444;
}

.submit-results {
display: flex;
gap: 0.75rem;
     align-items: center;
     margin-top: 0.75rem;
}

select, input[type="number"] {
padding: 0.5rem;
         font-size: 1rem;
         border-radius: 4px;
border: 1px solid #ccc;
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
