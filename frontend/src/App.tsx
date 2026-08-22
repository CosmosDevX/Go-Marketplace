import { useCallback, useEffect, useRef, useState } from 'react';
import { AlertCircle, PackageOpen, RefreshCw, Loader2 } from 'lucide-react';
import { Header } from './components/Header';
import { FiltersSidebar, type SortBy } from './components/FiltersSidebar';
import { ProductCard } from './components/ProductCard';
import { Pagination } from './components/Pagination';
import { AuthModal } from './components/AuthModal';
import { CartDrawer } from './components/CartDrawer';
import { SellerPanel } from './components/SellerPanel';
import { api, ApiError } from './api/client';
import { useAuth } from './hooks/useAuth';
import { useCart } from './hooks/useCart';
import type { Category, Product } from './types';

type View = 'catalog' | 'seller' | 'admin';

function App() {
  const auth = useAuth();
  const cart = useCart(auth.isAuthenticated);

  const [view, setView] = useState<View>('catalog');
  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [page, setPage] = useState(1);
  const [activeCategory, setActiveCategory] = useState<string | null>(null);
  const [sortBy, setSortBy] = useState<SortBy>(null);
  const [sortAsc, setSortAsc] = useState(true);

  const [categoriesLoading, setCategoriesLoading] = useState(true);
  const [productsLoading, setProductsLoading] = useState(false);
  const [showLoader, setShowLoader] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [authModal, setAuthModal] = useState<'login' | 'register' | null>(null);
  const [cartOpen, setCartOpen] = useState(false);

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
          setError(e instanceof ApiError ? e.message : 'Не удалось загрузить категории');
        }
      } finally {
        if (!cancelled) setCategoriesLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  // If user loses seller access, leave seller view
  useEffect(() => {
    if (view === 'seller' && !auth.canAccessSellerPanel) {
      setView('catalog');
    }
    if (view === 'admin' && !auth.canAccessAdminPanel) {
      setView('catalog');
    }
  }, [auth.canAccessSellerPanel, auth.canAccessAdminPanel, view]);

  const loadProducts = useCallback(async () => {
    const id = ++requestId.current;
    setProductsLoading(true);
    setError(null);

    if (loaderTimer.current) clearTimeout(loaderTimer.current);
    loaderTimer.current = setTimeout(() => {
      if (requestId.current === id) setShowLoader(true);
    }, 250);

    try {
      const data = await api.getProducts({
        page,
        category: activeCategory || undefined,
        sortBy: sortBy || undefined,
        asc: sortBy ? sortAsc : undefined,
      });
      if (requestId.current !== id) return;
      setProducts(Array.isArray(data?.products) ? data.products : []);
    } catch (e) {
      if (requestId.current !== id) return;
      setError(e instanceof ApiError ? e.message : 'Не удалось загрузить товары');
      setProducts([]);
    } finally {
      if (requestId.current === id) {
        if (loaderTimer.current) clearTimeout(loaderTimer.current);
        setProductsLoading(false);
        setShowLoader(false);
      }
    }
  }, [page, activeCategory, sortBy, sortAsc]);

  useEffect(() => {
    if (view !== 'catalog') return;
    loadProducts();
    return () => {
      if (loaderTimer.current) clearTimeout(loaderTimer.current);
    };
  }, [loadProducts, view]);

  const handleCategorySelect = (slug: string | null) => {
    setActiveCategory(slug);
    setPage(1);
  };

  const handleSortChange = (next: SortBy, asc: boolean) => {
    setSortBy(next);
    setSortAsc(asc);
    setPage(1);
  };

  const handleFiltersReset = () => {
    setActiveCategory(null);
    setSortBy(null);
    setSortAsc(true);
    setPage(1);
  };

  const handleCartClick = () => {
    if (!auth.isAuthenticated) {
      setAuthModal('login');
      return;
    }
    setCartOpen(true);
  };

  const hasProducts = products.length > 0;

  return (
    <div className="min-h-screen flex flex-col">
      <Header
        user={auth.user}
        cartCount={cart.totalCount}
        canAccessSellerPanel={auth.canAccessSellerPanel}
        canAccessAdminPanel={auth.canAccessAdminPanel}
        onLogoClick={() => setView('catalog')}
        onLoginClick={() => setAuthModal('login')}
        onRegisterClick={() => setAuthModal('register')}
        onLogout={() => {
          auth.logout();
          setView('catalog');
        }}
        onCartClick={handleCartClick}
        onSellerClick={() => setView('seller')}
        onAdminClick={() => setView('admin')}
      />

      {(view === 'admin' && auth.canAccessAdminPanel) ||
      (view === 'seller' && auth.canAccessSellerPanel) ? (
        <SellerPanel
          onBack={() => setView('catalog')}
          categories={categories}
          isAdmin={view === 'admin'}
          onCategoriesChange={setCategories}
        />
      ) : (
        <main className="flex-1 mx-auto w-full max-w-7xl px-4 sm:px-6 py-6 sm:py-8">
          <div className="mb-6">
            <h1 className="text-2xl sm:text-3xl font-semibold tracking-tight">
              Каталог товаров
            </h1>
            <p className="mt-1 text-sm text-[var(--color-text-muted)]">
              Фильтры слева · сортировка по цене и названию
            </p>
          </div>

          <div className="flex flex-col gap-6 lg:flex-row lg:items-start">
            <FiltersSidebar
              categories={categories}
              categoriesLoading={categoriesLoading}
              activeCategory={activeCategory}
              sortBy={sortBy}
              asc={sortAsc}
              onCategoryChange={handleCategorySelect}
              onSortChange={handleSortChange}
              onReset={handleFiltersReset}
            />

            <div className="min-w-0 flex-1">
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
                {products.map((product, index) => (
                  <ProductCard
                    key={product.product_id}
                    product={product}
                    index={index}
                    isAuthenticated={auth.isAuthenticated}
                    requireAuth={() => setAuthModal('login')}
                    onAddToCart={cart.addItem}
                  />
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
            </div>
          </div>
        </main>
      )}

      <footer className="border-t border-[var(--color-border)] py-6 text-center text-xs text-[var(--color-text-muted)]">
        © 2026 MarketPlace · Тёмная тема · Янтарный акцент
      </footer>

      <AuthModal
        mode={authModal}
        onClose={() => setAuthModal(null)}
        onSwitchMode={(mode) => setAuthModal(mode)}
        onSuccess={(token, username, roles) => {
          auth.login(token, username, roles ?? []);
          setAuthModal(null);
        }}
      />

      <CartDrawer
        open={cartOpen}
        onClose={() => setCartOpen(false)}
        items={cart.items}
        loading={cart.loading}
        error={cart.error}
        totalPrice={cart.totalPrice}
        onChangeQuantity={cart.changeQuantity}
        onRemove={cart.removeItem}
      />
    </div>
  );
}

export default App;
