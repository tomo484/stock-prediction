'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import Sidebar from '@/components/Sidebar';
import DateNavigation from '@/components/DateNavigation';
import StockDetailCard from '@/components/StockDetailCard';
import SyncButton from '@/components/SyncButton';
import { useLatestStocks, useDateStocks } from '@/hooks/useStocks';

function DashboardPageContent() {
  const searchParams = useSearchParams();
  const dateParam = searchParams.get('date');
  const { data: latestData, isLoading: latestLoading } = useLatestStocks();
  const { data: dateData, isLoading: dateLoading } = useDateStocks(dateParam || '');

  const stocks = dateParam ? dateData : latestData;
  const isLoading = dateParam ? dateLoading : latestLoading;

  return (
    <div className="flex min-h-screen bg-background-light dark:bg-background-dark">
      <Sidebar />
      <main className="flex-1 p-6 lg:p-8">
        <div className="max-w-7xl mx-auto">
          {/* Page Heading */}
          <div className="flex flex-wrap justify-between items-center gap-4 mb-6">
            <div className="flex flex-col gap-1">
              <p className="text-3xl font-bold tracking-tight">メインダッシュボード</p>
              <p className="text-text-secondary-light dark:text-text-secondary-dark text-base font-normal">
                AIによる急騰株分析
              </p>
            </div>
            <SyncButton />
          </div>

          {/* Date Navigation */}
          <DateNavigation variant="dashboard" dateParam={dateParam} />

          {/* Headline */}
          <h2 className="text-2xl font-bold tracking-tight px-1 pb-4 pt-2">本日の急騰株 Top 5</h2>

          {/* Card List */}
          {isLoading ? (
            <div className="grid grid-cols-1 gap-6">
              {Array.from({ length: 5 }).map((_, i) => (
                <div
                  key={i}
                  className="bg-surface-light dark:bg-surface-dark p-6 rounded-lg border border-gray-200 dark:border-gray-700 animate-pulse"
                >
                  <div className="h-6 bg-gray-700 rounded w-1/3 mb-4" />
                  <div className="h-4 bg-gray-700 rounded w-2/3" />
                </div>
              ))}
            </div>
          ) : stocks && stocks.length > 0 ? (
            <div className="grid grid-cols-1 gap-6">
              {stocks.map((stock) => (
                <StockDetailCard key={stock.ticker} stock={stock} />
              ))}
            </div>
          ) : (
            <div className="bg-surface-light dark:bg-surface-dark p-6 rounded-lg border border-gray-200 dark:border-gray-700">
              <p className="text-center text-text-secondary-light dark:text-text-secondary-dark">
                データがありません
              </p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default function DashboardPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background-light dark:bg-background-dark" />}>
      <DashboardPageContent />
    </Suspense>
  );
}

