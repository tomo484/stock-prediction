'use client';

import { useSyncStocks } from '@/hooks/useStocks';
import { RefreshCw } from 'lucide-react';

export default function SyncButton() {
  const { mutate: syncStocks, isPending } = useSyncStocks();

  return (
    <button
      onClick={() => syncStocks()}
      disabled={isPending}
      className="flex min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg h-10 px-4 bg-primary text-white text-sm font-bold leading-normal tracking-[0.015em] hover:bg-primary/80 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {isPending ? (
        <>
          <RefreshCw className="h-4 w-4 animate-spin" />
          <span className="truncate">同期中...</span>
        </>
      ) : (
        <>
          <RefreshCw className="h-4 w-4" />
          <span className="truncate">データを更新</span>
        </>
      )}
    </button>
  );
}


