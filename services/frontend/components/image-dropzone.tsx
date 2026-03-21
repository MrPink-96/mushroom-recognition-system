"use client";

import { useCallback, useState } from "react";
import { useDropzone } from "react-dropzone";
import { cn } from "@/lib/utils";
import { Upload, X, ImageIcon } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ImageDropzoneProps {
    onFileSelect: (file: File | null) => void;
    disabled?: boolean;
}

export function ImageDropzone({ onFileSelect, disabled }: ImageDropzoneProps) {
    const [preview, setPreview] = useState<string | null>(null);
    const [fileName, setFileName] = useState<string | null>(null);

    const onDrop = useCallback(
        (acceptedFiles: File[]) => {
            const file = acceptedFiles[0];
            if (file) {
                setFileName(file.name);
                const reader = new FileReader();
                reader.onload = () => {
                    setPreview(reader.result as string);
                };
                reader.readAsDataURL(file);
                onFileSelect(file);
            }
        },
        [onFileSelect]
    );

    const removeFile = () => {
        setPreview(null);
        setFileName(null);
        onFileSelect(null);
    };

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: {
            "image/jpeg": [".jpg", ".jpeg"],
            "image/png": [".png"],
        },
        maxFiles: 1,
        maxSize: 10 * 1024 * 1024, // 10MB
        disabled,
    });

    if (preview) {
        return (
            <div className="relative flex flex-col items-center">
                <div className="relative inline-block overflow-hidden rounded-xl border bg-card shadow-sm">
                    <img
                        src={preview}
                        alt="Загруженное изображение"
                        className="max-h-[500px] w-auto object-contain"
                    />
                </div>
                <div className="mt-3 flex items-center gap-4">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <ImageIcon className="h-4 w-4" />
                        <span className="truncate max-w-[200px]">{fileName}</span>
                    </div>
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={removeFile}
                        disabled={disabled}
                        className="text-muted-foreground hover:text-destructive"
                    >
                        <X className="h-4 w-4" />
                        <span className="ml-1">Удалить</span>
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div
            {...getRootProps()}
            className={cn(
                "flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed bg-card p-12 transition-colors",
                isDragActive
                    ? "border-primary bg-primary/5"
                    : "border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/50",
                disabled && "cursor-not-allowed opacity-50"
            )}
        >
            <input {...getInputProps()} />
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
                <Upload className="h-8 w-8 text-primary" />
            </div>
            <p className="mt-4 text-center text-lg font-medium text-foreground">
                {isDragActive ? "Отпустите файл здесь" : "Перетащите изображение сюда"}
            </p>
            <p className="mt-2 text-center text-sm text-muted-foreground">
                или нажмите для выбора файла
            </p>
            <p className="mt-4 text-center text-xs text-muted-foreground">
                Поддерживаются JPG и PNG до 10 МБ
            </p>
        </div>
    );
}
