'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import Header from '@/components/Header';
import Footer from '@/components/Footer';
import DateNavigation from '@/components/DateNavigation';
import StockCard from '@/components/StockCard';
import { useLatestStocks, useDateStocks } from '@/hooks/useStocks';

function HomePageContent() {
  const searchParams = useSearchParams();
  const dateParam = searchParams.get('date');
  const { data: latestData, isLoading: latestLoading } = useLatestStocks();
  const { data: dateData, isLoading: dateLoading } = useDateStocks(dateParam || '');

  const stocks = dateParam ? dateData : latestData;
  const isLoading = dateParam ? dateLoading : latestLoading;

  return (
    <div className="relative flex min-h-screen w-full flex-col bg-background-light dark:bg-background-dark group/design-root overflow-x-hidden">
      <div className="layout-container flex h-full grow flex-col">
        <Header />
        <main className="flex-grow">
          <div className="@container">
            {/* Hero Section */}
            <div className="flex min-h-[480px] flex-col gap-6 items-center justify-center p-4 text-center hero-bg @[480px]:gap-8 relative">
              <div
                className="absolute inset-0 bg-cover bg-center"
                style={{
                  backgroundImage:
                    'linear-gradient(rgba(16, 22, 34, 0.6) 0%, rgba(16, 22, 34, 1) 100%), url(https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?w=1920&q=80)',
                }}
              />
              <div className="relative z-10 flex flex-col gap-2">
                <h1 className="text-white text-4xl font-black leading-tight tracking-[-0.033em] @[480px]:text-5xl @[480px]:font-black @[480px]:leading-tight @[480px]:tracking-[-0.033em]">
                  AIが急騰株を即座に分析
                </h1>
                <h2 className="text-slate-300 text-sm font-normal leading-normal @[480px]:text-base @[480px]:font-normal @[480px]:leading-normal max-w-2xl mx-auto">
                  最新のデータに基づき、AIが選出した注目の急騰株を毎日更新。あなたの投資判断をサポートします。
                </h2>
              </div>
              <Link
                href="/dashboard"
                className="relative z-10 flex min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-10 px-4 @[480px]:h-12 @[480px]:px-5 bg-primary text-white text-sm font-bold leading-normal tracking-[0.015em] @[480px]:text-base @[480px]:font-bold @[480px]:leading-normal @[480px]:tracking-[0.015em] hover:bg-primary/80 transition-colors"
              >
                <span className="truncate">AIによる分析を読む（Dashboardへ）</span>
              </Link>
            </div>

            {/* Top 5 Section */}
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between px-4 sm:px-6 pt-12 sm:pt-16">
              <h2 className="text-white text-[22px] font-bold leading-tight tracking-[-0.015em]">
                Top 5銘柄 AI分析結果
              </h2>
              <DateNavigation variant="home" dateParam={dateParam} />
            </div>

            {/* Stock Cards Grid */}
            {isLoading ? (
              <div className="grid grid-cols-1 gap-4 px-4 sm:px-6 py-6 sm:grid-cols-2 lg:grid-cols-3">
                {Array.from({ length: 5 }).map((_, i) => (
                  <div
                    key={i}
                    className="flex flex-col gap-4 rounded-xl border border-[#324467] bg-gradient-to-br from-[#192233]/70 to-[#111722]/50 p-5 animate-pulse"
                  >
                    <div className="h-6 bg-gray-700 rounded w-1/3" />
                    <div className="h-4 bg-gray-700 rounded w-2/3" />
                  </div>
                ))}
              </div>
            ) : stocks && stocks.length > 0 ? (
              <div className="grid grid-cols-1 gap-4 px-4 sm:px-6 py-6 sm:grid-cols-2 lg:grid-cols-3">
                {stocks.map((stock) => (
                  <StockCard key={stock.ticker} stock={stock} variant="home" />
                ))}
              </div>
            ) : (
              <div className="px-4 sm:px-6 py-6">
                <p className="text-slate-400 text-center">データがありません</p>
              </div>
            )}
          </div>
        </main>
        <Footer />
      </div>
    </div>
  );
}

export default function HomePage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background-light dark:bg-background-dark" />}>
      <HomePageContent />
    </Suspense>
  );
}

