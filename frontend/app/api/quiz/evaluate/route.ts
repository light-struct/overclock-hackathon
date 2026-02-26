import { NextRequest, NextResponse } from "next/server"

export const maxDuration = 60

interface EvaluateRequest {
  question: string
  userAnswer: string
  correctAnswer: string
}

export async function POST(req: NextRequest) {
  try {
    const body: EvaluateRequest = await req.json()
    const { question, userAnswer, correctAnswer } = body

    // 🔎 Валидация входных данных
    if (
      !question?.trim() ||
      !userAnswer?.trim() ||
      !correctAnswer?.trim()
    ) {
      return NextResponse.json({ error: "Invalid input data" }, { status: 400 })
    }

    const apiKey = process.env.GEMINI_API_KEY
    if (!apiKey) {
      console.error("❌ GEMINI_API_KEY is missing")
      return NextResponse.json({ error: "Server configuration error" }, { status: 500 })
    }

    // ⚡ Локальная проверка полного совпадения
    const normalizedUser = userAnswer.toLowerCase().trim()
    const normalizedCorrect = correctAnswer.toLowerCase().trim()
    if (normalizedUser === normalizedCorrect) {
      return NextResponse.json({
        score: 1,
        correctAnswer,
        explanation: "Correct answer.",
      })
    }

    // 🤖 AI evaluation
    const prompt = `
You are a quiz evaluator.

Question: "${question}"
Student's Answer: "${userAnswer}"
Correct Answer: "${correctAnswer}"

Respond ONLY with valid JSON:
{
  "score": 0 | 0.5 | 1,
  "correctAnswer": "${correctAnswer}",
  "explanation": "1-2 sentence explanation"
}
`

    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=${apiKey}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          contents: [{ parts: [{ text: prompt }] }],
          generationConfig: {
            temperature: 0.1,
            maxOutputTokens: 512,
          },
        }),
      }
    )

    // fallback если AI упал
    if (!response.ok) {
      const errorText = await response.text()
      console.error("🔥 Gemini API error:", errorText)
      return NextResponse.json({
        score: 0,
        correctAnswer,
        explanation: `Correct answer: ${correctAnswer}`,
      })
    }

    const data = await response.json()
    let textContent = data?.candidates?.[0]?.content?.parts?.[0]?.text || ""
    textContent = textContent.trim()

    // удалить markdown
    if (textContent.startsWith("```")) {
      textContent = textContent.replace(/^```(?:json)?\n?/, "").replace(/\n?```$/, "").trim()
    }

    // извлечь JSON внутри текста
    let evaluation
    try {
      const match = textContent.match(/\{[\s\S]*\}/)
      evaluation = match ? JSON.parse(match[0]) : null
    } catch (err) {
      console.error("❌ JSON parse failed:", err)
    }

    // fallback если JSON не получился
    if (!evaluation) {
      console.warn("⚠ AI returned non-JSON, using fallback")
      return NextResponse.json({
        score: 0,
        correctAnswer,
        explanation: `Correct answer: ${correctAnswer}`,
      })
    }

    return NextResponse.json({
      score: Number(evaluation.score) || 0,
      correctAnswer: evaluation.correctAnswer || correctAnswer,
      explanation: evaluation.explanation || `Correct answer: ${correctAnswer}`,
    })
  } catch (error) {
    console.error("💀 Evaluation fatal error:", error)
    return NextResponse.json(
      { error: "Failed to evaluate answer" },
      { status: 500 }
    )
  }
}