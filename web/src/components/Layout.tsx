import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import BottomNav from './BottomNav'

export default function Layout() {
  return (
    <div className="flex min-h-screen bg-app-bg md:flex-row">
      <Sidebar />
      <main className="flex-1 overflow-y-auto px-4 py-6 pb-24 sm:px-6 md:px-8 md:py-8 md:pb-8">
        <div className="mx-auto max-w-6xl">
          <Outlet />
        </div>
      </main>
      <BottomNav />
    </div>
  )
}
