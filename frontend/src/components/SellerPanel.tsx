import { useCallback, useEffect, useRef, useState, FormEvent } from 'react';
import {
  ArrowLeft,
  Loader2,
  Plus,
  Trash2,
  PackageOpen,
  ImageOff,
  FolderPlus,
} from 'lucide-react';
import { api, ApiError } from '../api/client';
import type { Category, Product } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';
import { IMAGE_CONFIG } from '../config/images';
import { Pagination } from './Pagination';
import { ConfirmModal } from './ConfirmModal';
import { AdminRoles } from './AdminRoles';

interface Props {
  onBack: () => void;
  categories: Category[];
  isAdmin: boolean;
  onCategoriesChange: (categories: Category[]) => void;
}

export function SellerPanel({
  onBack,
  categories,
  isAdmin,
  onCategoriesChange,
}: Props) {
  const [products, setProducts] = useState<Product[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [showLoader, setShowLoader] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [showCategoryForm, setShowCategoryForm] = useState(false);

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [categorySlug, setCategorySlug] = useState('');
  const [price, setPrice] = useState('');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const [catName, setCatName] = useState('');
  const [catSlug, setCatSlug] = useState('');
  const [catSubmitting, setCatSubmitting] = useState(false);
  const [catError, setCatError] = useState<string | null>(null);
  const [catSuccess, setCatSuccess] = useState<string | null>(null);

  const [deleteTarget, setDeleteTarget] = useState<Product | null>(null);
  const [deleting, setDeleting] = useState(false);

  const loaderTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const requestId = useRef(0);
  const thumb = IMAGE_CONFIG.SELLER_THUMB_SIZE;

  const loadProducts = useCallback(async () => {
    const id = ++requestId.current;
    setLoading(true);
    setError(null);

    if (loaderTimer.current) clearTimeout(loaderTimer.current);
    loaderTimer.current = setTimeout(() => {
      if (requestId.current === id) setShowLoader(true);
    }, 250);

    try {
      // admin — все товары каталога; seller — только свои
      const data = isAdmin
        ? await api.getProducts({ page })
        : await api.getSellerProducts(page);
      if (requestId.current !== id) return;
      setProducts(Array.isArray(data?.products) ? data.products : []);
    } catch (e) {
      if (requestId.current !== id) return;
      setError(e instanceof ApiError ? e.message : 'Не удалось загрузить товары');
      setProducts([]);
    } finally {
      if (requestId.current === id) {
        if (loaderTimer.current) clearTimeout(loaderTimer.current);
        setLoading(false);
        setShowLoader(false);
      }
    }
  }, [page, isAdmin]);

  useEffect(() => {
    loadProducts();
    return () => {
      if (loaderTimer.current) clearTimeout(loaderTimer.current);
    };
  }, [loadProducts]);

  const resetForm = () => {
    setName('');
    setDescription('');
    setCategorySlug('');
    setPrice('');
    setImageFile(null);
    setFormError(null);
  };

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    setFormError(null);
    setSubmitting(true);

    try {
      const formData = new FormData();
      formData.append('product_name', name);
      formData.append('product_description', description);
      formData.append('category_slug', categorySlug);
      formData.append('product_price', price);
      if (imageFile) {
        formData.append('product_image', imageFile);
      }

      await api.createProduct(formData);
      resetForm();
      setShowForm(false);
      setPage(1);
      const data = await api.getSellerProducts(1);
      setProducts(Array.isArray(data?.products) ? data.products : []);
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : 'Ошибка создания');
    } finally {
      setSubmitting(false);
    }
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    setError(null);
    try {
      await api.deleteProduct(deleteTarget.product_id);
      setProducts((prev) =>
        prev.filter((p) => p.product_id !== deleteTarget.product_id)
      );
      setDeleteTarget(null);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Ошибка удаления');
      setDeleteTarget(null);
    } finally {
      setDeleting(false);
    }
  };

  const handleCreateCategory = async (e: FormEvent) => {
    e.preventDefault();
    setCatError(null);
    setCatSuccess(null);
    setCatSubmitting(true);
    try {
      await api.createCategory({
        category_name: catName,
        category_slug: catSlug,
      });
      const list = await api.getCategories();
      onCategoriesChange(Array.isArray(list) ? list : []);
      setCatName('');
      setCatSlug('');
      setCatSuccess('Категория создана');
    } catch (err) {
      setCatError(err instanceof ApiError ? err.message : 'Ошибка создания категории');
    } finally {
      setCatSubmitting(false);
    }
  };

  const onCatNameChange = (value: string) => {
    setCatName(value);
    setCatSlug(
      value
        .toLowerCase()
        .trim()
        .replace(/\s+/g, '-')
        .replace(/[^a-z0-9а-яё\-]/gi, '')
    );
  };

  const title = isAdmin ? 'Панель администрирования' : 'Панель продавца';

  return (
    <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 py-6 sm:py-8">
      <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={onBack}
            className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
          >
            <ArrowLeft size={16} />
            Каталог
          </button>
          <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        </div>

        <div className="flex flex-wrap gap-2">
          {isAdmin && (
            <button
              type="button"
              onClick={() => {
                setShowCategoryForm((v) => !v);
                setCatError(null);
                setCatSuccess(null);
              }}
              className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2 text-sm font-medium text-[var(--color-text)] hover:border-[var(--color-accent)]/50"
            >
              <FolderPlus size={16} />
              {showCategoryForm ? 'Закрыть' : 'Категория'}
            </button>
          )}
          {!isAdmin && (
            <button
              type="button"
              onClick={() => {
                resetForm();
                setShowForm((v) => !v);
              }}
              className="btn-gradient flex items-center gap-1.5 rounded-xl px-4 py-2 text-sm font-medium text-black"
            >
              <Plus size={16} />
              {showForm ? 'Закрыть форму' : 'Новый товар'}
            </button>
          )}
        </div>
      </div>

      {isAdmin && <AdminRoles />}

      {isAdmin && showCategoryForm && (
        <form
          onSubmit={handleCreateCategory}
          className="mb-6 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-5 space-y-4"
        >
          <h2 className="text-lg font-medium">Создать категорию</h2>
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Название
              </label>
              <input
                required
                value={catName}
                onChange={(e) => onCatNameChange(e.target.value)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
                placeholder="Электроника"
              />
            </div>
            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Slug
              </label>
              <input
                required
                value={catSlug}
                onChange={(e) => setCatSlug(e.target.value)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
                placeholder="electronics"
              />
            </div>
          </div>
          {catError && (
            <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-300">
              {catError}
            </div>
          )}
          {catSuccess && (
            <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-300">
              {catSuccess}
            </div>
          )}
          <button
            type="submit"
            disabled={catSubmitting}
            className="btn-gradient flex items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-semibold text-black disabled:opacity-60"
          >
            {catSubmitting && <Loader2 size={16} className="animate-spin" />}
            Создать категорию
          </button>
        </form>
      )}

      {!isAdmin && showForm && (
        <form
          onSubmit={handleCreate}
          className="mb-8 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-5 space-y-4"
        >
          <h2 className="text-lg font-medium">Создать товар</h2>

          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Название
              </label>
              <input
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
              />
            </div>

            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Цена
              </label>
              <input
                required
                type="number"
                step="0.01"
                min="0"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
                placeholder="199.99"
              />
            </div>
          </div>

          <div>
            <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
              Описание
            </label>
            <textarea
              required
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)] resize-y"
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Категория
              </label>
              <select
                required
                value={categorySlug}
                onChange={(e) => setCategorySlug(e.target.value)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
              >
                <option value="">Выберите категорию</option>
                {categories.map((c) => (
                  <option key={c.category_id} value={c.category_slug}>
                    {c.category_name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                Фото
              </label>
              <input
                type="file"
                accept="image/*"
                onChange={(e) => setImageFile(e.target.files?.[0] ?? null)}
                className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3 py-2 text-sm file:mr-3 file:rounded-lg file:border-0 file:bg-[var(--color-accent)] file:px-3 file:py-1 file:text-sm file:font-medium file:text-black"
              />
            </div>
          </div>

          {formError && (
            <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-300">
              {formError}
            </div>
          )}

          <button
            type="submit"
            disabled={submitting}
            className="btn-gradient flex items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-semibold text-black disabled:opacity-60"
          >
            {submitting && <Loader2 size={16} className="animate-spin" />}
            Создать
          </button>
        </form>
      )}

      {error && (
        <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
          {error}
        </div>
      )}

      <div className="relative min-h-[160px]">
        {showLoader && (
          <div className="flex items-center justify-center py-16">
            <Loader2 size={24} className="animate-spin text-[var(--color-accent)]" />
          </div>
        )}

        {!loading && products.length === 0 && !error && (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <PackageOpen size={40} className="mb-3 text-[var(--color-text-muted)]" />
            <p className="font-medium">Товаров пока нет</p>
            {!isAdmin && (
              <p className="mt-1 text-sm text-[var(--color-text-muted)]">
                Создайте первый товар
              </p>
            )}
          </div>
        )}

        {products.length > 0 && !showLoader && (
          <div className="overflow-x-auto rounded-2xl border border-[var(--color-border)]">
            <table className="w-full min-w-[640px] text-left text-sm">
              <thead className="border-b border-[var(--color-border)] bg-[var(--color-surface)] text-[var(--color-text-muted)]">
                <tr>
                  <th className="px-4 py-3 font-medium">Фото</th>
                  <th className="px-4 py-3 font-medium">Название</th>
                  <th className="px-4 py-3 font-medium">Категория</th>
                  <th className="px-4 py-3 font-medium">Цена</th>
                  <th className="px-4 py-3 font-medium text-right">Действия</th>
                </tr>
              </thead>
              <tbody>
                {products.map((p) => {
                  const src = resolveImageUrl(p.product_image);
                  return (
                    <tr
                      key={p.product_id}
                      className="border-b border-[var(--color-border)] last:border-0"
                    >
                      <td className="px-4 py-3">
                        <div
                          className="overflow-hidden rounded-lg bg-[#1a1a1a]"
                          style={{ width: thumb, height: thumb }}
                        >
                          {src ? (
                            <img
                              src={src}
                              alt=""
                              className="h-full w-full object-cover"
                            />
                          ) : (
                            <div className="flex h-full w-full items-center justify-center text-[var(--color-text-muted)]">
                              <ImageOff size={16} />
                            </div>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <p className="font-medium line-clamp-1">{p.product_name}</p>
                        <p className="text-xs text-[var(--color-text-muted)] line-clamp-1">
                          {p.product_description}
                        </p>
                      </td>
                      <td className="px-4 py-3 text-[var(--color-text-muted)]">
                        {p.category?.category_name ?? '—'}
                      </td>
                      <td className="px-4 py-3 font-medium text-[var(--color-accent)]">
                        {formatPrice(p.product_price)}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <button
                          type="button"
                          onClick={() => setDeleteTarget(p)}
                          className="inline-flex items-center gap-1 rounded-lg border border-[var(--color-border)] px-2.5 py-1.5 text-xs text-[var(--color-text-muted)] hover:border-red-500/40 hover:text-red-400"
                        >
                          <Trash2 size={14} />
                          Удалить
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <Pagination
        page={page}
        hasNext={true}
        onPrev={() => setPage((p) => Math.max(1, p - 1))}
        onNext={() => setPage((p) => p + 1)}
        loading={loading}
      />

      <ConfirmModal
        open={!!deleteTarget}
        title="Удалить товар?"
        message={
          deleteTarget
            ? `«${deleteTarget.product_name}» будет удалён без возможности восстановления.`
            : ''
        }
        confirmLabel="Удалить"
        cancelLabel="Отмена"
        loading={deleting}
        danger
        onConfirm={handleConfirmDelete}
        onCancel={() => !deleting && setDeleteTarget(null)}
      />
    </div>
  );
}
