"use client";

import { useState, useEffect } from "react";
import { ImageDropzone } from "@/components/image-dropzone";
import { PredictionResult } from "@/components/prediction-result";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { predictMushroom } from "@/lib/api";
import { MushroomPrediction } from "@/lib/types";
import { Loader2, Search, AlertCircle } from "lucide-react";

export default function RecognizePage() {
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [results, setResults] = useState<MushroomPrediction[] | null>(null);
    const [resetKey, setResetKey] = useState(0);

    // 👉 ВСТАВКА ЧЕРЕЗ CTRL + V
    useEffect(() => {
        const handlePaste = (event: ClipboardEvent) => {
            const items = event.clipboardData?.items;
            if (!items) return;

            for (const item of items) {
                if (item.type.startsWith("image/")) {
                    const file = item.getAsFile();
                    if (file) {
                        setSelectedFile(file);
                        setError(null);
                        setResults(null);

                        // сбрасываем dropzone (важно!)
                        setResetKey((prev) => prev + 1);
                    }
                }
            }
        };

        window.addEventListener("paste", handlePaste);

        return () => {
            window.removeEventListener("paste", handlePaste);
        };
    }, []);

    const handleFileSelect = (file: File | null) => {
        setSelectedFile(file);
        setError(null);
        setResults(null);
    };

    const handleSubmit = async () => {
        if (!selectedFile) return;

        setIsLoading(true);
        setError(null);

        const response = await predictMushroom(selectedFile);

        if (response.error) {
            setError(response.error);
        } else if (response.data) {
            setResults(response.data.data);
        }

        setIsLoading(false);
    };

    const topResult = results?.[0];
    const otherResults = results?.slice(1);

    return (
        <div className="container mx-auto px-4 py-12">
            <div className="mx-auto max-w-4xl">
                {/* Header */}
                <div className="mb-8 text-center">
                    <h1 className="text-3xl font-bold text-foreground md:text-4xl">
                        Распознавание гриба
                    </h1>
                    <p className="mt-2 text-muted-foreground">
                        Загрузите фотографию гриба
                    </p>
                    <p className="text-sm text-muted-foreground mt-1">
                        Или вставьте изображение через Ctrl+V
                    </p>
                </div>

                {/* Upload Section */}
                <div className="space-y-6">
                    <ImageDropzone
                        key={resetKey}
                        file={selectedFile}  // 
                        onFileSelect={handleFileSelect}
                        disabled={isLoading}
                    />

                    {selectedFile && !results && (
                        <div className="flex justify-center">
                            <Button
                                size="lg"
                                onClick={handleSubmit}
                                disabled={isLoading}
                                className="w-full max-w-xs"
                            >
                                {isLoading ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        Анализ...
                                    </>
                                ) : (
                                    <>
                                        <Search className="mr-2 h-4 w-4" />
                                        Распознать гриб
                                    </>
                                )}
                            </Button>
                        </div>
                    )}
                </div>

                {/* Error */}
                {error && (
                    <Alert variant="destructive" className="mt-6">
                        <AlertCircle className="h-4 w-4" />
                        <AlertTitle>Ошибка</AlertTitle>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                {/* Results */}
                {results && results.length > 0 && (
                    <div className="mt-12 space-y-8">
                        {/* Top Result */}
                        {topResult && (
                            <section>
                                <h2 className="mb-4 text-xl font-semibold text-foreground">
                                    Наиболее вероятный результат
                                </h2>
                                <PredictionResult prediction={topResult} isTopResult />
                            </section>
                        )}

                        {/* Other Results */}
                        {otherResults && otherResults.length > 0 && (
                            <section>
                                <h2 className="mb-4 text-xl font-semibold text-foreground">
                                    Другие варианты
                                </h2>
                                <div className="grid gap-4 sm:grid-cols-2">
                                    {otherResults.map((prediction) => (
                                        <PredictionResult
                                            key={prediction.id}
                                            prediction={prediction}
                                        />
                                    ))}
                                </div>
                            </section>
                        )}

                        {/* New Search Button */}
                        <div className="flex justify-center pt-4">
                            <Button
                                variant="outline"
                                size="lg"
                                onClick={() => {
                                    setResults(null);
                                    setSelectedFile(null);
                                    setError(null);
                                    setResetKey((prev) => prev + 1);
                                }}
                            >
                                Распознать другой гриб
                            </Button>
                        </div>
                    </div>
                )}

                {/* No Results */}
                {results && results.length === 0 && (
                    <Alert className="mt-6">
                        <AlertCircle className="h-4 w-4" />
                        <AlertTitle>Результаты не найдены</AlertTitle>
                        <AlertDescription>
                            Не удалось определить гриб на изображении. Попробуйте загрузить
                            другое фото с лучшим качеством или освещением.
                        </AlertDescription>
                    </Alert>
                )}
            </div>
        </div>
    );
}