import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../lib/api'
import { useAuth } from '../lib/AuthContext'

export default function Login() {
  const auth = useAuth()
  const navigate = useNavigate()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleLogin() {
    setError('')
    setLoading(true)
    try {
      const token = await login(email, password)
      auth.login(token)
      navigate('/')
    } catch {
      setError('Invalid email or password.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex justify-center">
      <div className="w-full max-w-sm rounded-2xl bg-brand-bg p-8 text-white shadow-sm">
        <h1 className="mb-6 text-xl font-semibold">🔐 Admin</h1>

        {auth.isAuthenticated ? (
          <>
            <p className="mb-4 text-sm text-emerald-400">✅ Logged in</p>
            <button
              onClick={() => auth.logout()}
              className="w-full rounded-lg border border-slate-600 px-4 py-2 text-sm font-medium transition-colors hover:bg-brand-bg-hover"
            >
              Logout
            </button>
          </>
        ) : (
          <>
            <div className="flex flex-col gap-3">
              <input
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && !loading && email && password && handleLogin()}
                className="rounded-lg border border-slate-600 bg-white/5 px-3 py-2 text-sm text-white placeholder:text-slate-400 focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
              />
              <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && !loading && email && password && handleLogin()}
                className="rounded-lg border border-slate-600 bg-white/5 px-3 py-2 text-sm text-white placeholder:text-slate-400 focus:border-brand-accent focus:outline-none focus:ring-2 focus:ring-brand-accent/30"
              />
              <button
                onClick={handleLogin}
                disabled={loading || !email || !password}
                className="rounded-lg bg-brand-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-accent-hover disabled:cursor-not-allowed disabled:opacity-50"
              >
                {loading ? 'Logging in...' : 'Login'}
              </button>
            </div>
            {error && <p className="mt-3 text-sm text-red-400">{error}</p>}
          </>
        )}
      </div>
    </div>
  )
}
