"use client"

import { useState } from "react"
import Link from "next/link"
import {
  ArrowLeft,
  GraduationCap,
  BookOpen,
  Flame,
  Sparkles,
  Zap,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useApp } from "@/lib/app-context"
import type { QuizConfig } from "@/lib/quiz-types"

const difficulties = [
  {
    value: "easy" as const,
    label: "easy",
    description: "basicFacts",
    icon: BookOpen,
  },
  {
    value: "medium" as const,
    label: "medium",
    description: "appliedKnowledge",
    icon: GraduationCap,
  },
  {
    value: "hard" as const,
    label: "hard",
    description: "advanced",
    icon: Flame,
  },
]

const aiProviders = [
  {
    value: "gemini" as const,
    label: "Gemini 2.5 Flash",
    icon: Sparkles,
  },
  {
    value: "groq" as const,
    label: "Groq Llama 3.3",
    icon: Zap,
  },
]

interface QuizSetupProps {
  onStart: (config: QuizConfig) => void
}

export function QuizSetup({ onStart }: QuizSetupProps) {
  const { t } = useApp()
  const [topic, setTopic] = useState("")
  const [numQuestions, setNumQuestions] = useState(5)
  const [difficulty, setDifficulty] = useState<"easy" | "medium" | "hard">(
    "medium"
  )
  const [aiProvider, setAiProvider] = useState<"gemini" | "groq">("gemini")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!topic.trim()) return
    onStart({ topic: topic.trim(), numQuestions, difficulty, aiProvider })
  }

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="border-b border-border px-6 py-4">
        <div className="mx-auto flex max-w-2xl items-center gap-3">
          <Button asChild variant="ghost" size="sm" className="text-muted-foreground hover:text-foreground">
            <Link href="/">
              <ArrowLeft className="mr-1 h-4 w-4" />
              {t.quiz.back}
            </Link>
          </Button>
          <div className="flex items-center gap-2">
            <div className="flex h-6 w-6 items-center justify-center rounded bg-primary">
              <GraduationCap className="h-4 w-4 text-primary-foreground" />
            </div>
            <span className="text-sm font-bold text-foreground">QuizAgent</span>
          </div>
        </div>
      </header>

      <main className="flex flex-1 items-center justify-center px-6 py-12">
        <form onSubmit={handleSubmit} className="w-full max-w-md space-y-6">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-foreground">{t.quiz.setup}</h1>
            <p className="mt-2 text-sm text-muted-foreground">
              {t.quiz.setupDescription}
            </p>
          </div>

          <div className="space-y-5">
            <div className="space-y-2">
              <Label htmlFor="topic" className="text-foreground">{t.quiz.topic}</Label>
              <Input
                id="topic"
                type="text"
                placeholder={t.quiz.topicPlaceholder}
                value={topic}
                onChange={(e) => setTopic(e.target.value)}
                className="bg-secondary text-foreground placeholder:text-muted-foreground"
                required
                autoFocus
              />
            </div>

            <div className="space-y-2">
              <Label className="text-foreground">{t.quiz.numQuestions}</Label>
              <div className="flex items-center gap-3">
                {[3, 5, 10, 15].map((n) => (
                  <button
                    key={n}
                    type="button"
                    onClick={() => setNumQuestions(n)}
                    className={`flex h-10 w-full items-center justify-center rounded-lg border text-sm font-semibold transition-colors ${
                      numQuestions === n
                        ? "border-primary bg-accent text-primary"
                        : "border-border bg-secondary text-muted-foreground hover:text-foreground hover:border-primary/30"
                    }`}
                  >
                    {n}
                  </button>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <Label className="text-foreground">{t.quiz.difficulty}</Label>
              <div className="grid grid-cols-3 gap-3">
                {difficulties.map((d) => (
                  <button
                    key={d.value}
                    type="button"
                    onClick={() => setDifficulty(d.value)}
                    className={`flex flex-col items-center gap-1.5 rounded-xl border p-4 transition-colors ${
                      difficulty === d.value
                        ? "border-primary bg-accent text-primary"
                        : "border-border bg-secondary text-muted-foreground hover:text-foreground hover:border-primary/30"
                    }`}
                  >
                    <d.icon className="h-5 w-5" />
                    <span className="text-xs font-semibold">{t.quiz[d.label]}</span>
                    <span className="text-[10px] opacity-70">{t.quiz[d.description]}</span>
                  </button>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <Label className="text-foreground">{t.quiz.aiProvider}</Label>
              <div className="grid grid-cols-2 gap-3">
                {aiProviders.map((provider) => (
                  <button
                    key={provider.value}
                    type="button"
                    onClick={() => setAiProvider(provider.value)}
                    className={`flex flex-col items-center gap-2 rounded-xl border p-4 transition-colors ${
                      aiProvider === provider.value
                        ? "border-primary bg-accent text-primary"
                        : "border-border bg-secondary text-muted-foreground hover:text-foreground hover:border-primary/30"
                    }`}
                  >
                    <provider.icon className="h-5 w-5" />
                    <span className="text-xs font-semibold text-center">{provider.label}</span>
                  </button>
                ))}
              </div>
            </div>
          </div>

          <Button
            type="submit"
            className="w-full"
            size="lg"
            disabled={!topic.trim()}
          >
            {t.quiz.startQuiz}
          </Button>
        </form>
      </main>
    </div>
  )
}
