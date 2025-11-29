'use client';

import Link from 'next/link';
import { ArrowRight } from 'lucide-react';
import { StockDisplayData } from '@/types/stock';

interface StockCardProps {
  stock: StockDisplayData;
  variant?: 'home' | 'dashboard';
}

export default function StockCard({ stock, variant = 'home' }: StockCardProps) {
  if (variant === 'home') {
    return (
      <Link
        href={`/stocks/${stock.ticker}`}
        className="group flex cursor-pointer flex-col gap-4 rounded-xl border border-[#324467] bg-gradient-to-br from-[#192233]/70 to-[#111722]/50 p-5 transition-all duration-300 hover:border-primary hover:shadow-2xl hover:shadow-primary/20"
      >
        <div className="flex items-center justify-between">
          <div className="flex items-baseline gap-3">
            <span className="text-2xl font-bold text-slate-400">{stock.rank}</span>
            <h3 className="text-base font-bold text-white">
              {stock.name} ({stock.ticker})
            </h3>
          </div>
          <div
            className="text-lg font-bold text-[#00FFAB]"
            style={{ textShadow: '0 0 8px rgba(0, 255, 171, 0.5)' }}
          >
            +{stock.changeRate.toFixed(1)}%
          </div>
        </div>
        <div className="mt-auto flex items-center justify-end gap-1 text-sm font-medium text-primary transition-colors group-hover:text-white">
          <span>詳細分析を見る</span>
          <ArrowRight className="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1" />
        </div>
      </Link>
    );
  }

  return null;
}


