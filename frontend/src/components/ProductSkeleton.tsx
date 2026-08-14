export function ProductSkeleton() {
  return (
    <div className="flex flex-col rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] overflow-hidden">
      <div className="h-40 sm:h-44 animate-pulse bg-[#222]" />
      <div className="flex flex-1 flex-col gap-3 p-4">
        <div className="h-4 w-3/4 animate-pulse rounded bg-[#2a2a2a]" />
        <div className="h-3 w-full animate-pulse rounded bg-[#2a2a2a]" />
        <div className="h-3 w-2/3 animate-pulse rounded bg-[#2a2a2a]" />
        <div className="mt-auto flex items-center justify-between pt-2">
          <div className="h-6 w-20 animate-pulse rounded bg-[#2a2a2a]" />
          <div className="h-9 w-24 animate-pulse rounded-xl bg-[#2a2a2a]" />
        </div>
      </div>
    </div>
  );
}
