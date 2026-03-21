import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Camera, Shield, BookOpen } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertTriangle } from "lucide-react";

export default function HomePage() {
  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-b from-primary/5 to-background py-20 md:py-32">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-3xl text-center">
            <h1 className="text-balance text-4xl font-bold tracking-tight text-foreground md:text-5xl lg:text-6xl">
              Распознайте гриб по фотографии
            </h1>
            <p className="mt-6 text-pretty text-lg text-muted-foreground md:text-xl">
              Загрузите фото гриба и получите информацию о его виде, съедобности
              и уровне токсичности с помощью машинного обучения
            </p>
            <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
              <Button asChild size="lg" className="w-full sm:w-auto">
                <Link href="/recognize">Распознать гриб</Link>
              </Button>
            </div>
          </div>
        </div>

        {/* Decorative background */}
        <div className="absolute inset-0 -z-10 overflow-hidden">
          <div className="absolute -top-1/2 left-1/2 h-96 w-96 -translate-x-1/2 rounded-full bg-primary/10 blur-3xl" />
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-center text-3xl font-bold text-foreground">
            Возможности системы
          </h2>
          <div className="mt-12 grid gap-6 md:grid-cols-3">
            <Card className="border-0 bg-muted/50">
              <CardContent className="flex flex-col items-center p-6 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <Camera className="h-6 w-6 text-primary" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-foreground">
                  Быстрое распознавание
                </h3>
                <p className="mt-2 text-sm text-muted-foreground">
                  Загрузите фото и получите результат за несколько секунд.
                  Поддержка форматов JPG и PNG.
                </p>
              </CardContent>
            </Card>

            <Card className="border-0 bg-muted/50">
              <CardContent className="flex flex-col items-center p-6 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <Shield className="h-6 w-6 text-primary" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-foreground">
                  Оценка безопасности
                </h3>
                <p className="mt-2 text-sm text-muted-foreground">
                  Система показывает уровень токсичности и съедобности с
                  предупреждениями об опасности.
                </p>
              </CardContent>
            </Card>

            <Card className="border-0 bg-muted/50">
              <CardContent className="flex flex-col items-center p-6 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <BookOpen className="h-6 w-6 text-primary" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-foreground">
                  Подробная информация
                </h3>
                <p className="mt-2 text-sm text-muted-foreground">
                  Научное и русское название, описание, категория и другие
                  характеристики гриба.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Warning Section */}
      <section className="py-12">
        <div className="container mx-auto px-4">
          <Alert variant="warning" className="mx-auto max-w-2xl">
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription className="ml-2">
              <strong>Важно:</strong> Это приложение создано в учебных целях и
              не заменяет консультацию специалиста-миколога. Никогда не
              употребляйте грибы, в безопасности которых вы не уверены на 100%.
            </AlertDescription>
          </Alert>
        </div>
      </section>
    </div>
  );
}
