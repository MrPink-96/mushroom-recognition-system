export function Footer() {
  return (
    <footer className="border-t bg-muted/30">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
          <p className="text-sm text-muted-foreground">
            MushroomID - Система распознавания грибов
          </p>
          <p className="text-xs text-muted-foreground">
            Внимание: приложение носит информационный характер. Всегда
            консультируйтесь со специалистом перед употреблением грибов.
          </p>
        </div>
      </div>
    </footer>
  );
}
