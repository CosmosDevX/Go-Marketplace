/**
 * Нормализует URL картинки:
 * - полный http://localhost:8080/api/... → /api/... (через Vite proxy)
 * - уже относительный /api/... → как есть
 * - пустое/null → null
 */
export function resolveImageUrl(url: string | null | undefined): string | null {
  if (!url || typeof url !== 'string' || url.trim() === '') return null;

  const trimmed = url.trim();

  // http://localhost:8080/api/v1/uploads/... → /api/v1/uploads/...
  if (trimmed.startsWith('http://localhost:8080/')) {
    return trimmed.replace('http://localhost:8080', '');
  }

  // https variant or other host — оставляем как есть
  if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
    return trimmed;
  }

  // относительный путь без leading slash
  if (trimmed.startsWith('api/')) {
    return `/${trimmed}`;
  }

  // просто имя файла
  if (!trimmed.startsWith('/')) {
    return `/api/v1/uploads/${trimmed}`;
  }

  return trimmed;
}
