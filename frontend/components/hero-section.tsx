"use client"

import Link from "next/link"
import { ArrowRight, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useApp } from "@/lib/app-context"

export function HeroSection() {
  const { t } = useApp()
  
  return (
    <section className="relative overflow-hidden px-6 py-24 md:py-36">
      <div className="pointer-events-none absolute -right-40 -top-40 h-[500px] w-[500px] rounded-full bg-primary/5" />
      <div className="pointer-events-none absolute -bottom-20 -left-20 h-[300px] w-[300px] rounded-full bg-primary/5" />

      <div className="relative mx-auto max-w-4xl text-center">
        <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/20 bg-accent px-4 py-1.5">
          <Sparkles className="h-3.5 w-3.5 text-primary" />
          <span className="text-xs font-semibold text-accent-foreground">
            {t.hero.badge}
          </span>
        </div>

        <h1 className="text-balance text-4xl font-bold leading-tight tracking-tight text-foreground md:text-6xl lg:text-7xl">
           {t.hero.title1} <span className="text-primary">{t.hero.title2}</span>
        </h1>

        <p className="mx-auto mt-6 max-w-2xl text-pretty text-lg leading-relaxed text-muted-foreground">
          {t.hero.description}
        </p>

        <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Button asChild size="lg" className="gap-2 px-8">
            <Link href="/quiz">
              {t.hero.startTesting}
              <ArrowRight className="h-4 w-4" />
            </Link>
          </Button>
          <Button asChild variant="outline" size="lg" className="gap-2 px-8 border-primary/20 text-foreground hover:bg-accent">
            <a href="#how-it-works">{t.hero.howItWorks}</a>
          </Button>
        </div>

        <div className="mt-20 grid grid-cols-3 gap-8 border-t border-border pt-10">
          <div>
            <p className="text-lg font-bold text-foreground md:text-xl">{t.hero.stat1Value}</p>
            <p className="mt-1 text-xs text-muted-foreground md:text-sm">{t.hero.stat1Label}</p>
          </div>
          <div>
            <p className="text-lg font-bold text-foreground md:text-xl">{t.hero.stat2Value}</p>
            <p className="mt-1 text-xs text-muted-foreground md:text-sm">{t.hero.stat2Label}</p>
          </div>
          <div>
            <p className="text-lg font-bold text-foreground md:text-xl">{t.hero.stat3Value}</p>
            <p className="mt-1 text-xs text-muted-foreground md:text-sm">{t.hero.stat3Label}</p>
          </div>
        </div>
      </div>
    </section>
  )
}
