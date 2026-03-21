export interface Category {
  id: number;
  name: string;
}

export interface MushroomPrediction {
  id: number;
  scientific_name: string;
  common_name: string;
  description: string;
  edibility: number; // 0: несъедобный, 1: условно-съедобный, 2: съедобный
  toxicity_level: number; // 0: безопасный, 1: слабо-ядовитый, 2: ядовитый, 3: смертельный
  images: string[];
  category: Category;
  confidence: number; // 0.0 - 1.0
}

export interface PredictResponse {
  data: MushroomPrediction[];
}

export interface ApiError {
  error: string;
}

// Helper functions for display
export function getEdibilityLabel(edibility: number): string {
  switch (edibility) {
    case 0:
      return "Несъедобный";
    case 1:
      return "Условно-съедобный";
    case 2:
      return "Съедобный";
    default:
      return "Неизвестно";
  }
}

export function getToxicityLabel(level: number): string {
  switch (level) {
    case 0:
      return "Безопасный";
    case 1:
      return "Слабо-ядовитый";
    case 2:
      return "Ядовитый";
    case 3:
      return "Смертельно ядовитый";
    default:
      return "Неизвестно";
  }
}

export function getEdibilityColor(edibility: number): "safe" | "warning" | "danger" {
  switch (edibility) {
    case 2:
      return "safe";
    case 1:
      return "warning";
    default:
      return "danger";
  }
}

export function getToxicityColor(level: number): "safe" | "warning" | "danger" {
  switch (level) {
    case 0:
      return "safe";
    case 1:
      return "warning";
    default:
      return "danger";
  }
}
