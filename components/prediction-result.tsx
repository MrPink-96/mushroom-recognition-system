"use client";

import { MushroomPrediction } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ToxicityBadge } from "@/components/toxicity-badge";
import { EdibilityBadge } from "@/components/edibility-badge";
import { ConfidenceCircle } from "@/components/confidence-circle";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertTriangle, ShieldAlert, Info } from "lucide-react";

interface PredictionResultProps {
  prediction: MushroomPrediction;
  isTopResult?: boolean;
}

export function PredictionResult({
  prediction,
  isTopResult = false,
}: PredictionResultProps) {
  const isDangerous = prediction.toxicity_level >= 2;
  const isLowConfidence = prediction.confidence < 0.7;

  if (!isTopResult) {
    // Compact card for other results
    return (
      <Card className="overflow-hidden transition-shadow hover:shadow-md">
        <div className="flex">
          {prediction.images && prediction.images.length > 0 && (
            <div className="h-24 w-24 flex-shrink-0">
              <img
                src={prediction.images[0]}
                alt={prediction.common_name}
                className="h-full w-full object-cover"
              />
            </div>
          )}
          <CardContent className="flex flex-1 items-center justify-between p-4">
            <div className="flex-1">
              <h4 className="font-semibold text-foreground">
                {prediction.common_name}
              </h4>
              <p className="text-sm italic text-muted-foreground">
                {prediction.scientific_name}
              </p>
              <div className="mt-2 flex gap-2">
                <ToxicityBadge level={prediction.toxicity_level} />
              </div>
            </div>
            <ConfidenceCircle confidence={prediction.confidence} size="sm" />
          </CardContent>
        </div>
      </Card>
    );
  }

  // Full card for top result
  return (
    <div className="space-y-4">
      {/* Warnings */}
      {isDangerous && (
        <Alert variant="destructive">
          <ShieldAlert className="h-4 w-4" />
          <AlertTitle>Внимание! Опасный гриб</AlertTitle>
          <AlertDescription>
            Этот гриб потенциально ядовит или смертельно опасен. Ни в коем случае
            не употребляйте его в пищу!
          </AlertDescription>
        </Alert>
      )}

      {isLowConfidence && (
        <Alert variant="warning">
          <AlertTriangle className="h-4 w-4" />
          <AlertTitle>Низкая точность распознавания</AlertTitle>
          <AlertDescription>
            Точность распознавания ниже 70%. Рекомендуем сделать дополнительное
            фото или проконсультироваться со специалистом.
          </AlertDescription>
        </Alert>
      )}

      {/* Main result card */}
      <Card className="overflow-hidden">
        <div className="grid md:grid-cols-2">
          {/* Image */}
          <div className="relative aspect-square bg-muted md:aspect-auto">
            {prediction.images && prediction.images.length > 0 ? (
              <img
                src={prediction.images[0]}
                alt={prediction.common_name}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full items-center justify-center">
                <Info className="h-16 w-16 text-muted-foreground" />
              </div>
            )}
          </div>

          {/* Info */}
          <div className="flex flex-col">
            <CardHeader>
              <div className="flex items-start justify-between gap-4">
                <div>
                  <CardTitle className="text-2xl">
                    {prediction.common_name}
                  </CardTitle>
                  <p className="mt-1 text-lg italic text-muted-foreground">
                    {prediction.scientific_name}
                  </p>
                </div>
                <ConfidenceCircle confidence={prediction.confidence} size="lg" />
              </div>
            </CardHeader>

            <CardContent className="flex-1 space-y-4">
              {/* Badges */}
              <div className="flex flex-wrap gap-2">
                <Badge variant="secondary">{prediction.category.name}</Badge>
                <EdibilityBadge edibility={prediction.edibility} />
                <ToxicityBadge level={prediction.toxicity_level} />
              </div>

              {/* Description */}
              <div>
                <h4 className="mb-2 font-semibold text-foreground">Описание</h4>
                <p className="text-sm leading-relaxed text-muted-foreground">
                  {prediction.description}
                </p>
              </div>
            </CardContent>
          </div>
        </div>
      </Card>
    </div>
  );
}
