import Image from 'next/image'
import Link from 'next/link'
import { Header } from '@/components/header'
import { Footer } from '@/components/footer'
import { ArrowRight, Brain, Sparkles, Target, Users, Mail } from 'lucide-react'

export default function AboutPage() {
  return (
    <main className="min-h-screen bg-background">
      <Header />


      {/* Hero */}
      <section className="relative overflow-hidden px-6 py-24 md:py-36">
        <div className="pointer-events-none absolute -right-40 -top-40 h-[500px] w-[500px] rounded-full bg-primary/5" />
        <div className="pointer-events-none absolute -bottom-20 -left-20 h-[300px] w-[300px] rounded-full bg-primary/5" />

        <div className="relative mx-auto max-w-4xl text-center">
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/20 bg-accent px-4 py-1.5">
            <Sparkles className="h-3.5 w-3.5 text-primary" />
            <span className="text-xs font-semibold text-accent-foreground">
              Overclock Team
            </span>
          </div>

          <h1 className="text-balance text-4xl font-bold leading-tight tracking-tight text-foreground md:text-6xl lg:text-7xl">
            {"Интеллектуальная система "}<span className="text-primary">тестирования</span>
          </h1>

          <p className="mx-auto mt-6 max-w-2xl text-pretty text-lg leading-relaxed text-muted-foreground">
            Мы создаём платформу, которая использует AI для адаптивных тестов.
            Система подбирает вопросы под уровень каждого студента, помогая учиться эффективнее.
          </p>

          <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Link
              href="/quiz"
              className="inline-flex items-center gap-2 rounded-lg bg-primary px-8 py-3 text-sm font-bold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:shadow-xl hover:shadow-primary/30 hover:brightness-110"
            >
              Начать тестирование
              <ArrowRight className="h-4 w-4" />
            </Link>
            <a
              href="#about"
              className="inline-flex items-center gap-2 rounded-lg border border-primary/20 bg-background px-8 py-3 text-sm font-medium text-foreground transition-colors hover:bg-accent"
            >
              Подробнее
            </a>
          </div>

          <div className="mt-20 grid grid-cols-3 gap-8 border-t border-border pt-10">
            <div>
              <p className="text-lg font-bold text-foreground md:text-xl">AI-адаптация</p>
              <p className="mt-1 text-xs text-muted-foreground md:text-sm">Персональные вопросы</p>
            </div>
            <div>
              <p className="text-lg font-bold text-foreground md:text-xl">5 человек</p>
              <p className="mt-1 text-xs text-muted-foreground md:text-sm">Overclock Team</p>
            </div>
            <div>
              <p className="text-lg font-bold text-foreground md:text-xl">Нархоз</p>
              <p className="mt-1 text-xs text-muted-foreground md:text-sm">Колледж, Алматы</p>
            </div>
          </div>
        </div>
      </section>

      {/* Team Photo */}
      <section className="border-y border-border bg-primary/5 px-6 py-16 lg:py-20">
        <div className="mx-auto max-w-4xl">
          <div className="overflow-hidden rounded-2xl border border-primary/20 shadow-2xl shadow-primary/10">
            <Image
              src="/team.jpg"
              alt="Команда Overclock на хакатоне BilimHack Almaty"
              width={1200}
              height={700}
              className="h-auto w-full object-cover opacity-90"
              priority
            />
          </div>
          <p className="mt-4 text-center font-mono text-xs tracking-wide text-muted-foreground uppercase">
            {"BilimHack Almaty \u00b7 Хакатон среди студентов колледжей"}
          </p>
        </div>
      </section>

      {/* О проекте */}
      <section id="about" className="px-6 py-20 lg:py-28">
        <div className="mx-auto max-w-6xl">
          <div className="flex flex-col gap-16 lg:flex-row lg:gap-20">
            <div className="flex-1">
              <span className="font-mono text-xs font-bold tracking-wider text-primary uppercase">О проекте</span>
              <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
                QuizAgent — тесты, которые думают
              </h2>
              <p className="mt-5 text-pretty leading-relaxed text-muted-foreground">
                Мы разрабатываем платформу, которая использует искусственный интеллект для создания
                адаптивных тестов. Система анализирует уровень знаний каждого студента и подбирает
                вопросы индивидуально, помогая эффективнее усваивать материал.
              </p>
              <p className="mt-4 text-pretty leading-relaxed text-muted-foreground">
                Наша цель — сделать процесс оценки знаний не просто проверкой, а инструментом
                обучения, который адаптируется под каждого.
              </p>
            </div>
            <div className="flex flex-1 flex-col gap-5">
              {[
                {
                  icon: Brain,
                  title: "AI-адаптация",
                  desc: "Вопросы подбираются под уровень знаний каждого студента в реальном времени.",
                },
                {
                  icon: Target,
                  title: "Точная оценка",
                  desc: "Система выявляет пробелы в знаниях и помогает сфокусироваться на слабых местах.",
                },
                {
                  icon: Users,
                  title: "Для преподавателей",
                  desc: "Преподаватели получают детальную аналитику по каждому студенту и группе.",
                },
              ].map((feature) => (
                <div
                  key={feature.title}
                  className="group flex gap-4 rounded-xl border border-border bg-card p-5 transition-colors hover:border-primary/30 hover:bg-accent"
                >
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-primary/10">
                    <feature.icon className="h-5 w-5 text-primary" />
                  </div>
                  <div>
                    <h3 className="text-sm font-bold text-foreground">{feature.title}</h3>
                    <p className="mt-1 text-sm leading-relaxed text-muted-foreground">{feature.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Путь */}
      <section className="border-y border-border bg-card px-6 py-20 lg:py-28">
        <div className="mx-auto max-w-3xl">
          <div className="text-center">
            <span className="font-mono text-xs font-bold tracking-wider text-primary uppercase">Наш путь</span>
            <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
              От идеи до продукта
            </h2>
          </div>
          <div className="mt-14 flex flex-col">
            {[
              {
                step: "01",
                title: "Идея",
                desc: "Всё началось с простого вопроса: почему тесты одинаковые для всех, если каждый учится по-разному?",
              },
              {
                step: "02",
                title: "Команда",
                desc: "Пять студентов Нархоз колледжа объединились, чтобы создать решение на основе AI.",
              },
              {
                step: "03",
                title: "BilimHack",
                desc: "Представили прототип на хакатоне BilimHack Almaty и получили обратную связь от экспертов.",
              },
              {
                step: "04",
                title: "Разработка",
                desc: "Сейчас мы активно работаем над платформой, внедряя AI-алгоритмы адаптивного тестирования.",
              },
            ].map((item, i) => (
              <div key={item.step} className="flex gap-6">
                <div className="flex flex-col items-center">
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-primary font-mono text-sm font-bold text-primary-foreground">
                    {item.step}
                  </div>
                  {i < 3 && <div className="mt-2 h-full w-px bg-border" />}
                </div>
                <div className="pb-10">
                  <h3 className="text-base font-bold text-foreground">{item.title}</h3>
                  <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Контакт */}
      <section className="px-6 py-20 lg:py-28">
        <div className="mx-auto max-w-2xl text-center">
          <span className="font-mono text-xs font-bold tracking-wider text-primary uppercase">Связаться с нами</span>
          <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
            Хотите узнать больше?
          </h2>
          <p className="mx-auto mt-4 max-w-md text-pretty leading-relaxed text-muted-foreground">
            Если вам интересен наш проект или вы хотите сотрудничать — напишите нам.
          </p>
          <a
            href="mailto:overclock@narxoz.edu.kz"
            className="mt-8 inline-flex items-center gap-2.5 rounded-xl bg-primary px-7 py-3.5 text-sm font-bold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:shadow-xl hover:shadow-primary/30 hover:brightness-110"
          >
            <Mail className="h-4 w-4" />
            overclocknarxoz@gmail.com
            <ArrowRight className="h-4 w-4" />
          </a> <div></div>
          <a
              href="https://github.com/light-struct/overclock-hackathon"
              className="mt-8 inline-flex items-center gap-2.5 rounded-xl bg-primary px-7 py-3.5 text-sm font-bold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:shadow-xl hover:shadow-primary/30 hover:brightness-110"
          >
            <Image
                src="/GitHub_Invertocat_White.svg"
                alt="GitHub logo"
                width={15}
                height={15}
                priority
            />
            GitHub репозитории проекта
            <ArrowRight className="h-4 w-4" />
          </a>
        </div>
      </section>

      


      <Footer />
    </main>
  )
}