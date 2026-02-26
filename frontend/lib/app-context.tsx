"use client"

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { Language, translations } from './translations'

interface AppContextType {
  language: Language
  setLanguage: (lang: Language) => void
  t: typeof translations.en
  token: string | null
  setToken: (token: string | null) => void
  logout: () => void
}

const AppContext = createContext<AppContextType | undefined>(undefined)

export function AppProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>('en')
  const [token, setToken] = useState<string | null>(null)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    const savedLang = localStorage.getItem('language') as Language
    if (savedLang && ['kk', 'ru', 'en'].includes(savedLang)) {
      setLanguage(savedLang)
    }
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      setToken(savedToken)
    }
    setMounted(true)
  }, [])

  const handleSetLanguage = (lang: Language) => {
    setLanguage(lang)
    localStorage.setItem('language', lang)
    // Force re-render
    setMounted(false)
    setTimeout(() => setMounted(true), 0)
  }

  const handleSetToken = (newToken: string | null) => {
    setToken(newToken)
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  const logout = () => {
    handleSetToken(null)
  }

  if (!mounted) {
    return null
  }

  return (
    <AppContext.Provider value={{
      language,
      setLanguage: handleSetLanguage,
      t: translations[language],
      token,
      setToken: handleSetToken,
      logout
    }}>
      {children}
    </AppContext.Provider>
  )
}

export function useApp() {
  const context = useContext(AppContext)
  if (!context) {
    throw new Error('useApp must be used within AppProvider')
  }
  return context
}
