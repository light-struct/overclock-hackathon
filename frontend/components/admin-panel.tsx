"use client"

import { useEffect, useState } from 'react'
import { useApp } from '@/lib/app-context'
import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { Users, FileText, Shield } from 'lucide-react'

interface User {
  id: string
  name: string
  email: string
  role?: 'student' | 'admin'
  created_at?: string
}

interface Attempt {
  id: string
  user_id: string
  subject: string
  topic: string
  score: number
  created_at: string
}

export function AdminPanel() {
  const { token, t } = useApp()
  const [user, setUser] = useState<User | null>(null)
  const [allUsers, setAllUsers] = useState<User[]>([])
  const [attempts, setAttempts] = useState<Attempt[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchData = async () => {
      const tok = token || localStorage.getItem('token')
      if (!tok) {
        setError('Not authenticated')
        setLoading(false)
        return
      }

      try {
        // Check if user is admin
        const meRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/me`, {
          headers: { Authorization: `Bearer ${tok}` },
        })
        
        if (!meRes.ok) throw new Error('Failed to load profile')
        const meData: User = await meRes.json()
        setUser(meData)

        if (meData.role !== 'admin') {
          setError('Access denied. Admin only.')
          setLoading(false)
          return
        }

        // Fetch all users
        const usersRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/users`, {
          headers: { Authorization: `Bearer ${tok}` },
        })
        
        if (usersRes.ok) {
          const usersData = await usersRes.json()
          setAllUsers(usersData.users || [])
        }

        // Fetch all attempts
        const attemptsRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/exam/attempts`, {
          headers: { Authorization: `Bearer ${tok}` },
        })
        
        if (attemptsRes.ok) {
          const attemptsData = await attemptsRes.json()
          setAttempts(attemptsData.attempts || [])
        }

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [token])

  const getUserName = (userId: string) => {
    const user = allUsers.find(u => u.id === userId)
    return user?.name || 'Unknown'
  }

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return '-'
    return d.toLocaleString()
  }

  if (loading) {
    return <div className="text-center">{t.admin.loading}</div>
  }

  if (error) {
    return (
      <div className="mx-auto max-w-md text-center">
        <Shield className="mx-auto h-12 w-12 text-destructive mb-4" />
        <h2 className="text-2xl font-bold mb-4">{t.admin.accessDenied}</h2>
        <p className="text-muted-foreground mb-4">{error}</p>
        <Link href="/">
          <Button>{t.admin.goHome}</Button>
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold mb-2">{t.admin.title}</h1>
        <p className="text-muted-foreground">{t.admin.subtitle}</p>
      </div>

      {/* Users Section */}
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2 mb-4">
          <Users className="h-5 w-5 text-primary" />
          <h2 className="text-xl font-semibold">{t.admin.allUsers} ({allUsers.length})</h2>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left py-2 px-4">{t.admin.name}</th>
                <th className="text-left py-2 px-4">{t.admin.email}</th>
                <th className="text-left py-2 px-4">{t.admin.role}</th>
                <th className="text-left py-2 px-4">{t.admin.registered}</th>
              </tr>
            </thead>
            <tbody>
              {allUsers.map(u => (
                <tr key={u.id} className="border-b hover:bg-accent">
                  <td className="py-3 px-4">{u.name}</td>
                  <td className="py-3 px-4 text-muted-foreground">{u.email}</td>
                  <td className="py-3 px-4">
                    <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                      u.role === 'admin' ? 'bg-primary/10 text-primary' : 'bg-secondary text-secondary-foreground'
                    }`}>
                      {u.role || 'student'}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-sm text-muted-foreground">
                    {u.created_at ? formatDate(u.created_at) : '-'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Attempts Section */}
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2 mb-4">
          <FileText className="h-5 w-5 text-primary" />
          <h2 className="text-xl font-semibold">{t.admin.allAttempts} ({attempts.length})</h2>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left py-2 px-4">{t.admin.student}</th>
                <th className="text-left py-2 px-4">{t.admin.subject}</th>
                <th className="text-left py-2 px-4">{t.admin.difficulty}</th>
                <th className="text-left py-2 px-4">{t.admin.score}</th>
                <th className="text-left py-2 px-4">{t.admin.date}</th>
              </tr>
            </thead>
            <tbody>
              {attempts.map(a => (
                <tr key={a.id} className="border-b hover:bg-accent">
                  <td className="py-3 px-4 font-medium">{getUserName(a.user_id)}</td>
                  <td className="py-3 px-4">{a.subject}</td>
                  <td className="py-3 px-4">
                    <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                      a.topic === 'easy' ? 'bg-green-100 text-green-800' :
                      a.topic === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-red-100 text-red-800'
                    }`}>
                      {a.topic}
                    </span>
                  </td>
                  <td className="py-3 px-4">
                    <span className={`font-bold ${
                      a.score >= 80 ? 'text-green-600' :
                      a.score >= 60 ? 'text-yellow-600' :
                      'text-red-600'
                    }`}>
                      {a.score}%
                    </span>
                  </td>
                  <td className="py-3 px-4 text-sm text-muted-foreground">
                    {formatDate(a.created_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
