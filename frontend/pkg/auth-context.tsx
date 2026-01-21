"use client"

import React, { createContext, useContext, useState, useEffect } from 'react'

interface AuthContextType {
    user: any | null
    loading: boolean
    login: (phone: string, password: string) => Promise<void>
    logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<any | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Basic session recovery logic
        const savedUser = localStorage.getItem('bjdms_user')
        if (savedUser) {
            setUser(JSON.parse(savedUser))
        }
        setLoading(false)
    }, [])

    const login = async (phone: string, password: string) => {
        // In production, this would call the API
        // const res = await fetch('/api/v1/auth/login', { ... })

        // Mock for now
        const mockUser = { id: '1', name: 'Demo User', role: 'admin' }
        setUser(mockUser)
        localStorage.setItem('bjdms_user', JSON.stringify(mockUser))
    }

    const logout = () => {
        setUser(null)
        localStorage.removeItem('bjdms_user')
    }

    return (
        <AuthContext.Provider value={{ user, loading, login, logout }}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
