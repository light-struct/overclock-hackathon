import { NextRequest, NextResponse } from "next/server"

export const maxDuration = 60

interface GenerateRequest {
  topic: string
  numQuestions: number
  difficulty: "easy" | "medium" | "hard"
}

export async function POST(req: NextRequest) {
  try {
    const authHeader = req.headers.get('authorization')
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return NextResponse.json(
        { error: "Unauthorized" },
        { status: 401 }
      )
    }

    const apiKey = process.env.GEMINI_API_KEY
    if (!apiKey) {
      return NextResponse.json(
        { error: "Server configuration error" },
        { status: 500 }
      )
    }

    const { topic, numQuestions, difficulty }: GenerateRequest =
      await req.json()

    const difficultyGuide =
      difficulty === "easy"
        ? "Simple recall questions about basic facts and definitions."
        : difficulty === "medium"
          ? "Applied knowledge questions requiring reasoning and connections between concepts."
          : "Advanced analytical questions involving edge cases, synthesis, and deep understanding."

    const prompt = `Generate a quiz about "${topic}" with exactly ${numQuestions} questions at ${difficulty} difficulty level.

${difficultyGuide}

You MUST respond with a valid JSON array only. No markdown, no explanation, no code blocks. Just a raw JSON array.

Each element must be an object with these exact fields:
- "questionNumber": integer starting from 1
- "question": the question text (string)
- "options": array of 4 answer options as strings (A, B, C, D) — even for true/false, provide 4 options
- "correctAnswer": the exact text of the correct option (must match one of the options exactly)

Example format:
[{"questionNumber":1,"question":"What is 2+2?","options":["3","4","5","6"],"correctAnswer":"4"}]

Return ONLY the JSON array.`

    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=${apiKey}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          contents: [{ parts: [{ text: prompt }] }],
          generationConfig: {
            temperature: 0.7,
            maxOutputTokens: 4096,
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

    // Parse JSON from response - handle markdown code blocks
    let jsonStr = textContent.trim()
    if (jsonStr.startsWith("```")) {
      jsonStr = jsonStr.replace(/^```(?:json)?\n?/, "").replace(/\n?```$/, "")
    }

    const questions = JSON.parse(jsonStr)

    if (!Array.isArray(questions) || questions.length === 0) {
      return NextResponse.json(
        { error: "Failed to generate valid questions" },
        { status: 500 }
      )
    }

    return NextResponse.json({ questions })
  } catch (error) {
    console.error("Quiz generation error:", error)
    return NextResponse.json(
      {
        error:
          error instanceof Error ? error.message : "Failed to generate quiz",
      },
      { status: 500 }
    )
  }
}
