import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'

type AuthContextValue = {
    token: string
    isAuthenticated: boolean
    login: (token: string) => void
    logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
    const [token, setToken] = useState(() => localStorage.getItem('access_token') ?? '')

    const login = useCallback((newToken: string) => {
        setToken(newToken)
        localStorage.setItem('access_token', newToken)
    }, [])

    const logout = useCallback(() => {
        setToken('')
        localStorage.removeItem('access_token')
    }, [])

    const value: AuthContextValue = {
        token,
        isAuthenticated: token !== '',
        login,
        logout,
    }

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
    const ctx = useContext(AuthContext)
    if (!ctx) throw new Error('useAuth must be used within an AuthProvider')
        return ctx
}
