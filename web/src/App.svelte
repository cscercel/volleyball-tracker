<script lang="ts">
import Home from './pages/Home.svelte'
import Players from './pages/Players.svelte'
import Matches from './pages/Matches.svelte'
import Login from './pages/Login.svelte'
import { auth } from './lib/auth.svelte.ts'

const routes: Record<string, any> = {
    '/': Home,
    '/players': Players,
    '/matches': Matches,
    '/login': Login,
}

let currentPath = $state(window.location.pathname)

    window.addEventListener('popstate', () => {
            currentPath = window.location.pathname
            })

function navigate(path: string) {
    window.history.pushState({}, '', path)
        currentPath = path
}
</script>

<nav>
<div class="nav-brand" onclick={() => navigate('/')}>
🏐 Volleyball Tracker
</div>

<div class="nav-links">
<button
class:active={currentPath === '/'}
onclick={() => navigate('/')}
>Leaderboard</button>
<button
class:active={currentPath === '/players'}
onclick={() => navigate('/players')}
>Players</button>
<button
class:active={currentPath === '/matches'}
onclick={() => navigate('/matches')}
>Matches</button>
</div>

<div class="nav-auth">
{#if auth.isAuthenticated}
<span class="nav-status">✅ Admin</span>
<button class="nav-logout" onclick={() => auth.logout()}>Logout</button>
{:else}
<button class="nav-login" onclick={() => navigate('/login')}>Login</button>
{/if}
</div>
</nav>

<main>
{#if currentPath === '/'}
<Home />
{:else if currentPath === '/players'}
<Players />
{:else if currentPath === '/matches'}
<Matches />
{:else if currentPath === '/login'}
<Login />
{:else}
<Home />
{/if}
</main>

<style>
nav {
display: flex;
         align-items: center;
         justify-content: space-between;
padding: 0 2rem;
height: 60px;
background: #1a1a2e;
color: white;
position: sticky;
top: 0;
     z-index: 100;
}

.nav-brand {
    font-size: 1.2rem;
    font-weight: 700;
cursor: pointer;
        letter-spacing: 0.5px;
}

.nav-links {
display: flex;
gap: 0.5rem;
}

.nav-links button {
background: none;
border: none;
color: #aaa;
       font-size: 0.95rem;
padding: 0.4rem 0.85rem;
         border-radius: 6px;
cursor: pointer;
transition: color 0.2s, background 0.2s;
}

.nav-links button:hover {
color: white;
background: rgba(255,255,255,0.08);
}

.nav-links button.active {
color: #ff6b35;
       font-weight: 600;
}

.nav-auth {
display: flex;
         align-items: center;
gap: 0.75rem;
}

.nav-status {
    font-size: 0.85rem;
color: #aaa;
}

.nav-login, .nav-logout {
padding: 0.4rem 1rem;
         border-radius: 6px;
border: 1px solid #ff6b35;
background: transparent;
color: #ff6b35;
cursor: pointer;
        font-size: 0.9rem;
transition: background 0.2s, color 0.2s;
}

.nav-login:hover, .nav-logout:hover {
background: #ff6b35;
color: white;
}

main {
    min-height: calc(100vh - 60px);
}
</style>
