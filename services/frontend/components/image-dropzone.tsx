"use client";

import { useCallback, useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import { cn } from "@/lib/utils";
import { Upload, X, ImageIcon } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ImageDropzoneProps {
    file: File | null;
    onFileSelect: (file: File | null) => void;
    disabled?: boolean;
}

export function ImageDropzone({
    file,
    onFileSelect,
    disabled,
}: ImageDropzoneProps) {
    const [preview, setPreview] = useState<string | null>(null);

    useEffect(() => {
        if (!file) {
            setPreview(null);
            return;
        }

        const url = URL.createObjectURL(file);
        setPreview(url);

        return () => {
            URL.revokeObjectURL(url);
        };
    }, [file]);

    const onDrop = useCallback(
        (acceptedFiles: File[]) => {
            if (acceptedFiles[0]) {
                onFileSelect(acceptedFiles[0]);
            }
        },
        [onFileSelect]
    );

    const removeFile = () => {
        onFileSelect(null);
    };

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: {
            "image/jpeg": [".jpg", ".jpeg"],
            "image/png": [".png"],
        },
        maxFiles: 1,
        disabled,
    });

    if (preview) {
        return (
            <div className="flex flex-col items-center">
                <div className="overflow-hidden rounded-xl border bg-card shadow-sm">
                    <img
                        src={preview}
                        alt="preview"
                        className="max-h-[500px] object-contain"
                    />
                </div>

                <div className="mt-3 flex items-center gap-4">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <ImageIcon className="h-4 w-4" />
                        <span className="truncate max-w-[200px]">
                            {file?.name || "image.png"}
                        </span>
                    </div>

                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={removeFile}
                        disabled={disabled}
                    >
                        <X className="h-4 w-4" />
                        Удалить
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div
            {...getRootProps()}
            className={cn(
                "flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed bg-card p-12",
                isDragActive
                    ? "border-primary bg-primary/5"
                    : "border-muted-foreground/25",
                disabled && "opacity-50"
            )}
        >
            <input {...getInputProps()} />

            <Upload className="h-8 w-8 text-primary" />

            <p className="mt-4 text-lg">
                {isDragActive
                    ? "Отпустите файл здесь"
                    : "Перетащите изображение сюда"}
            </p>

            <p className="text-sm text-muted-foreground">
                или нажмите для выбора файла
            </p>
        </div>
    );
}