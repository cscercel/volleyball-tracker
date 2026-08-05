import { NavLink } from 'react-router-dom'
import { useAuth } from '../lib/AuthContext'

const navItems = [
  { to: '/', label: 'Leaderboard', icon: TrophyIcon, end: true },
  { to: '/players', label: 'Players', icon: UsersIcon, end: false },
  { to: '/matches', label: 'Matches', icon: VolleyballIcon, end: false },
]

export default function Sidebar() {
  const auth = useAuth()

  return (
    <aside className="hidden h-screen w-64 shrink-0 flex-col bg-brand-bg text-slate-300 md:flex">
      <div className="flex items-center gap-2 px-6 py-6">
        <span className="text-2xl">🏐</span>
        <span className="text-lg font-semibold tracking-tight text-white">
          Volleyball Tracker
        </span>
      </div>

      <nav className="flex-1 space-y-1 px-3">
        {navItems.map(({ to, label, icon: Icon, end }) => (
          <NavLink
            key={to}
            to={to}
            end={end}
            className={({ isActive }) =>
              [
                'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-brand-accent text-white'
                  : 'text-slate-300 hover:bg-brand-bg-hover hover:text-white',
              ].join(' ')
            }
          >
            <Icon className="h-5 w-5 shrink-0" />
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="border-t border-slate-700/60 px-3 py-4">
        {auth.isAuthenticated ? (
          <div className="space-y-2 px-3">
            <p className="flex items-center gap-2 text-xs font-medium text-emerald-400">
              <span className="h-1.5 w-1.5 rounded-full bg-emerald-400" />
              Admin
            </p>
            <button
              onClick={() => auth.logout()}
              className="w-full rounded-lg border border-slate-600 px-3 py-2 text-left text-sm text-slate-300 transition-colors hover:bg-brand-bg-hover hover:text-white"
            >
              Logout
            </button>
          </div>
        ) : (
          <NavLink
            to="/login"
            className="block rounded-lg border border-brand-accent px-3 py-2 text-center text-sm font-medium text-brand-accent transition-colors hover:bg-brand-accent hover:text-white"
          >
            Login
          </NavLink>
        )}
      </div>
    </aside>
  )
}

function TrophyIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" {...props}>
      <path d="M8 21h8M12 17v4M7 4h10v4a5 5 0 0 1-10 0V4Z" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M7 5H4a1 1 0 0 0-1 1v1a4 4 0 0 0 4 4M17 5h3a1 1 0 0 1 1 1v1a4 4 0 0 1-4 4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function UsersIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" {...props}>
      <path d="M17 21v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2M10 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8ZM21 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function VolleyballIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" {...props}>
      <circle cx="12" cy="12" r="9" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M12 3c3 2.5 3 15.5 0 18M4.5 8c3.5 1.5 12 1.5 15 0M4.5 16c3.5-1.5 12-1.5 15 0" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}
