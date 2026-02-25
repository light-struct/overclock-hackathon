"use client"

import {
  GraduationCap,
  CheckCircle2,
  XCircle,
  RotateCcw,
  Trophy,
  Target,
  TrendingUp,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"
import type { QuizConfig, QuestionResult } from "@/lib/quiz-types"

interface QuizResultsProps {
  config: QuizConfig
  results: QuestionResult[]
  onRestart: () => void
}

export function QuizResults({ config, results, onRestart }: QuizResultsProps) {
  const totalScore = results.reduce((sum, r) => sum + r.score, 0)
  const maxScore = results.length
  const percentage = Math.round((totalScore / maxScore) * 100)

  const correctCount = results.filter((r) => r.score >= 1).length
  const partialCount = results.filter((r) => r.score > 0 && r.score < 1).length
  const incorrectCount = results.filter((r) => r.score === 0).length

  const getRating = () => {
    if (percentage >= 90) return { label: "Excellent", color: "text-[var(--success)]" }
    if (percentage >= 70) return { label: "Good", color: "text-[var(--success)]" }
    if (percentage >= 50) return { label: "Fair", color: "text-amber-500" }
    return { label: "Needs Improvement", color: "text-destructive" }
  }

  const rating = getRating()

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Header */}
      <header className="shrink-0 border-b border-border px-4 py-3">
        <div className="mx-auto flex max-w-3xl items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-6 w-6 items-center justify-center rounded bg-primary">
              <GraduationCap className="h-4 w-4 text-primary-foreground" />
            </div>
            <span className="text-sm font-bold text-foreground">QuizAgent</span>
          </div>
          <Button variant="outline" size="sm" onClick={onRestart} className="gap-1.5">
            <RotateCcw className="h-3.5 w-3.5" />
            New Quiz
          </Button>
        </div>
      </header>

      <main className="flex-1 px-4 py-8">
        <div className="mx-auto max-w-3xl">
          {/* Score Card */}
          <div className="mb-8 rounded-2xl border border-border bg-card p-8 text-center">
            <div className="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-full bg-accent">
              <Trophy className="h-8 w-8 text-primary" />
            </div>
            <h1 className="text-3xl font-bold text-foreground">Quiz Complete!</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              {config.topic} - {config.difficulty} difficulty
            </p>

            {/* Score circle */}
            <div className="mx-auto mt-8 flex h-32 w-32 flex-col items-center justify-center rounded-full border-4 border-primary">
              <span className="text-4xl font-bold text-primary">{percentage}%</span>
              <span className="text-xs text-muted-foreground">
                {totalScore}/{maxScore}
              </span>
            </div>

            <p className={`mt-4 text-lg font-bold ${rating.color}`}>{rating.label}</p>

            {/* Stats */}
            <div className="mt-8 grid grid-cols-3 gap-4">
              <div className="rounded-xl border border-border bg-background p-4">
                <div className="flex items-center justify-center gap-1.5 text-[var(--success)]">
                  <CheckCircle2 className="h-4 w-4" />
                  <span className="text-2xl font-bold">{correctCount}</span>
                </div>
                <p className="mt-1 text-xs text-muted-foreground">Correct</p>
              </div>
              <div className="rounded-xl border border-border bg-background p-4">
                <div className="flex items-center justify-center gap-1.5 text-amber-500">
                  <Target className="h-4 w-4" />
                  <span className="text-2xl font-bold">{partialCount}</span>
                </div>
                <p className="mt-1 text-xs text-muted-foreground">Partial</p>
              </div>
              <div className="rounded-xl border border-border bg-background p-4">
                <div className="flex items-center justify-center gap-1.5 text-destructive">
                  <XCircle className="h-4 w-4" />
                  <span className="text-2xl font-bold">{incorrectCount}</span>
                </div>
                <p className="mt-1 text-xs text-muted-foreground">Incorrect</p>
              </div>
            </div>
          </div>

          {/* Detailed Results */}
          <div className="mb-6 flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-primary" />
            <h2 className="text-lg font-bold text-foreground">Question Breakdown</h2>
          </div>

          <div className="space-y-4">
            {results.map((result, i) => (
              <div
                key={i}
                className={`rounded-xl border p-5 ${
                  result.score >= 1
                    ? "border-[var(--success)]/20 bg-[var(--success)]/5"
                    : result.score > 0
                      ? "border-amber-400/20 bg-amber-50"
                      : "border-destructive/20 bg-destructive/5"
                }`}
              >
                <div className="mb-3 flex items-start justify-between gap-3">
                  <div className="flex items-start gap-2">
                    <span className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded bg-accent text-xs font-bold text-primary">
                      {result.questionNumber}
                    </span>
                    <p className="text-sm font-medium leading-relaxed text-foreground">
                      {result.question}
                    </p>
                  </div>
                  <div className="flex shrink-0 items-center gap-1">
                    {result.score >= 1 ? (
                      <CheckCircle2 className="h-5 w-5 text-[var(--success)]" />
                    ) : result.score > 0 ? (
                      <CheckCircle2 className="h-5 w-5 text-amber-500" />
                    ) : (
                      <XCircle className="h-5 w-5 text-destructive" />
                    )}
                    <span className="text-xs font-bold text-muted-foreground">
                      {result.score}/1
                    </span>
                  </div>
                </div>

                <div className="ml-8 space-y-1.5">
                  <div className="flex items-start gap-2 text-sm">
                    <span className="shrink-0 font-semibold text-muted-foreground">
                      Your answer:
                    </span>
                    <span
                      className={
                        result.score >= 1
                          ? "text-[var(--success)]"
                          : "text-destructive"
                      }
                    >
                      {result.userAnswer}
                    </span>
                  </div>
                  {result.score < 1 && (
                    <div className="flex items-start gap-2 text-sm">
                      <span className="shrink-0 font-semibold text-muted-foreground">
                        Correct:
                      </span>
                      <span className="text-[var(--success)]">{result.correctAnswer}</span>
                    </div>
                  )}
                  <p className="text-xs leading-relaxed text-muted-foreground">
                    {result.explanation}
                  </p>
                </div>
              </div>
            ))}
          </div>

          {/* Overall progress bar */}
          <div className="mt-8 rounded-xl border border-border bg-card p-5">
            <div className="mb-2 flex items-center justify-between text-sm">
              <span className="font-semibold text-foreground">Overall Score</span>
              <span className="font-bold text-primary">{percentage}%</span>
            </div>
            <Progress value={percentage} className="h-3" />
          </div>

          {/* CTA */}
          <div className="mt-8 text-center">
            <Button onClick={onRestart} size="lg" className="gap-2">
              <RotateCcw className="h-4 w-4" />
              Take Another Quiz
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
