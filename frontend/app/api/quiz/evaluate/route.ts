import { NextRequest, NextResponse } from "next/server"

export const maxDuration = 60

interface EvaluateRequest {
  question: string
  userAnswer: string
  correctAnswer: string
  apiKey: string
}

export async function POST(req: NextRequest) {
  try {
    const { question, userAnswer, correctAnswer, apiKey }: EvaluateRequest =
      await req.json()

    if (!apiKey) {
      return NextResponse.json(
        { error: "Gemini API key is required" },
        { status: 400 }
      )
    }

    const prompt = `You are a quiz evaluator. Evaluate the student's answer.

Question: "${question}"
Student's Answer: "${userAnswer}"
Correct Answer: "${correctAnswer}"

Respond with ONLY a valid JSON object (no markdown, no code blocks, no extra text):
{
  "score": <number from 0 to 1, where 1 = fully correct, 0.5 = partially correct, 0 = wrong>,
  "correctAnswer": "${correctAnswer}",
  "explanation": "<brief 1-2 sentence explanation of why the answer is right or wrong>"
}

Rules:
- If the student's answer matches the correct answer (even with slight wording differences), score = 1
- If partially correct, score = 0.5
- If incorrect, score = 0
- Return ONLY the JSON object.`

    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=${apiKey}`,
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

    if (!response.ok) {
      const errorData = await response.json().catch(() => null)
      const errorMessage =
        errorData?.error?.message || `Gemini API error: ${response.status}`
      return NextResponse.json({ error: errorMessage }, { status: 500 })
    }

    const data = await response.json()
    const textContent =
      data?.candidates?.[0]?.content?.parts?.[0]?.text || ""

    let jsonStr = textContent.trim()
    if (jsonStr.startsWith("```")) {
      jsonStr = jsonStr.replace(/^```(?:json)?\n?/, "").replace(/\n?```$/, "")
    }

    const evaluation = JSON.parse(jsonStr)

    return NextResponse.json({
      score: Number(evaluation.score) || 0,
      correctAnswer: evaluation.correctAnswer || correctAnswer,
      explanation: evaluation.explanation || "",
    })
  } catch (error) {
    console.error("Evaluation error:", error)
    return NextResponse.json(
      {
        error:
          error instanceof Error ? error.message : "Failed to evaluate answer",
      },
      { status: 500 }
    )
  }
}
