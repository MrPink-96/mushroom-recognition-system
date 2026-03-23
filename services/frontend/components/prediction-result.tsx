"use client";

import { useState } from "react";
import { MushroomPrediction } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ToxicityBadge } from "@/components/toxicity-badge";
import { EdibilityBadge } from "@/components/edibility-badge";
import { ConfidenceCircle } from "@/components/confidence-circle";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertTriangle, ShieldAlert, Info, X } from "lucide-react";
import { ImageCarousel } from "@/components/image-carousel";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";

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
    const [isModalOpen, setIsModalOpen] = useState(false);

    // =========================
    // COMPACT CARD (другие результаты)
    // =========================
    if (!isTopResult) {
        return (
            <>
                <Card
                    className="overflow-hidden transition-shadow hover:shadow-md cursor-pointer"
                    onClick={() => setIsModalOpen(true)}
                >
                    {/* Высота определяется контентом + фиксированной высотой изображения */}
                    <div className="flex">
                        {/* Изображение */}
                        {prediction.images && prediction.images.length > 0 ? (
                            <div className="w-[38%] max-w-[180px] min-w-[120px] h-48 flex-shrink-0 overflow-hidden rounded-l-lg">
                                <img
                                    src={prediction.images[0]}
                                    alt={prediction.common_name}
                                    className="h-full w-full object-cover"
                                />
                            </div>
                        ) : (
                            <div className="w-28 h-48 flex items-center justify-center bg-muted rounded-l-lg flex-shrink-0">
                                <Info className="h-10 w-10 text-muted-foreground" />
                            </div>
                        )}

                        {/* Контент */}
                        <CardContent className="flex flex-1 flex-col justify-between p-3 pr-4">
                            <div>
                                <h4 className="font-semibold text-foreground leading-tight">
                                    {prediction.common_name}
                                </h4>
                                <p className="text-sm italic text-muted-foreground">
                                    {prediction.scientific_name}
                                </p>

                                <div className="mt-2 flex gap-2">
                                    <ToxicityBadge level={prediction.toxicity_level} />
                                </div>
                            </div>

                            <div className="flex items-center justify-between pb-1">
                                <span className="text-sm text-primary font-medium">
                                    Подробнее →
                                </span>
                                <ConfidenceCircle confidence={prediction.confidence} size="sm" />
                            </div>
                        </CardContent>
                    </div>
                </Card>

                {/* Модальное окно */}
                <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                    <DialogContent className="max-w-3xl p-0">
                        <div className="overflow-hidden rounded-lg">
                            <div className="max-h-[90vh] overflow-y-auto p-6">
                                <button
                                    className="absolute right-6 top-6 z-50 rounded-full p-1.5 text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
                                    onClick={() => setIsModalOpen(false)}
                                    type="button"
                                >
                                    <X className="h-5 w-5" />
                                </button>

                                <DialogHeader>
                                    <div className="flex items-start justify-between">
                                        <div>
                                            <DialogTitle className="text-2xl">
                                                {prediction.common_name}
                                            </DialogTitle>
                                            <p className="text-lg italic text-muted-foreground mt-1">
                                                {prediction.scientific_name}
                                            </p>
                                        </div>
                                        <ConfidenceCircle confidence={prediction.confidence} size="md" />
                                    </div>
                                </DialogHeader>

                                <div className="space-y-4 mt-4">
                                    {prediction.images && prediction.images.length > 0 && (
                                        <div className="rounded-lg overflow-hidden">
                                            <ImageCarousel
                                                images={prediction.images}
                                                alt={prediction.common_name}
                                            />
                                        </div>
                                    )}

                                    <div className="flex flex-wrap gap-2">
                                        <Badge variant="secondary">{prediction.category.name}</Badge>
                                        <EdibilityBadge edibility={prediction.edibility} />
                                        <ToxicityBadge level={prediction.toxicity_level} />
                                    </div>

                                    <div>
                                        <h4 className="mb-2 font-semibold text-foreground">Описание</h4>
                                        <p className="text-sm leading-relaxed text-muted-foreground">
                                            {prediction.description}
                                        </p>
                                    </div>

                                    {prediction.toxicity_level >= 2 && (
                                        <Alert variant="destructive">
                                            <ShieldAlert className="h-4 w-4" />
                                            <AlertTitle>Опасно!</AlertTitle>
                                            <AlertDescription>
                                                Этот гриб потенциально ядовит или смертельно опасен.
                                            </AlertDescription>
                                        </Alert>
                                    )}

                                    {prediction.confidence < 0.7 && (
                                        <Alert variant="warning">
                                            <AlertTriangle className="h-4 w-4" />
                                            <AlertTitle>Низкая точность</AlertTitle>
                                            <AlertDescription>
                                                Точность распознавания ниже 70%.
                                            </AlertDescription>
                                        </Alert>
                                    )}
                                </div>
                            </div>
                        </div>
                    </DialogContent>
                </Dialog>
            </>
        );
    }

    // =========================
    // TOP RESULT (главный результат)
    // =========================
    return (
        <div className="space-y-4">
            {/* Предупреждения */}
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

            {/* Основная карточка */}
            <Card className="overflow-hidden">
                <div className="grid md:grid-cols-2">
                    {/* Изображение / Карусель */}
                    <div className="relative aspect-square bg-muted md:aspect-auto">
                        {prediction.images && prediction.images.length > 0 ? (
                            <ImageCarousel
                                images={prediction.images}
                                alt={prediction.common_name}
                            />
                        ) : (
                            <div className="flex h-full items-center justify-center">
                                <Info className="h-16 w-16 text-muted-foreground" />
                            </div>
                        )}
                    </div>

                    {/* Информация */}
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
                            {/* Бейджи */}
                            <div className="flex flex-wrap gap-2">
                                <Badge variant="secondary">{prediction.category.name}</Badge>
                                <EdibilityBadge edibility={prediction.edibility} />
                                <ToxicityBadge level={prediction.toxicity_level} />
                            </div>

                            {/* Описание */}
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