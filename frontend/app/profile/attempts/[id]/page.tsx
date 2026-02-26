"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { Header } from "@/components/header"
import { Footer } from "@/components/footer"
import { Button } from "@/components/ui/button"
import { useApp } from "@/lib/app-context"
import { ArrowLeft, CheckCircle2, XCircle, Target, Loader2, MessageSquare } from "lucide-react"
import type { QuestionResult } from "@/lib/quiz-types"

interface Attempt {
  id: number
  user_id: number
  subject: string
  topic: string
  score: number
  results: string
  created_at: string
}

export default function AttemptDetailPage() {
  const params = useParams()
  const id = params?.id as string | undefined
  const { token, t } = useApp()
  const [attempt, setAttempt] = useState<Attempt | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [explainingIndex, setExplainingIndex] = useState<number | null>(null)
  const [aiExplanations, setAiExplanations] = useState<Record<number, string>>({})

  useEffect(() => {
    const tok = token ?? (typeof window !== "undefined" ? localStorage.getItem("token") : null)
    if (!tok) {
      setLoading(false)
      setError("Not authenticated")
      return
    }
    if (!id || id === "undefined" || id === "null") {
      setLoading(false)
      setError("Invalid attempt")
      return
    }
    const apiUrl = process.env.NEXT_PUBLIC_API_URL
    if (!apiUrl) {
      setLoading(false)
      setError("API URL not configured")
      return
    }
    fetch(`${apiUrl}/exam/attempts/${id}`, {
      headers: { Authorization: `Bearer ${tok}` },
    })
      .then(async (res) => {
        if (!res.ok) {
          const errBody = await res.json().catch(() => ({}))
          throw new Error((errBody as { error?: string }).error || `HTTP ${res.status}`)
        }
        return res.json()
      })
      .then((data: Attempt) => {
        setAttempt(data)
        setError("")
      })
      .catch((e) => {
        setError(e instanceof Error ? e.message : "Failed to load attempt")
        setAttempt(null)
      })
      .finally(() => setLoading(false))
  }, [id, token])

  const results: QuestionResult[] = attempt?.results
    ? (() => {
        try {
          const parsed = JSON.parse(attempt.results)
          return Array.isArray(parsed) ? parsed : []
        } catch {
          return []
        }
      })()
    : []

  const handleExplain = async (index: number, r: QuestionResult) => {
    if (aiExplanations[index]) return
    setExplainingIndex(index)
    try {
      const res = await fetch("/api/quiz/explain", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          question: r.question,
          userAnswer: r.userAnswer,
          correctAnswer: r.correctAnswer,
          score: r.score,
          existingExplanation: r.explanation,
        }),
      })
      const data = await res.json()
      setAiExplanations((prev) => ({ ...prev, [index]: data.explanation || r.explanation }))
    } catch {
      setAiExplanations((prev) => ({ ...prev, [index]: r.explanation || "Error loading explanation." }))
    } finally {
      setExplainingIndex(null)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container py-12 flex items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </main>
        <Footer />
      </div>
    )
  }

  if (error || !attempt) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <main className="container py-12 max-w-md">
          <p className="text-destructive font-medium">{error || "Attempt not found"}</p>
          <p className="text-sm text-muted-foreground mt-2">
            Попытка не найдена или у вас нет доступа. Откройте попытку по кнопке «Смотреть» в разделе «Результаты тестов» в профиле.
          </p>
          <Button asChild variant="outline" className="mt-4">
            <Link href="/profile">{t.profile.backToProfile}</Link>
          </Button>
        </main>
        <Footer />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container py-8 max-w-3xl mx-auto">
        <Button asChild variant="ghost" size="sm" className="mb-6 gap-2">
          <Link href="/profile">
            <ArrowLeft className="h-4 w-4" />
            {t.profile.backToProfile}
          </Link>
        </Button>

        <div className="rounded-xl border bg-card p-6 mb-8">
          <h1 className="text-xl font-bold">
            {attempt.subject} — {attempt.topic}
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            {new Date(attempt.created_at).toLocaleString()} · {t.profile.score}: {attempt.score}%
          </p>
        </div>

        {results.length === 0 ? (
          <p className="text-muted-foreground">
            Detailed results were not saved for this attempt. New attempts will show each question and answer here.
          </p>
        ) : (
          <div className="space-y-4">
            {results.map((r, i) => (
              <div
                key={i}
                className={`rounded-xl border p-5 ${
                  r.score >= 1
                    ? "border-[var(--success)]/20 bg-[var(--success)]/5"
                    : r.score > 0
                      ? "border-amber-400/20 bg-amber-500/5"
                      : "border-destructive/20 bg-destructive/5"
                }`}
              >
                <div className="mb-3 flex items-start justify-between gap-3">
                  <div className="flex items-start gap-2">
                    <span className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded bg-accent text-xs font-bold text-primary">
                      {r.questionNumber}
                    </span>
                    <p className="text-sm font-medium leading-relaxed text-foreground">{r.question}</p>
                  </div>
                  <div className="flex shrink-0 items-center gap-1">
                    {r.score >= 1 ? (
                      <CheckCircle2 className="h-5 w-5 text-[var(--success)]" />
                    ) : r.score > 0 ? (
                      <Target className="h-5 w-5 text-amber-500" />
                    ) : (
                      <XCircle className="h-5 w-5 text-destructive" />
                    )}
                    <span className="text-xs font-bold text-muted-foreground">{r.score}/1</span>
                  </div>
                </div>

                <div className="ml-8 space-y-1.5">
                  <div className="flex items-start gap-2 text-sm">
                    <span className="shrink-0 font-semibold text-muted-foreground">
                      {t.profile.yourAnswer}:
                    </span>
                    <span className={r.score >= 1 ? "text-[var(--success)]" : "text-destructive"}>
                      {r.userAnswer || "—"}
                    </span>
                  </div>
                  {r.score < 1 && (
                    <div className="flex items-start gap-2 text-sm">
                      <span className="shrink-0 font-semibold text-muted-foreground">
                        {t.profile.correctAnswerLabel}:
                      </span>
                      <span className="text-[var(--success)]">{r.correctAnswer}</span>
                    </div>
                  )}
                  <p className="text-xs leading-relaxed text-muted-foreground">{r.explanation}</p>

                  {aiExplanations[i] && (
                    <div className="mt-3 rounded-lg border border-border bg-muted/30 p-3">
                      <p className="text-xs font-semibold text-muted-foreground mb-1">
                        {t.profile.aiExplanation}
                      </p>
                      <p className="text-sm leading-relaxed">{aiExplanations[i]}</p>
                    </div>
                  )}

                  <Button
                    variant="outline"
                    size="sm"
                    className="mt-2 gap-2"
                    onClick={() => handleExplain(i, r)}
                    disabled={explainingIndex === i}
                  >
                    {explainingIndex === i ? (
                      <Loader2 className="h-3.5 w-3.5 animate-spin" />
                    ) : (
                      <MessageSquare className="h-3.5 w-3.5" />
                    )}
                    {t.profile.explain}
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
      <Footer />
    </div>
  )
}
