'use client';

import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import Header from '@/components/Header';
import TimelineItem from '@/components/TimelineItem';
import { useStockHistory } from '@/hooks/useStocks';
import { ChevronLeft } from 'lucide-react';

export default function StockDetailPage() {
  const params = useParams();
  const router = useRouter();
  const ticker = params.ticker as string;
  const { data, isLoading } = useStockHistory(ticker);

  return (
    <div className="relative flex h-auto min-h-screen w-full flex-col">
      <Header />
      <main className="flex w-full flex-1 justify-center py-8 sm:py-12">
        <div className="w-full max-w-4xl flex-1 px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col gap-8">
            {/* Breadcrumb */}
            <div className="flex flex-col gap-4">
              <div className="flex flex-wrap items-center gap-2">
                <Link
                  href="/dashboard"
                  className="text-sm font-medium text-zinc-500 hover:text-primary dark:text-zinc-400 dark:hover:text-primary"
                >
                  ランキング
                </Link>
                <span className="text-sm font-medium text-zinc-400 dark:text-zinc-500">/</span>
                <span className="text-sm font-medium text-zinc-900 dark:text-zinc-100">
                  銘柄詳細
                </span>
              </div>

              {/* Stock Info */}
              {isLoading ? (
                <div className="animate-pulse">
                  <div className="h-8 bg-gray-700 rounded w-1/2 mb-2" />
                  <div className="h-4 bg-gray-700 rounded w-1/4" />
                </div>
              ) : data ? (
                <div className="flex flex-wrap items-center justify-between gap-4">
                  <div className="flex flex-col gap-1">
                    <h2 className="text-3xl font-extrabold tracking-tighter text-zinc-900 dark:text-white sm:text-4xl">
                      {data.stockInfo.name} ({data.stockInfo.ticker})
                    </h2>
                    {data.stockInfo.sector && (
                      <p className="text-base text-zinc-500 dark:text-zinc-400">
                        {data.stockInfo.sector}
                        {data.stockInfo.industry && ` / ${data.stockInfo.industry}`}
                      </p>
                    )}
                  </div>
                </div>
              ) : (
                <div>
                  <h2 className="text-3xl font-extrabold tracking-tighter text-zinc-900 dark:text-white sm:text-4xl">
                    {ticker}
                  </h2>
                  <p className="text-base text-zinc-500 dark:text-zinc-400">データが見つかりません</p>
                </div>
              )}
            </div>

            {/* Timeline */}
            <div className="flex flex-col rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
              <div className="border-b border-zinc-200 p-4 dark:border-zinc-800 sm:p-6">
                <h3 className="text-lg font-bold text-zinc-900 dark:text-white">
                  過去のTop 5ランクイン履歴
                </h3>
              </div>
              <div className="p-4 sm:p-6">
                {isLoading ? (
                  <div className="animate-pulse">
                    <div className="h-20 bg-gray-700 rounded mb-4" />
                    <div className="h-20 bg-gray-700 rounded mb-4" />
                    <div className="h-20 bg-gray-700 rounded" />
                  </div>
                ) : data && data.history.length > 0 ? (
                  <div className="grid grid-cols-[auto_1fr] gap-x-4 sm:gap-x-6">
                    {data.history.map((item, index) => (
                      <TimelineItem
                        key={`${item.date}-${item.rank}`}
                        history={item}
                        isFirst={index === 0}
                        isLast={index === data.history.length - 1}
                      />
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <p className="text-zinc-500 dark:text-zinc-400">
                      ランクイン履歴がありません
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}


