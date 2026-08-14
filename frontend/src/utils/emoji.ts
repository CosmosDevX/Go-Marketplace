const CATEGORY_EMOJIS: Record<string, string> = {
  electronics: '📱',
  'ready-made-food': '🍱',
  books: '📚',
  clothing: '👕',
  default: '📦',
};

const PRODUCT_KEYWORDS: [string, string][] = [
  ['телефон', '📱'],
  ['phone', '📱'],
  ['наушник', '🎧'],
  ['headphone', '🎧'],
  ['ноутбук', '💻'],
  ['laptop', '💻'],
  ['книга', '📖'],
  ['book', '📖'],
  ['еда', '🍕'],
  ['food', '🍕'],
  ['пицца', '🍕'],
  ['одежда', '👕'],
  ['shirt', '👕'],
];

export function getProductEmoji(productName: string, categorySlug: string): string {
  const lowerName = productName.toLowerCase();
  for (const [keyword, emoji] of PRODUCT_KEYWORDS) {
    if (lowerName.includes(keyword)) return emoji;
  }
  return CATEGORY_EMOJIS[categorySlug] || CATEGORY_EMOJIS.default;
}

export function formatPrice(price: string): string {
  const num = parseFloat(price);
  if (isNaN(num)) return price;
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(num);
}
