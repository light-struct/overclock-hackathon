"use client";

import Link from "next/link";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";
import { useApp } from "@/lib/app-context";

export default function NotFound() {
  const { t } = useApp();

  return (
    <main className="min-h-screen bg-background">
      <Header />
      <div className="container py-24">
        <div className="max-w-3xl mx-auto text-center">
          <div className="inline-flex items-center justify-center mb-6 h-32 w-32 rounded-full bg-[linear-gradient(180deg,#fff,#ffecec)] mx-auto" style={{border: '4px solid var(--brand-red)'}}>
            <span className="text-3xl font-bold text-red-700">404</span>
          </div>

          <h1 className="text-3xl font-bold mb-3">{t.notFound.title}</h1>
          <p className="muted mb-6">{t.notFound.message}</p>

          <div className="flex items-center justify-center gap-3">
            <Link href="/" className="bg-primary text-primary-foreground px-4 py-2 rounded">{t.notFound.goHome}</Link>
            <Link href="/register" className="px-4 py-2 border rounded">{t.notFound.login}</Link>
          </div>

          <div className="mt-10 card text-left">
            <h4 className="font-semibold mb-2">{t.notFound.suggestionsHeading}</h4>
            <ul className="muted-sm list-disc ml-5">
              <li>{t.notFound.suggestion1}</li>
              <li>{t.notFound.suggestion2}</li>
              <li>{t.notFound.suggestion3}</li>
            </ul>
          </div>
        </div>
      </div>
      <Footer />
    </main>
  );
}
