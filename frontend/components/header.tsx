"use client"

import Link from "next/link"
import { GraduationCap, LogOut, Globe, User } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useApp } from "@/lib/app-context"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export function Header() {
  const { token, logout, language, setLanguage, t } = useApp()
  

  return (
    
    <header className="sticky top-0 z-50 border-b border-border bg-background/90 backdrop-blur-md">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
            <GraduationCap className="h-5 w-5 text-primary-foreground" />
          </div>
          <span className="text-lg font-bold text-foreground">
            QuizAgent
          </span>
        </Link>
        <nav className="hidden items-center gap-8 md:flex">
          <a
            href="#features"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            {t.header.features}
          </a>
          <a
            href="#how-it-works"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            {t.header.howItWorks}
          </a>
          <a
            href="/about"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            {t.header.about}
          </a>
        </nav>
        
        <div className="flex items-center gap-3">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm">
                <Globe className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem
                className={language === 'en' ? 'font-bold bg-gray-100' : ''}
                onClick={() => setLanguage('en')}
              >
                {language === 'en' && '✓ '}English
              </DropdownMenuItem>
              <DropdownMenuItem
                className={language === 'ru' ? 'font-bold bg-gray-100' : ''}
                onClick={() => setLanguage('ru')}
              >
                {language === 'ru' && '✓ '}Русский
              </DropdownMenuItem>
              <DropdownMenuItem
                className={language === 'kk' ? 'font-bold bg-gray-100' : ''}
                onClick={() => setLanguage('kk')}
              >
                {language === 'kk' && '✓ '}Қазақша
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          {token ? (
            <>
              <Button asChild size="sm">
                <Link href="/quiz">{t.header.startQuiz}</Link>
              </Button>
              <Button asChild size="sm" variant="ghost">
                <Link href="/profile" className="flex items-center gap-2">
                  <User className="h-4 w-4" />
                  <span className="hidden sm:inline">{t.header.profile}</span>
                </Link>
              </Button>
              <Button variant="ghost" size="sm" onClick={logout}>
                <LogOut className="h-4 w-4" />
              </Button>
            </>
          ) : (
            <Button asChild size="sm">
              <Link href="/login">{t.auth.login}</Link>
            </Button>
          )}
        </div>
      </div>
    </header>
  )
}
