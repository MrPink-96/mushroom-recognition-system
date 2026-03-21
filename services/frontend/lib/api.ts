import { PredictResponse, ApiError } from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function predictMushroom(
    file: File
): Promise<{ data?: PredictResponse; error?: string }> {
    try {
        console.log("IS FILE:", file instanceof File);
        console.log("TYPE:", file.type);
        console.log("NAME:", file.name);

        const formData = new FormData();
        const fixedFile = new File([file], file.name, {
            type: file.type || "image/jpeg",
        });

        formData.append("file", fixedFile);

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
