import { useCallback, useEffect, useRef, useState } from 'react';
import { AlertCircle, PackageOpen, RefreshCw, Loader2 } from 'lucide-react';
import { Header } from './components/Header';
import { CategoryChips } from './components/CategoryChips';
import { ProductCard } from './components/ProductCard';
import { Pagination } from './components/Pagination';
import { AuthModal } from './components/AuthModal';
import { api, ApiError } from './api/client';
import { useAuth } from './hooks/useAuth';
import type { Category, Product } from './types';

function App() {
  const auth = useAuth();

  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [page, setPage] = useState(1);
  const [activeCategory, setActiveCategory] = useState<string | null>(null);

  const [categoriesLoading, setCategoriesLoading] = useState(true);
  const [productsLoading, setProductsLoading] = useState(false);
  const [showLoader, setShowLoader] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [authModal, setAuthModal] = useState<'login' | 'register' | null>(null);

  const loaderTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const requestId = useRef(0);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        setCategoriesLoading(true);
        const data = await api.getCategories();
        if (!cancelled) setCategories(Array.isArray(data) ? data : []);
      } catch (e) {
        if (!cancelled) {
          const msg = e instanceof ApiError ? e.message : 'Не удалось загрузить категории';
          setError(msg);
        }
      } finally {
        if (!cancelled) setCategoriesLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  const loadProducts = useCallback(async () => {
    const id = ++requestId.current;
    setProductsLoading(true);
    setError(null);

    // Показываем лоадер только если запрос длится > 250мс
    if (loaderTimer.current) clearTimeout(loaderTimer.current);
    loaderTimer.current = setTimeout(() => {
      if (requestId.current === id) setShowLoader(true);
    }, 250);

    try {
      const data = await api.getProducts({
        page,
        category: activeCategory || undefined,
      });
      if (requestId.current !== id) return;
      setProducts(Array.isArray(data?.products) ? data.products : []);
    } catch (e) {
      if (requestId.current !== id) return;
      const msg = e instanceof ApiError ? e.message : 'Не удалось загрузить товары';
      setError(msg);
      setProducts([]);
    } finally {
      if (requestId.current === id) {
        if (loaderTimer.current) clearTimeout(loaderTimer.current);
        setProductsLoading(false);
        setShowLoader(false);
      }
    }
  }, [page, activeCategory]);

  useEffect(() => {
    loadProducts();
    return () => {
      if (loaderTimer.current) clearTimeout(loaderTimer.current);
    };
  }, [loadProducts]);

  const handleCategorySelect = (slug: string | null) => {
    setActiveCategory(slug);
    setPage(1);
  };

  const hasProducts = products.length > 0;

  return (
    <div className="min-h-screen flex flex-col">
      <Header
        user={auth.user}
        onLoginClick={() => setAuthModal('login')}
        onRegisterClick={() => setAuthModal('register')}
        onLogout={auth.logout}
      />

      <main className="flex-1 mx-auto w-full max-w-7xl px-4 sm:px-6 py-6 sm:py-8">
        <div className="mb-6">
          <h1 className="text-2xl sm:text-3xl font-semibold tracking-tight">
            Каталог товаров
          </h1>
          <p className="mt-1 text-sm text-[var(--color-text-muted)]">
            Выберите категорию или смотрите все товары
          </p>
        </div>

        <div className="mb-8">
          <CategoryChips
            categories={categories}
            activeSlug={activeCategory}
            onSelect={handleCategorySelect}
            loading={categoriesLoading}
          />
        </div>

        {error && (
          <div className="mb-6 flex items-start gap-3 rounded-2xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
            <AlertCircle size={18} className="mt-0.5 shrink-0" />
            <div className="flex-1">
              <p>{error}</p>
              <button
                type="button"
                onClick={loadProducts}
                className="mt-2 inline-flex items-center gap-1.5 text-[var(--color-accent)] hover:underline"
              >
                <RefreshCw size={14} />
                Попробовать снова
              </button>
            </div>
          </div>
        )}

        <div className="relative min-h-[200px]">
          {/* Лоадер только если запрос реально долгий */}
          {showLoader && (
            <div
              className={
                hasProducts
                  ? 'absolute inset-0 z-10 flex items-start justify-center pt-16 bg-[var(--color-background)]/40'
                  : 'flex items-center justify-center py-24'
              }
            >
              <div className="flex items-center gap-2 rounded-full bg-[var(--color-surface)] px-4 py-2.5 border border-[var(--color-border)]">
                <Loader2 size={18} className="animate-spin text-[var(--color-accent)]" />
                <span className="text-sm text-[var(--color-text-muted)]">Загрузка...</span>
              </div>
            </div>
          )}

          {!productsLoading && !hasProducts && !error && (
            <div className="flex flex-col items-center justify-center py-20 text-center">
              <PackageOpen size={48} className="mb-4 text-[var(--color-text-muted)]" />
              <p className="text-lg font-medium">Товары не найдены</p>
              <p className="mt-1 text-sm text-[var(--color-text-muted)]">
                На этой странице товаров нет. Попробуйте предыдущую страницу или другую категорию.
              </p>
            </div>
          )}

          {hasProducts && (
            <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-5">
              {products.map((product) => (
                <ProductCard key={product.product_id} product={product} />
              ))}
            </div>
          )}
        </div>

        <Pagination
          page={page}
          hasNext={true}
          onPrev={() => setPage((p) => Math.max(1, p - 1))}
          onNext={() => setPage((p) => p + 1)}
          loading={productsLoading}
        />
      </main>

      <footer className="border-t border-[var(--color-border)] py-6 text-center text-xs text-[var(--color-text-muted)]">
        © 2026 MarketPlace · Тёмная тема · Янтарный акцент
      </footer>

      <AuthModal
        mode={authModal}
        onClose={() => setAuthModal(null)}
        onSwitchMode={(mode) => setAuthModal(mode)}
        onSuccess={(token, username) => {
          auth.login(token, username);
          setAuthModal(null);
        }}
      />
    </div>
  );
}

export default App;
