export interface QuizConfig {
  topic: string
  numQuestions: number
  difficulty: "easy" | "medium" | "hard"
}

export interface QuestionResult {
  questionNumber: number
  question: string
  userAnswer: string
  correctAnswer: string
  score: number
  explanation: string
}

export interface QuizSummary {
  totalScore: number
  maxScore: number
  percentage: number
  results: QuestionResult[]
}
