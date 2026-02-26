"use client"

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useApp } from '@/lib/app-context'
import { Button } from '@/components/ui/button'

interface User {
  id: string
  name: string
  email: string
  role?: 'student' | 'admin'
}

interface Attempt {
  id: string
  user_id: string
  subject?: string
  topic?: string
  score: number
  created_at: string | null
  student_name?: string
}

export function ProfileCard() {
  const { token, logout, t } = useApp()

  const [user, setUser] = useState<User | null>(null)
  const [allUsers, setAllUsers] = useState<User[]>([])
  const [attempts, setAttempts] = useState<Attempt[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [hydrated, setHydrated] = useState(false)
  const [hasLocalToken, setHasLocalToken] = useState<boolean | null>(null)

  useEffect(() => {
    if (typeof window !== 'undefined') {
      setHasLocalToken(!!localStorage.getItem('token'))
      setHydrated(true)
    }
  }, [])

  useEffect(() => {
    const fetchProfile = async () => {
      const tok = token ?? localStorage.getItem('token')
      if (!tok) return

      setLoading(true)
      setError('')

      try {
        // Получаем текущего пользователя
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/me`, {
          headers: { Authorization: `Bearer ${tok}` },
        })

        if (!res.ok) throw new Error('Failed to load profile')
        const data: User = await res.json()
        setUser(data)

        // Если админ - загружаем всех пользователей
        if (data.role === 'admin') {
          const usersRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/users`, {
            headers: { Authorization: `Bearer ${tok}` },
          })
          if (usersRes.ok) {
            const usersData = await usersRes.json()
            setAllUsers(usersData.users || [])
          }
        }

        // Получаем все попытки
        const attemptsRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/exam/attempts`, {
          headers: { Authorization: `Bearer ${tok}` },
        })

        if (attemptsRes.ok) {
          const d = await attemptsRes.json()
          let arr: Attempt[] = []

          if (Array.isArray(d.attempts)) arr = d.attempts
          else if (Array.isArray(d)) arr = d
          else if (d && typeof d === 'object') arr = Object.values(d)

          arr = arr.filter(a => a)
          // Убираем дубликаты по id (одна попытка — одна запись)
          const seen = new Set<string>()
          arr = arr.filter(a => {
            const key = String(a.id ?? '')
            if (seen.has(key)) return false
            seen.add(key)
            return true
          })

          if (data.role === 'admin') {
            const userMap = new Map(allUsers.map(u => [u.id, u.name]))
            const adminAttempts = arr.map(a => ({
              id: a.id ?? 'unknown',
              user_id: a.user_id ?? 'unknown',
              subject: a.subject,
              topic: a.topic,
              score: a.score ?? 0,
              created_at: a.created_at ?? null,
              student_name: userMap.get(String(a.user_id)) || 'Unknown'
            }))
            setAttempts(adminAttempts)
          } else {
            const ownAttempts = arr
              .filter(a => a.user_id === data.id)
              .map(a => ({
                id: a.id ?? 'unknown',
                user_id: a.user_id ?? data.id,
                subject: a.subject,
                topic: a.topic,
                score: a.score ?? 0,
                created_at: a.created_at ?? null
              }))
            setAttempts(ownAttempts)
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load profile')
      } finally {
        setLoading(false)
      }
    }

    fetchProfile()
  }, [token, allUsers])

  const formatDateSafe = (dateStr: string | null | undefined) => {
    if (!dateStr) return '-'
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return '-'
    return d.toLocaleString()
  }

  if (!token && hydrated && !hasLocalToken) {
    return (
      <div className="mx-auto max-w-md text-center">
        <h2 className="text-2xl font-bold mb-4">{t.profile.title}</h2>
        <p className="text-sm text-muted-foreground mb-4">{t.profile.notLoggedIn}</p>
        <Link href="/login">
          <Button className="w-full">{t.profile.signIn}</Button>
        </Link>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-md">
      <h2 className="text-2xl font-bold mb-6">{t.profile.title}</h2>

      <div className="space-y-4">
        {loading && <p className="text-sm text-muted-foreground">Loading...</p>}
        {error && <p className="text-sm text-destructive">{error}</p>}

        {user && (
          <div className="rounded-lg border bg-card p-6">
            <p className="text-sm text-muted-foreground">Name</p>
            <p className="text-lg font-medium mb-3">{user.name}</p>

            <p className="text-sm text-muted-foreground">Email</p>
            <p className="text-sm mb-4">{user.email}</p>

            <div className="flex gap-3">
              <Button onClick={logout} variant="ghost" className="flex-1">
                {t.header.logout}
              </Button>
            </div>
          </div>
        )}

        {attempts && (
          <div className="mt-6">
            <h3 className="text-lg font-semibold mb-3">{t.profile.attemptsHeading}</h3>

            {attempts.length === 0 ? (
              <p className="text-sm text-muted-foreground">{t.profile.noAttempts}</p>
            ) : (
              <ul className="space-y-3">
                {attempts.map(a => {
                  const displayDate = formatDateSafe(a.created_at)
                  const studentName = a.student_name ?? ''
                  return (
                    <li key={a.id ?? JSON.stringify(a)} className="rounded-md border p-3 bg-card">
                      <div className="flex items-center justify-between gap-2 flex-wrap">
                        <div>
                          {user?.role === 'admin' && studentName && (
                            <div className="text-sm text-muted-foreground">Student: {studentName}</div>
                          )}
                          <div className="text-sm text-muted-foreground">{t.profile.date}</div>
                          <div className="text-sm">{displayDate}</div>
                          {(a.subject || a.topic) && (
                            <div className="text-xs text-muted-foreground mt-1">
                              {[a.subject, a.topic].filter(Boolean).join(' · ')}
                            </div>
                          )}
                        </div>
                        <div className="flex items-center gap-3">
                          <div className="text-right">
                            <div className="text-sm text-muted-foreground">{t.profile.score}</div>
                            <div className="text-lg font-medium">
                              {typeof a.score === 'number' ? a.score : a.score ? String(a.score) : '-'}
                            </div>
                          </div>
                          {user?.role !== 'admin' && (
                            <Link href={`/profile/attempts/${a.id}`}>
                              <Button variant="outline" size="sm">{t.profile.viewAttempt}</Button>
                            </Link>
                          )}
                        </div>
                      </div>
                    </li>
                  )
                })}
              </ul>
            )}
          </div>
        )}
      </div>
    </div>
  )
}