"use client"

import { useState } from "react"
import {
  GraduationCap,
  CheckCircle2,
  XCircle,
  Loader2,
  ArrowRight,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import type { QuizConfig, QuestionResult } from "@/lib/quiz-types"

// app/api/quiz/evaluate/route.ts
const apiKey = process.env.GEMINI_API_KEY 

interface Question {
  questionNumber: number
  question: string
  options: string[]
  correctAnswer: string
}

interface QuizActiveProps {
  config: QuizConfig
  questions: Question[]
  onFinish: (results: QuestionResult[]) => void
  onReset: () => void
}

export function QuizActive({
  config,
  questions,
  onFinish,
  onReset,
}: QuizActiveProps) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null)
  const [isEvaluating, setIsEvaluating] = useState(false)
  const [feedback, setFeedback] = useState<{
    score: number
    correctAnswer: string
    explanation: string
  } | null>(null)
  const [results, setResults] = useState<QuestionResult[]>([])

  const currentQuestion = questions[currentIndex]
  const isLastQuestion = currentIndex === questions.length - 1
  const progress = ((currentIndex + (feedback ? 1 : 0)) / questions.length) * 100

  const handleSelectAnswer = async (answer: string) => {
    if (isEvaluating || feedback) return
    setSelectedAnswer(answer)
    setIsEvaluating(true)

    try {
      const response = await fetch("/api/quiz/evaluate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          question: currentQuestion.question,
          userAnswer: answer,
          correctAnswer: currentQuestion.correctAnswer,
        }),
      })

      if (!response.ok) {
        throw new Error("Failed to evaluate answer")
      }

      const evaluation: { score: number; correctAnswer: string; explanation: string } =
        await response.json()

      setFeedback(evaluation)

      const result: QuestionResult = {
        questionNumber: currentQuestion.questionNumber,
        question: currentQuestion.question,
        userAnswer: answer,
        correctAnswer: evaluation.correctAnswer,
        score: evaluation.score,
        explanation: evaluation.explanation,
      }

      setResults((prev) => [...prev, result])
    } catch {
      // Fallback: simple comparison
      const isCorrect =
        answer.toLowerCase().trim() ===
        currentQuestion.correctAnswer.toLowerCase().trim()
      const fallbackFeedback = {
        score: isCorrect ? 1 : 0,
        correctAnswer: currentQuestion.correctAnswer,
        explanation: isCorrect
          ? "Your answer is correct!"
          : `The correct answer is: ${currentQuestion.correctAnswer}`,
      }
      setFeedback(fallbackFeedback)

      setResults((prev) => [
        ...prev,
        {
          questionNumber: currentQuestion.questionNumber,
          question: currentQuestion.question,
          userAnswer: answer,
          correctAnswer: currentQuestion.correctAnswer,
          score: fallbackFeedback.score,
          explanation: fallbackFeedback.explanation,
        },
      ])
    } finally {
      setIsEvaluating(false)
    }
  }

  const handleNext = () => {
    if (isLastQuestion) {
      onFinish([...results])
    } else {
      setCurrentIndex((prev) => prev + 1)
      setSelectedAnswer(null)
      setFeedback(null)
    }
  }

  const getOptionStyle = (option: string) => {
    if (!feedback) {
      if (selectedAnswer === option && isEvaluating) {
        return "border-primary bg-accent text-primary"
      }
      return "border-border bg-card text-card-foreground hover:border-primary/40 hover:bg-accent/50 cursor-pointer"
    }

    const isSelected = selectedAnswer === option
    const isCorrect = option === feedback.correctAnswer

    if (isCorrect) {
      return "border-[var(--success)] bg-[var(--success)]/10 text-[var(--success)]"
    }
    if (isSelected && !isCorrect) {
      return "border-destructive bg-destructive/10 text-destructive"
    }
    return "border-border bg-muted/30 text-muted-foreground opacity-60"
  }

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Header */}
      <header className="shrink-0 border-b border-border px-4 py-3">
        <div className="mx-auto flex max-w-2xl items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-6 w-6 items-center justify-center rounded bg-primary">
              <GraduationCap className="h-4 w-4 text-primary-foreground" />
            </div>
            <span className="text-sm font-bold text-foreground">QuizAgent</span>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="border-primary/20 text-xs capitalize">
              {config.difficulty}
            </Badge>
            <Badge variant="secondary" className="max-w-[140px] truncate text-xs">
              {config.topic}
            </Badge>
          </div>
        </div>
      </header>

      {/* Progress */}
      <div className="border-b border-border px-4 py-3">
        <div className="mx-auto max-w-2xl">
          <div className="mb-2 flex items-center justify-between text-xs text-muted-foreground">
            <span>
              Question {currentIndex + 1} of {questions.length}
            </span>
            <span>{Math.round(progress)}% complete</span>
          </div>
          <Progress value={progress} className="h-2" />
        </div>
      </div>

      {/* Question */}
      <main className="flex flex-1 flex-col px-4 py-8">
        <div className="mx-auto w-full max-w-2xl flex-1">
          <div className="mb-8">
            <span className="mb-2 inline-block rounded bg-accent px-2.5 py-1 text-xs font-semibold text-primary">
              Question {currentQuestion.questionNumber}
            </span>
            <h2 className="text-xl font-bold leading-relaxed text-foreground">
              {currentQuestion.question}
            </h2>
          </div>

          {/* Options */}
          <div className="space-y-3">
            {currentQuestion.options.map((option, i) => {
              const letter = String.fromCharCode(65 + i)
              return (
                <button
                  key={i}
                  type="button"
                  disabled={!!feedback || isEvaluating}
                  onClick={() => handleSelectAnswer(option)}
                  className={`flex w-full items-start gap-3 rounded-xl border p-4 text-left transition-all ${getOptionStyle(option)}`}
                >
                  <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border border-current/20 text-xs font-bold">
                    {letter}
                  </span>
                  <span className="pt-0.5 text-sm font-medium leading-relaxed">
                    {option}
                  </span>
                  {feedback && option === feedback.correctAnswer && (
                    <CheckCircle2 className="ml-auto h-5 w-5 shrink-0 text-[var(--success)]" />
                  )}
                  {feedback &&
                    selectedAnswer === option &&
                    option !== feedback.correctAnswer && (
                      <XCircle className="ml-auto h-5 w-5 shrink-0 text-destructive" />
                    )}
                </button>
              )
            })}
          </div>

          {/* Evaluating */}
          {isEvaluating && (
            <div className="mt-6 flex items-center gap-2 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin text-primary" />
              Evaluating your answer...
            </div>
          )}

          {/* Feedback */}
          {feedback && (
            <div
              className={`mt-6 rounded-xl border p-4 ${
                feedback.score >= 1
                  ? "border-[var(--success)]/30 bg-[var(--success)]/5"
                  : feedback.score > 0
                    ? "border-amber-400/30 bg-amber-50"
                    : "border-destructive/30 bg-destructive/5"
              }`}
            >
              <div className="mb-1 flex items-center gap-2">
                {feedback.score >= 1 ? (
                  <CheckCircle2 className="h-5 w-5 text-[var(--success)]" />
                ) : feedback.score > 0 ? (
                  <CheckCircle2 className="h-5 w-5 text-amber-500" />
                ) : (
                  <XCircle className="h-5 w-5 text-destructive" />
                )}
                <span
                  className={`text-sm font-bold ${
                    feedback.score >= 1
                      ? "text-[var(--success)]"
                      : feedback.score > 0
                        ? "text-amber-600"
                        : "text-destructive"
                  }`}
                >
                  {feedback.score >= 1
                    ? "Correct!"
                    : feedback.score > 0
                      ? "Partially Correct"
                      : "Incorrect"}
                </span>
                <span className="ml-auto text-xs font-semibold text-muted-foreground">
                  Score: {feedback.score}/1
                </span>
              </div>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {feedback.explanation}
              </p>
            </div>
          )}
        </div>

        {/* Next / Finish button */}
        {feedback && (
          <div className="mx-auto mt-6 w-full max-w-2xl">
            <Button
              onClick={handleNext}
              size="lg"
              className="w-full gap-2"
            >
              {isLastQuestion ? "View Results" : "Next Question"}
              <ArrowRight className="h-4 w-4" />
            </Button>
          </div>
        )}
      </main>
    </div>
  )
}
