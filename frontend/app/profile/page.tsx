"use client"

import { Header } from '@/components/header'
import { Footer } from '@/components/footer'
import { ProfileCard } from '@/components/profile'

export default function ProfilePage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />

      <main className="container py-12">
        <ProfileCard />
      </main>

      <Footer />
    </div>
  )
}
