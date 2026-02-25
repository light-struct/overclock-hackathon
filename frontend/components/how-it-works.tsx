import Link from "next/link"
import { ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"

const steps = [
  {
    step: "01",
    title: "Choose Your Topic",
    description:
      "Enter any subject you want to be tested on. The AI covers everything from math to philosophy.",
  },
  {
    step: "02",
    title: "Configure the Quiz",
    description:
      "Set the number of questions and difficulty level. Provide your Gemini API key for the AI engine.",
  },
  {
    step: "03",
    title: "Answer Questions",
    description:
      "The AI presents questions one by one, evaluates your answers, and provides the correct answer with a score.",
  },
  {
    step: "04",
    title: "Review Your Results",
    description:
      "After all questions are answered, view your total score, every question with the correct answers, and a final performance summary.",
  },
]

export function HowItWorks() {
  return (
    <section
      id="how-it-works"
      className="border-t border-border bg-secondary px-6 py-24"
    >
      <div className="mx-auto max-w-4xl">
        <div className="mb-16 text-center">
          <p className="mb-3 text-sm font-bold uppercase tracking-widest text-primary">
            How it works
          </p>
          <h2 className="text-balance text-3xl font-bold text-foreground md:text-4xl">
            Four simple steps to get started
          </h2>
        </div>

        <div className="space-y-6">
          {steps.map((item, i) => (
            <div key={item.step} className="flex gap-6">
              <div className="flex flex-col items-center">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary font-mono text-sm font-bold text-primary-foreground">
                  {item.step}
                </div>
                {i < steps.length - 1 && (
                  <div className="my-2 h-full w-px bg-border" />
                )}
              </div>
              <div className="pb-6">
                <h3 className="text-lg font-semibold text-foreground">
                  {item.title}
                </h3>
                <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                  {item.description}
                </p>
              </div>
            </div>
          ))}
        </div>

        <div className="mt-12 text-center">
          <Button asChild size="lg" className="gap-2">
            <Link href="/quiz">
              Start Your Quiz
              <ArrowRight className="h-4 w-4" />
            </Link>
          </Button>
        </div>
      </div>
    </section>
  )
}
