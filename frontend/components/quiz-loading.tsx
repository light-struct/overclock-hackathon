"use client"

import { GraduationCap, Loader2 } from "lucide-react"
import type { QuizConfig } from "@/lib/quiz-types"

interface QuizLoadingProps {
  config: QuizConfig
}

export function QuizLoading({ config }: QuizLoadingProps) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background px-6">
      <div className="text-center">
        <div className="mb-6 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-accent">
          <GraduationCap className="h-8 w-8 text-primary" />
        </div>
        <div className="mb-4 flex items-center justify-center gap-2">
          <Loader2 className="h-5 w-5 animate-spin text-primary" />
          <span className="text-lg font-bold text-foreground">
            Generating Your Quiz...
          </span>
        </div>
        <p className="text-sm text-muted-foreground">
          Creating {config.numQuestions} questions about{" "}
          <span className="font-semibold text-foreground">{config.topic}</span>{" "}
          at {config.difficulty} difficulty
        </p>
      </div>
    </div>
  )
}
