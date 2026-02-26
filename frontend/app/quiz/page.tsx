"use client"

import { useState } from "react"
import { QuizSetup } from "@/components/quiz-setup"
import { QuizLoading } from "@/components/quiz-loading"
import { QuizActive } from "@/components/quiz-active"
import { QuizResults } from "@/components/quiz-results"
import type { QuizConfig, QuestionResult } from "@/lib/quiz-types"

type QuizState =
  | { phase: "setup" }
  | { phase: "loading"; config: QuizConfig }
  | { phase: "active"; config: QuizConfig; questions: Question[] }
  | { phase: "results"; config: QuizConfig; results: QuestionResult[] }

interface Question {
  questionNumber: number
  question: string
  options: string[]
  correctAnswer: string
}

export default function QuizPage() {
  const [state, setState] = useState<QuizState>({ phase: "setup" })

  const handleStart = async (config: QuizConfig) => {
    setState({ phase: "loading", config })

    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/exam/quiz/generate`, {
        method: "POST",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
          topic: config.topic,
          numQuestions: config.numQuestions,
          difficulty: config.difficulty,
        }),
      })

      if (!response.ok) {
        const data = await response.json().catch(() => ({ error: "Failed to generate quiz" }))
        throw new Error(data.error || "Failed to generate quiz")
      }

      const data = await response.json()
      setState({ phase: "active", config, questions: data.questions })
    } catch (error) {
      alert(error instanceof Error ? error.message : "Failed to generate quiz. Check your API key.")
      setState({ phase: "setup" })
    }
  }

  const handleFinish = (results: QuestionResult[]) => {
    if (state.phase !== "active") return
    setState({ phase: "results", config: state.config, results })
  }

  const handleReset = () => {
    setState({ phase: "setup" })
  }

  switch (state.phase) {
    case "setup":
      return <QuizSetup onStart={handleStart} />
    case "loading":
      return <QuizLoading config={state.config} />
    case "active":
      return (
        <QuizActive
          config={state.config}
          questions={state.questions}
          onFinish={handleFinish}
          onReset={handleReset}
        />
      )
    case "results":
      return (
        <QuizResults
          config={state.config}
          results={state.results}
          onRestart={handleReset}
        />
      )
  }
}
