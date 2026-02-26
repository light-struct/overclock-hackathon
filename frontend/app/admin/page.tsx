"use client"

import { Header } from '@/components/header'
import { Footer } from '@/components/footer'
import { AdminPanel } from '@/components/admin-panel'

export default function AdminPage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />

      <main className="container py-12">
        <AdminPanel />
      </main>

      <Footer />
    </div>
  )
}
