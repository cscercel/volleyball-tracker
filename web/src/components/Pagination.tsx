type PaginationProps = {
  page: number
  hasMore: boolean
  onPageChange: (page: number) => void
  loading?: boolean
}

export default function Pagination({ page, hasMore, onPageChange, loading }: PaginationProps) {
  return (
    <div className="mt-4 flex items-center justify-center gap-4">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1 || loading}
        className="rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
      >
        ← Prev
      </button>
      <span className="text-sm text-slate-500">Page {page}</span>
      <button
        onClick={() => onPageChange(page + 1)}
        disabled={!hasMore || loading}
        className="rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
      >
        Next →
      </button>
    </div>
  )
}
