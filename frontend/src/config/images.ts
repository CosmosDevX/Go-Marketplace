/**
 * Размеры изображений товаров — меняй здесь в одном месте.
 *
 * PRODUCT_CARD_ASPECT — соотношение сторон блока на карточке каталога
 *   '1 / 1' | '4 / 3' | '3 / 4' | '16 / 9' и т.д.
 *
 * PRODUCT_CARD_MAX_HEIGHT — опциональный потолок высоты (px). null = без ограничения
 *
 * CART_THUMB_SIZE — миниатюра в корзине (px)
 *
 * SELLER_THUMB_SIZE — миниатюра в таблице продавца (px)
 */
export const IMAGE_CONFIG = {
  PRODUCT_CARD_ASPECT: '1 / 1',
  PRODUCT_CARD_MAX_HEIGHT: null as number | null,
  CART_THUMB_SIZE: 80,
  SELLER_THUMB_SIZE: 48,
} as const;
