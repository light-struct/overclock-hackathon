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

    // 🤖 AI evaluation — запрос полноценного объяснения от ИИ
    const prompt = `You are a tutor. Evaluate this quiz answer and explain in 2-4 sentences WHY the correct answer is correct and, if the student was wrong, WHY their answer is wrong and what they should remember. Use the same language as the question.

Question: "${question.replace(/"/g, '\\"')}"
Student's Answer: "${userAnswer.replace(/"/g, '\\"')}"
Correct Answer: "${correctAnswer.replace(/"/g, '\\"')}"

Respond ONLY with valid JSON (no markdown, no other text). Use the exact key names. For "correctAnswer" copy the Correct Answer from above.
{
  "score": 0 or 0.5 or 1,
  "correctAnswer": "<same as Correct Answer above>",
  "explanation": "Your clear educational explanation: why the right answer is right; if wrong — what mistake the student made and what to remember."
}`

    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=${apiKey}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          contents: [{ parts: [{ text: prompt }] }],
          generationConfig: {
            temperature: 0.2,
            maxOutputTokens: 800,
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