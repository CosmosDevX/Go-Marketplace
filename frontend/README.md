# MarketPlace Frontend

Фронтенд маркетплейса на React + TypeScript + Tailwind CSS + Framer Motion.

## Стек
- React 19 + TypeScript
- Vite
- Tailwind CSS v4
- Framer Motion (анимации)
- Lucide React (иконки)

## Запуск

```bash
cd marketplace-frontend
npm install
npm run dev
```

Откроется на http://localhost:5173

API проксируется на `http://localhost:8080` (см. `vite.config.ts`).

## Эндпоинты, которые использует главная страница

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/v1/categories` | Список категорий |
| GET | `/api/v1/products?page=1` | Товары (все) |
| GET | `/api/v1/products?page=1&category=electronics` | Товары по категории |

## Что реализовано
- Тёмная тема + янтарный акцент
- Адаптивная сетка товаров
- Фильтр по категориям (чипы)
- Пагинация (Предыдущая / Следующая)
- Эмодзи вместо фото
- Скелетоны загрузки
- Обработка ошибок
- Плавные анимации появления карточек

## Структура
```
src/
  api/client.ts       — API-клиент
  components/         — UI-компоненты
  types/              — TypeScript-типы
  utils/emoji.ts      — эмодзи + форматирование цены
  App.tsx             — главная страница
```
