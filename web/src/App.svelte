<script lang="ts">
  import Home from './pages/Home.svelte'
  import Players from './pages/Players.svelte'
  import Matches from './pages/Matches.svelte'

  const routes: Record<string, any> = {
    '/': Home,
    '/players': Players,
    '/matches': Matches,
  }

  let currentPath = $state(window.location.pathname)

  window.addEventListener('popstate', () => {
    currentPath = window.location.pathname
  })

  function navigate(path: string) {
    window.history.pushState({}, '', path)
    currentPath = path
  }

  const currentComponent = $derived(routes[currentPath] ?? Home)
</script>

<nav>
  <button onclick={() => navigate('/')}>🏐 Volleyball Tracker</button>
  <button onclick={() => navigate('/players')}>Players</button>
  <button onclick={() => navigate('/matches')}>Matches</button>
</nav>

<main>
  <svelte:component this={currentComponent} />
</main>
