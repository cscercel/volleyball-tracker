<script lang="ts">
import { login } from '../lib/api/index.ts'
import { auth } from '../lib/auth.svelte.ts'

let email = $state('')
let password = $state('')
    let error = $state('')
let loading = $state(false)

    async function handleLogin() {
        error = ''
            loading = true
            try {
                const token = await login(email, password)
                    auth.login(token)
            } catch {
                error = 'Invalid email or password.'
            } finally {
                loading = false
            }
    }
</script>

<div class="page">
<div class="login-card">
<h1>🔐 Admin </h1>

{#if auth.isAuthenticated}
<p class="success">✅ Logged in</p>
<button onclick={() => auth.logout()}>Logout</button>
{:else}
<div class="form">
<input
type="email"
placeholder="Email"
bind:value={email}
/>
<input
type="password"
placeholder="Password"
bind:value={password}
/>
<button onclick={handleLogin} disabled={loading || !email || !password}>
{loading ? 'Logging in...' : 'Login'}
</button>
</div>
{#if error}<p class="error">{error}</p>{/if}
{/if}
</div>
</div>

<style>
.page {
padding: 2rem;
display: flex;
         justify-content: center;
}

.login-card {
width: 100%;
       max-width: 400px;
padding: 2rem;
border: 1px solid #eee;
        border-radius: 12px;
background: #1a1a2e;
}

.form {
display: flex;
         flex-direction: column;
gap: 1rem;
}

input {
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

.success { color: green; }
.error { color: red; }
</style>
