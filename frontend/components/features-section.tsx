import {
  BrainCircuit,
  MessageSquareText,
  BarChart3,
  Zap,
  BookOpen,
  Target,
} from "lucide-react"

const features = [
  {
    icon: BrainCircuit,
    title: "AI Question Generation",
    description:
      "Gemini AI creates unique questions on any topic, from science to history. Every quiz is different.",
  },
  {
    icon: MessageSquareText,
    title: "Natural Conversation",
    description:
      "Answer in your own words. The AI understands context and evaluates meaning, not just keywords.",
  },
  {
    icon: BarChart3,
    title: "Instant Scoring",
    description:
      "Each answer is scored immediately with the correct answer shown so you learn as you go.",
  },
  {
    icon: Zap,
    title: "Adaptive Difficulty",
    description:
      "Choose easy, medium, or hard. The AI adjusts question complexity to match your level.",
  },
  {
    icon: BookOpen,
    title: "Multiple Formats",
    description:
      "Multiple choice, true/false, short answer, and open-ended questions for thorough assessment.",
  },
  {
    icon: Target,
    title: "Final Results Summary",
    description:
      "After all questions, see your total score, correct answers, and a detailed performance breakdown.",
  },
]

export function FeaturesSection() {
  return (
    <section id="features" className="border-t border-border px-6 py-24">
      <div className="mx-auto max-w-6xl">
        <div className="mb-16 text-center">
          <p className="mb-3 text-sm font-bold uppercase tracking-widest text-primary">
            Features
          </p>
          <h2 className="text-balance text-3xl font-bold text-foreground md:text-4xl">
            Everything you need to learn effectively
          </h2>
        </div>

        <div className="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="group rounded-xl border border-border bg-card p-6 transition-all hover:border-primary/30 hover:shadow-sm"
            >
              <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-lg bg-accent">
                <feature.icon className="h-5 w-5 text-primary" />
              </div>
              <h3 className="mb-2 text-base font-semibold text-card-foreground">
                {feature.title}
              </h3>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
