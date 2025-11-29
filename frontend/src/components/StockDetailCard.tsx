'use client';

import Link from 'next/link';
import { ArrowRight, Star } from 'lucide-react';
import { StockDisplayData } from '@/types/stock';

interface StockDetailCardProps {
  stock: StockDisplayData;
}

// モック財務指標データ
const getMockMetrics = (ticker: string) => {
  const mockData: Record<string, { per: string; pbr: string; dividend: string; margin: string }> = {
    default: { per: '15.2x', pbr: '1.8x', dividend: '2.5%', margin: '3.1x' },
  };
  return mockData[ticker] || mockData.default;
};

// モックスコア（85/100など）
const getMockScore = (rank: number) => {
  const scores = [85, 82, 78, 75, 72];
  return scores[rank - 1] || 70;
};

export default function StockDetailCard({ stock }: StockDetailCardProps) {
  const metrics = getMockMetrics(stock.ticker);
  const score = getMockScore(stock.rank);
  const filledStars = Math.floor(score / 20);
  const hasHalfStar = score % 20 >= 10;

  return (
    <div className="bg-surface-light dark:bg-surface-dark p-6 rounded-lg border border-gray-200 dark:border-gray-700">
      <div className="flex flex-wrap justify-between items-start gap-4 mb-4">
        <div className="flex items-center gap-4">
          <span className="text-2xl font-bold text-primary">{stock.rank}位</span>
          <div>
            <Link
              href={`/stocks/${stock.ticker}`}
              className="text-xl font-bold hover:text-primary transition-colors"
            >
              {stock.name} ({stock.ticker})
            </Link>
            <div className="flex items-baseline gap-4 mt-1">
              <p className="text-2xl font-bold">${stock.price.toFixed(2)}</p>
              <p className="text-lg font-semibold text-positive">
                +{stock.changeAmount.toFixed(2)} (+{stock.changeRate.toFixed(2)}%)
              </p>
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <Star
              key={i}
              className={`h-5 w-5 ${
                i < filledStars
                  ? 'fill-yellow-400 text-yellow-400'
                  : i === filledStars && hasHalfStar
                    ? 'fill-yellow-200 text-yellow-400'
                    : 'text-gray-400 dark:text-gray-600'
              }`}
            />
          ))}
          <p className="font-bold text-lg ml-2">{score}/100</p>
        </div>
      </div>
      <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
        <p className="text-sm font-semibold text-text-secondary-light dark:text-text-secondary-dark mb-2">
          AI分析サマリー
        </p>
        <p className="text-base mb-4">{stock.aiAnalysis}</p>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center mb-4 p-4 bg-gray-50 dark:bg-black/20 rounded-DEFAULT">
          <div>
            <p className="text-xs text-text-secondary-light dark:text-text-secondary-dark">PER</p>
            <p className="font-bold">{metrics.per}</p>
          </div>
          <div>
            <p className="text-xs text-text-secondary-light dark:text-text-secondary-dark">PBR</p>
            <p className="font-bold">{metrics.pbr}</p>
          </div>
          <div>
            <p className="text-xs text-text-secondary-light dark:text-text-secondary-dark">
              配当利回り
            </p>
            <p className="font-bold">{metrics.dividend}</p>
          </div>
          <div>
            <p className="text-xs text-text-secondary-light dark:text-text-secondary-dark">
              信用倍率
            </p>
            <p className="font-bold">{metrics.margin}</p>
          </div>
        </div>
        <div className="flex items-center gap-2 text-sm font-medium text-primary hover:underline w-fit cursor-not-allowed opacity-50">
          <span>関連ニュースを見る</span>
          <ArrowRight className="h-4 w-4" />
        </div>
      </div>
    </div>
  );
}


