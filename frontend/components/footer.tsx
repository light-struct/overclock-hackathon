import { GraduationCap } from "lucide-react"

export function Footer() {
  return (
    <footer className="border-t border-border px-6 py-10">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 sm:flex-row">
        <div className="flex items-center gap-2.5">
          <div className="flex h-6 w-6 items-center justify-center rounded bg-primary">
            <GraduationCap className="h-4 w-4 text-primary-foreground" />
          </div>
          <span className="text-sm font-bold text-foreground">QuizAgent</span>
        </div>
        <p className="text-xs text-muted-foreground">
          Powered by Gemini AI and Next.js
        </p>
      </div>
    </footer>
  )
}
