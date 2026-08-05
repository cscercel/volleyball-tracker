import { NavLink } from 'react-router-dom'
import { useAuth } from '../lib/AuthContext'

const navItems = [
  { to: '/', label: 'Board', icon: TrophyIcon, end: true },
  { to: '/players', label: 'Players', icon: UsersIcon, end: false },
  { to: '/matches', label: 'Matches', icon: VolleyballIcon, end: false },
]

export default function BottomNav() {
  const auth = useAuth()

  const authItem = auth.isAuthenticated
    ? { to: '/login', label: 'Admin', icon: AdminIcon, end: false }
    : { to: '/login', label: 'Login', icon: LoginIcon, end: false }

  const items = [...navItems, authItem]

  return (
    <nav className="fixed inset-x-0 bottom-0 z-10 flex border-t border-slate-800 bg-brand-bg pb-[env(safe-area-inset-bottom)] md:hidden">
      {items.map(({ to, label, icon: Icon, end }) => (
        <NavLink
          key={to}
          to={to}
          end={end}
          className={({ isActive }) =>
            [
              'flex flex-1 flex-col items-center gap-0.5 py-2.5 text-[11px] font-medium transition-colors',
              isActive ? 'text-brand-accent' : 'text-slate-400',
            ].join(' ')
          }
        >
          <Icon className="h-5 w-5 shrink-0" />
          {label}
        </NavLink>
      ))}
    </nav>
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

function LoginIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" {...props}>
      <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M15 12H3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function AdminIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" {...props}>
      <circle cx="12" cy="8" r="4" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M5 21c0-4 3-7 7-7s7 3 7 7" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}
