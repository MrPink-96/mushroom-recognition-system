import { PredictResponse, ApiError } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function predictMushroom(
  file: File
): Promise<{ data?: PredictResponse; error?: string }> {
  try {
    const formData = new FormData();
    formData.append("file", file, file.name);

    const response = await fetch(`${API_URL}/predict`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      const errorData: ApiError = await response.json().catch(() => ({
        error: "Ошибка сервера",
      }));
      return { error: errorData.error || `Ошибка: ${response.status}` };
    }

    const data: PredictResponse = await response.json();
    return { data };
  } catch (error) {
    console.error("Predict error:", error);
    return {
      error:
        "Не удалось подключиться к серверу. Убедитесь, что API Gateway запущен.",
    };
  }
}

export async function checkHealth(): Promise<boolean> {
  try {
    const response = await fetch(`${API_URL}/health`);
    return response.ok;
  } catch {
    return false;
  }
}
