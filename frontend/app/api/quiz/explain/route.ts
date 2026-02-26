import { NextRequest, NextResponse } from "next/server"

export const maxDuration = 60

interface ExplainRequest {
  question: string
  userAnswer: string
  correctAnswer: string
  score: number
  existingExplanation?: string
}

export async function POST(req: NextRequest) {
  try {
    const body: ExplainRequest = await req.json()
    const { question, userAnswer, correctAnswer, score, existingExplanation } = body

    if (!question?.trim()) {
      return NextResponse.json({ error: "Invalid input data" }, { status: 400 })
    }

    const apiKey = process.env.GEMINI_API_KEY
    if (!apiKey) {
      return NextResponse.json({ error: "Server configuration error" }, { status: 500 })
    }

    const wasCorrect = score >= 1
    const prompt = `You are a tutor. Give a clear, educational explanation (3–6 sentences) in the same language as the question.

Question: "${question.replace(/"/g, '\\"')}"
Student's answer: "${(userAnswer || "(no answer)").replace(/"/g, '\\"')}"
Correct answer: "${(correctAnswer || "(none)").replace(/"/g, '\\"')}"
The student's answer was ${wasCorrect ? "correct" : "incorrect"} (score: ${score}).

Explain:
1) Why the correct answer is correct and what concept it reflects.
2) If the student was wrong: what was wrong in their reasoning or answer, and what they should remember for next time.
3) If the student was right: briefly reinforce the key idea.

Write only the explanation text, no JSON, no markdown, no "Explanation:" prefix.`

    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=${apiKey}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          contents: [{ parts: [{ text: prompt }] }],
          generationConfig: {
            temperature: 0.3,
            maxOutputTokens: 1500,
          },
        }),
      }
    )

    if (!response.ok) {
      const errText = await response.text()
      console.error("Gemini explain error:", errText)
      return NextResponse.json({
        explanation: existingExplanation || `Correct answer: ${correctAnswer}.`,
      })
    }

    const data = await response.json()
    let text = data?.candidates?.[0]?.content?.parts?.[0]?.text || ""
    text = text.trim()
    if (!text) text = existingExplanation || `Correct answer: ${correctAnswer}.`

    return NextResponse.json({ explanation: text })
  } catch (error) {
    console.error("Explain error:", error)
    return NextResponse.json(
      { error: "Failed to get explanation" },
      { status: 500 }
    )
  }
}
