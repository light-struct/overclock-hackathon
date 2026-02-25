"use client"

import {
  BrainCircuit,
  MessageSquareText,
  BarChart3,
  Zap,
  BookOpen,
  Target,
} from "lucide-react"
import { useApp } from "@/lib/app-context"

export function FeaturesSection() {
  const { t } = useApp()

  const features = [
    {
      icon: BrainCircuit,
      title: t.features.feature1Title,
      description: t.features.feature1Desc,
    },
    {
      icon: MessageSquareText,
      title: t.features.feature2Title,
      description: t.features.feature2Desc,
    },
    {
      icon: BarChart3,
      title: t.features.feature3Title,
      description: t.features.feature3Desc,
    },
    {
      icon: Zap,
      title: t.features.feature4Title,
      description: t.features.feature4Desc,
    },
    {
      icon: BookOpen,
      title: t.features.feature5Title,
      description: t.features.feature5Desc,
    },
    {
      icon: Target,
      title: t.features.feature6Title,
      description: t.features.feature6Desc,
    },
  ]

  return (
    <section id="features" className="border-t border-border px-6 py-24">
      <div className="mx-auto max-w-6xl">
        <div className="mb-16 text-center">
          <p className="mb-3 text-sm font-bold uppercase tracking-widest text-primary">
            {t.features.title}
          </p>
          <h2 className="text-balance text-3xl font-bold text-foreground md:text-4xl">
            {t.features.subtitle}
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
